package adminapi

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/go-resty/resty/v2"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/infra/conf"

	"github.com/jare-abc/XrayR-dev/api/adminapi"
)

// APIClient create an api client to the panel.
type APIClient struct {
	client        *resty.Client
	APIHost       string
	NodeID        uint64
	Key           string
	NodeType      string
	EnableVless   bool
	VlessFlow     string
	SpeedLimit    float64
	DeviceLimit   int
	LocalRuleList []api.DetectRule
	resp          atomic.Value
	eTag          string
}

// New create an api instance
func New(apiConfig *api.Config) *APIClient {
	client := resty.New()
	client.SetRetryCount(3)
	if apiConfig.Timeout > 0 {
		client.SetTimeout(time.Duration(apiConfig.Timeout) * time.Second)
	} else {
		client.SetTimeout(5 * time.Second)
	}
	client.OnError(func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			// v.Response contains the last response from the server
			// v.Err contains the original error
			log.Print(v.Err)
		}
	})
	client.SetBaseURL(apiConfig.APIHost)
	// Create Key for each requests
	client.SetQueryParams(map[string]string{
		"node_id":   strconv.Itoa(apiConfig.NodeID),
		"node_type": strings.ToLower(apiConfig.NodeType),
		"token":     apiConfig.Key,
	})
	// Read local rule list
	localRuleList := readLocalRuleList(apiConfig.RuleListPath)
	apiClient := &APIClient{
		client:        client,
		NodeID:        apiConfig.NodeID,
		Key:           apiConfig.Key,
		APIHost:       apiConfig.APIHost,
		NodeType:      apiConfig.NodeType,
		EnableVless:   apiConfig.EnableVless,
		VlessFlow:     apiConfig.VlessFlow,
		SpeedLimit:    apiConfig.SpeedLimit,
		DeviceLimit:   apiConfig.DeviceLimit,
		LocalRuleList: localRuleList,
	}
	return apiClient
}

// START GetNodeInfo----------------------------------------------------------------------------------------------------------------
func (c *APIClient) GetNodeInfo() (nodeInfo *api.NodeInfo, err error) {
	server := new(serverConfig)
	path := "/api/node/nodesetting"

	res, err := c.client.R().
		SetHeader("token", c.Key).
		ForceContentType("application/json").
		Get(path)

	nodeInfoResp, err := c.parseResponse(res, path, err)
	if err != nil {
		return nil, err
	}
	b, _ := nodeInfoResp.Encode()
	json.Unmarshal(b, server)

	if server.ServerPort == 0 {
		return nil, errors.New("server port must > 0")
	}

	c.resp.Store(server)

	switch c.NodeType {
	case "Shadowsocks":
		nodeInfo, err = c.parseSSNodeResponse(server)
	default:
		return nil, fmt.Errorf("unsupported node type: %s", c.NodeType)
	}

	if err != nil {
		return nil, fmt.Errorf("parse node info failed: %s, \nError: %v", res.String(), err)
	}

	return nodeInfo, nil
}

func (c *APIClient) ParseSSNodeResponse(nodeInfoResponse *json.RawMessage) (*api.NodeInfo, error) {
	var speedLimit uint64 = 0
	shadowsocksNodeInfo := new(ShadowsocksNodeInfo)
	if err := json.Unmarshal(*nodeInfoResponse, shadowsocksNodeInfo); err != nil {
		return nil, fmt.Errorf("unmarshal %s failed: %s", reflect.TypeOf(*nodeInfoResponse), err)
	}

	if c.DeviceLimit == 0 && shadowsocksNodeInfo.ClientLimit > 0 {
		c.DeviceLimit = shadowsocksNodeInfo.ClientLimit
	}

	c.NodeID = shadowsocksNodeInfo.NodeID

	nodeInfo := &api.NodeInfo{
		NodeType:          c.NodeType,
		NodeID:            shadowsocksNodeInfo.NodeID,
		Port:              shadowsocksNodeInfo.Port,
		SpeedLimit:        speedLimit,
		TransportProtocol: "tcp",
		CypherMethod:      shadowsocksNodeInfo.Setting.Method,
	}

	return nodeInfo, nil
}

func (c *APIClient) ParseV2rayNodeResponse(nodeInfoResponse *json.RawMessage) (*api.NodeInfo, error) {
	var speedLimit uint64 = 0

	v2rayNodeInfo := new(V2rayNodeInfo)
	if err := json.Unmarshal(*nodeInfoResponse, v2rayNodeInfo); err != nil {
		return nil, fmt.Errorf("unmarshal %s failed: %s", reflect.TypeOf(*nodeInfoResponse), err)
	}

	if c.SpeedLimit > 0 {
		speedLimit = uint64((c.SpeedLimit * 1000000) / 8)
	} else {
		speedLimit = (v2rayNodeInfo.SpeedLimit * 1000000) / 8
	}

	if c.DeviceLimit == 0 && v2rayNodeInfo.ClientLimit > 0 {
		c.DeviceLimit = v2rayNodeInfo.ClientLimit
	}

	// Create GeneralNodeInfo
	nodeInfo := &api.NodeInfo{
		NodeType:          c.NodeType,
		NodeID:            c.NodeID,
		Port:              v2rayNodeInfo.V2Port,
		SpeedLimit:        speedLimit,
		AlterID:           v2rayNodeInfo.V2AlterID,
		TransportProtocol: v2rayNodeInfo.V2Net,
		FakeType:          v2rayNodeInfo.V2Type,
		EnableTLS:         v2rayNodeInfo.V2TLS,
		Path:              v2rayNodeInfo.V2Path,
		Host:              v2rayNodeInfo.V2Host,
		EnableVless:       c.EnableVless,
		VlessFlow:         c.VlessFlow,
	}

	return nodeInfo, nil
}

func (c *APIClient) ParseTrojanNodeResponse(nodeInfoResponse *json.RawMessage) (*api.NodeInfo, error) {
	var speedLimit uint64 = 0

	trojanNodeInfo := new(TrojanNodeInfo)
	if err := json.Unmarshal(*nodeInfoResponse, trojanNodeInfo); err != nil {
		return nil, fmt.Errorf("unmarshal %s failed: %s", reflect.TypeOf(*nodeInfoResponse), err)
	}
	if c.SpeedLimit > 0 {
		speedLimit = uint64((c.SpeedLimit * 1000000) / 8)
	} else {
		speedLimit = (trojanNodeInfo.SpeedLimit * 1000000) / 8
	}

	if c.DeviceLimit == 0 && trojanNodeInfo.ClientLimit > 0 {
		c.DeviceLimit = trojanNodeInfo.ClientLimit
	}

	// Create GeneralNodeInfo
	nodeInfo := &api.NodeInfo{
		NodeType:          c.NodeType,
		NodeID:            c.NodeID,
		Port:              trojanNodeInfo.TrojanPort,
		SpeedLimit:        speedLimit,
		TransportProtocol: "tcp",
		EnableTLS:         true,
	}

	return nodeInfo, nil
}

//END GetNodeInfo----------------------------------------------------------------------------------------------------------------

// START GetUserList----------------------------------------------------------------------------------------------------------------
func (c *APIClient) GetUserList() (UserList *[]api.UserInfo, err error) {
	path := "/api/node/usenodeuser"
	var nodeType = ""
	switch c.NodeType {
	case "Shadowsocks":
		nodeType = "ss"
	case "V2ray":
		nodeType = "v2ray"
	case "Trojan":
		nodeType = "trojan"
	default:
		return nil, fmt.Errorf("NodeType Error: %s", c.NodeType)
	}
	res, err := c.client.R().
		SetResult(&Response{}).
		ForceContentType("application/json").
		Get(path)

	response, err := c.parseResponse(res, path, err)
	if err != nil {
		return nil, err
	}

	var userListResponse *[]UserResponse
	if err := json.Unmarshal(response.Data, &userListResponse); err != nil {
		return nil, fmt.Errorf("unmarshal %s failed: %s", reflect.TypeOf(userListResponse), err)
	}
	userList, err := c.ParseUserListResponse(userListResponse)
	if err != nil {
		res, _ := json.Marshal(userListResponse)
		return nil, fmt.Errorf("parse user list failed: %s", string(res))
	}
	return userList, nil
}

func (c *APIClient) ParseUserListResponse(userInfoResponse *[]UserResponse) (*[]api.UserInfo, error) {
	var deviceLimit = 0
	var speedLimit uint64 = 0
	userList := make([]api.UserInfo, len(*userInfoResponse))
	for i, user := range *userInfoResponse {
		if c.DeviceLimit > 0 {
			deviceLimit = c.DeviceLimit
		} else {
			deviceLimit = user.DeviceLimit
		}

		if c.SpeedLimit > 0 {
			speedLimit = uint64((c.SpeedLimit * 1000000) / 8)
		} else if user.SpeedLimit > 0 {
			speedLimit = uint64((user.SpeedLimit * 1000000) / 8)
		}
		userList[i] = api.UserInfo{
			UID:         user.ID,
			Passwd:      user.Passwd,
			UUID:        user.Passwd,
			SpeedLimit:  speedLimit,
			DeviceLimit: deviceLimit,
		}
	}

	return &userList, nil
}

//END GetUserList----------------------------------------------------------------------------------------------------------------

// START GetNodeRule----------------------------------------------------------------------------------------------------------------
func (c *APIClient) GetNodeRule() (*[]api.DetectRule, error) {
	ruleList := c.LocalRuleList
	path := "/api/node/noderules"
	res, err := c.client.R().
		SetResult(&Response{}).
		ForceContentType("application/json").
		Get(path)

	response, err := c.parseResponse(res, path, err)
	if err != nil {
		return nil, err
	}

	ruleListResponse := new(NodeRule)

	if err := json.Unmarshal(response.Data, ruleListResponse); err != nil {
		return nil, fmt.Errorf("unmarshal %s failed: %s", reflect.TypeOf(ruleListResponse), err)
	}
	ruleList := c.LocalRuleList
	// Only support reject rule type
	if ruleListResponse.Mode != "reject" {
		return &ruleList, nil
	} else {
		for _, r := range ruleListResponse.Rules {
			if r.Type == "reg" {
				ruleList = append(ruleList, api.DetectRule{
					ID:      r.ID,
					Pattern: regexp.MustCompile(r.Pattern),
				})
			}

		}
	}

	return &ruleList, nil
}

//END GetNodeRule----------------------------------------------------------------------------------------------------------------

// START ReportNodeStatus--------------------------------------------------------------------------------------------------------------
func (c *APIClient) ReportNodeStatus(nodeStatus *api.NodeStatus) (err error) {
	path := "/api/node/checkonline"
	var nodeType = ""
	switch c.NodeType {
	case "Shadowsocks":
		nodeType = "ss"
	case "V2ray":
		nodeType = "v2ray"
	case "Trojan":
		nodeType = "trojan"
	default:
		return nil, fmt.Errorf("NodeType Error: %s", c.NodeType)
	}

	systemload := NodeStatus{
		Uptime: int(nodeStatus.Uptime),
		CPU:    fmt.Sprintf("%d%%", int(nodeStatus.CPU)),
		Mem:    fmt.Sprintf("%d%%", int(nodeStatus.Mem)),
		Disk:   fmt.Sprintf("%d%%", int(nodeStatus.Disk)),
	}

	res, err := c.createCommonRequest().
		SetBody(systemload).
		SetResult(&Response{}).
		ForceContentType("application/json").
		Post(path)

	_, err = c.parseResponse(res, path, err)
	if err != nil {
		return err
	}

	return nil
}

//END ReportNodeStatus----------------------------------------------------------------------------------------------------------------

// START ReportNodeOnlineUsers--------------------------------------------------------------------------------------------------------------
func (c *APIClient) ReportNodeOnlineUsers(onlineUserList *[]api.OnlineUser) error {
	return nil
}

//END ReportNodeOnlineUsers----------------------------------------------------------------------------------------------------------------

// START ReportNodeOnlineUsers--------------------------------------------------------------------------------------------------------------
func (c *APIClient) ReportUserTraffic(userTraffic *[]api.UserTraffic) error {

	data := make([]UserTraffic, len(*userTraffic))
	for i, traffic := range *userTraffic {
		data[i] = UserTraffic{
			UserId:   traffic.UID,
			Upload:   traffic.Upload,
			Download: traffic.Download}
	}
	postData := &PostData{data: data, nodeid: strconv.Itoa(c.NodeID)}
	path := "/api/node/usedtraffic"
	res, err := c.client.R().
		SetBody(postData).
		SetResult(&Response{}).
		ForceContentType("application/json").
		Post(path)
	_, err = c.parseResponse(res, path, err)
	if err != nil {
		return err
	}

	return nil
}

//END ReportNodeOnlineUsers----------------------------------------------------------------------------------------------------------------

// START Describe--------------------------------------------------------------------------------------------------------------
func (c *APIClient) Describe() api.ClientInfo {
	return api.ClientInfo{APIHost: c.APIHost, NodeID: c.NodeID, Key: c.Key, NodeType: c.NodeType}
}

//END Describe----------------------------------------------------------------------------------------------------------------

// START GetNodeRule--------------------------------------------------------------------------------------------------------------
func (c *APIClient) GetNodeRule() (*[]api.DetectRule, error) {
	path := "/api/node/rules"
	var nodeType = ""
	switch c.NodeType {
	case "Shadowsocks":
		nodeType = "ss"
	case "V2ray":
		nodeType = "v2ray"
	case "Trojan":
		nodeType = "trojan"
	default:
		return nil, fmt.Errorf("NodeType Error: %s", c.NodeType)
	}

	res, err := c.createCommonRequest().
		SetResult(&Response{}).
		ForceContentType("application/json").
		Get(path)

	response, err := c.parseResponse(res, path, err)
	if err != nil {
		return nil, err
	}

	ruleListResponse := new(NodeRule)

	if err := json.Unmarshal(response.Data, ruleListResponse); err != nil {
		return nil, fmt.Errorf("unmarshal %s failed: %s", reflect.TypeOf(ruleListResponse), err)
	}
	ruleList := c.LocalRuleList
	// Only support reject rule type
	if ruleListResponse.Mode != "reject" {
		return &ruleList, nil
	} else {
		for _, r := range ruleListResponse.Rules {
			if r.Type == "reg" {
				ruleList = append(ruleList, api.DetectRule{
					ID:      r.ID,
					Pattern: regexp.MustCompile(r.Pattern),
				})
			}

		}
	}

	return &ruleList, nil
}

//END GetNodeRule----------------------------------------------------------------------------------------------------------------

// START ReportIllegal--------------------------------------------------------------------------------------------------------------
func (c *APIClient) ReportIllegal(detectResultList *[]api.DetectResult) error {

	data := make([]IllegalItem, len(*detectResultList))
	for i, r := range *detectResultList {
		data[i] = IllegalItem{
			ID:     r.RuleID,
			UserId: r.UID,
		}
	}
	postData := &PostData{data}
	path := "/api/node/detectlog"
	res, err := c.client.R().
		SetBody(postData).
		SetResult(&Response{}).
		ForceContentType("application/json").
		Post(path)
	_, err = c.parseResponse(res, path, err)
	if err != nil {
		return err
	}
	return nil
}

//END ReportIllegal----------------------------------------------------------------------------------------------------------------

func (c *APIClient) Debug() {
	c.client.SetDebug(true)
}
