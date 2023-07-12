package adminapi_test

import (
	"fmt"
	"testing"

	"github.com/jare-abc/XrayR-dev/api"
	"github.com/jare-abc/XrayR-dev/api/adminapi"
)

func CreateClient() api.API {
	apiConfig := &api.Config{
		APIHost:  "http://127.0.0.1:5005",
		Key:      "123456789",
		NodeID:   1,
		NodeType: "Shadowsocks",
	}
	client := adminapi.New(apiConfig)
	return client
}

func TestGetV2rayNodeinfo(t *testing.T) {
	apiConfig := &api.Config{
		APIHost:  "http://127.0.0.1:8888",
		Key:      "naBDpLvREiwY9qPr",
		NodeID:   1,
		NodeType: "V2ray",
	}
	client := adminapi.New(apiConfig)

	nodeInfo, err := client.GetNodeInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodeInfo)
}

func TestGetSSNodeinfo(t *testing.T) {
	apiConfig := &api.Config{
		APIHost:  "http://127.0.0.1:5005",
		Key:      "123456789",
		NodeID:   3,
		NodeType: "Shadowsocks",
	}
	client := adminapi.New(apiConfig)
	nodeInfo, err := client.GetNodeInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodeInfo)
}

func TestGetTrojanNodeinfo(t *testing.T) {
	apiConfig := &api.Config{
		APIHost:  "http://127.0.0.1:8888",
		Key:      "kgnO2O66FmvP8rDV",
		NodeID:   2,
		NodeType: "Trojan",
	}
	client := adminapi.New(apiConfig)
	nodeInfo, err := client.GetNodeInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodeInfo)
}

func TestGetSSinfo(t *testing.T) {
	client := CreateClient()

	nodeInfo, err := client.GetNodeInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodeInfo)
}

func TestGetUserList(t *testing.T) {
	client := CreateClient()

	userList, err := client.GetUserList()
	if err != nil {
		t.Error(err)
	}

	t.Log(userList)
}

func TestReportNodeStatus(t *testing.T) {
	client := CreateClient()
	nodeStatus := &api.NodeStatus{
		CPU: 1, Mem: 1, Disk: 1, Uptime: 256,
	}
	err := client.ReportNodeStatus(nodeStatus)
	if err != nil {
		t.Error(err)
	}
}

func TestReportReportNodeOnlineUsers(t *testing.T) {
	client := CreateClient()
	userList, err := client.GetUserList()
	if err != nil {
		t.Error(err)
	}
	client.Debug()
	onlineUserList := make([]api.OnlineUser, len(*userList))
	for i, userInfo := range *userList {
		onlineUserList[i] = api.OnlineUser{
			UID: userInfo.UID,
			IP:  fmt.Sprintf("1.1.1.%d", i),
		}
	}
	// client.Debug()
	err = client.ReportNodeOnlineUsers(&onlineUserList)
	if err != nil {
		t.Error(err)
	}
}

func TestReportReportUserTraffic(t *testing.T) {
	client := CreateClient()
	userList, err := client.GetUserList()
	if err != nil {
		t.Error(err)
	}
	generalUserTraffic := make([]api.UserTraffic, len(*userList))
	for i, userInfo := range *userList {
		generalUserTraffic[i] = api.UserTraffic{
			UID:      userInfo.UID,
			Upload:   100,
			Download: 100,
		}
	}
	client.Debug()
	err = client.ReportUserTraffic(&generalUserTraffic)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNodeRule(t *testing.T) {
	client := CreateClient()
	client.Debug()
	ruleList, err := client.GetNodeRule()
	if err != nil {
		t.Error(err)
	}

	t.Log(ruleList)
}

func TestReportIllegal(t *testing.T) {
	client := CreateClient()

	detectResult := []api.DetectResult{
		{1, 1},
		{1, 2},
	}
	client.Debug()
	err := client.ReportIllegal(&detectResult)
	if err != nil {
		t.Error(err)
	}
}
