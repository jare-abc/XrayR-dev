package adminapi_test

import (
	"fmt"
	"testing"

	"github.com/jare-abc/XrayR-dev/api"
	"github.com/jare-abc/XrayR-dev/api/adminapi"
)

func CreateClient() api.API {
	apiConfig := &api.Config{
		APIHost:  "http://localhost:5005",
		Key:      "qwertyuiopasdfghjkl",
		NodeID:   1,
		NodeType: "Shadowsocks",
	}
	client := adminapi.New(apiConfig)
	return client
}

func TestGetV2rayNodeInfo(t *testing.T) {
	client := CreateClient()
	nodeInfo, err := client.GetNodeInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodeInfo)
}

func TestGetSSNodeInfo(t *testing.T) {
	apiConfig := &api.Config{
		APIHost:  "http://127.0.0.1:5005",
		Key:      "qwertyuiopasdfghjkl",
		NodeID:   1,
		NodeType: "Shadowsocks",
	}
	client := adminapi.New(apiConfig)
	nodeInfo, err := client.GetNodeInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodeInfo)
}

func TestGetTrojanNodeInfo(t *testing.T) {
	apiConfig := &api.Config{
		APIHost:  "http://127.0.0.1:5005",
		Key:      "qwertyuiopasdfghjkl",
		NodeID:   1,
		NodeType: "Trojan",
	}
	client := adminapi.New(apiConfig)
	nodeInfo, err := client.GetNodeInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodeInfo)
}

func TestReportReportNodeOnlineUsers(t *testing.T) {
	apiConfig := &api.Config{
		APIHost:  "http://127.0.0.1:5005",
		Key:      "qwertyuiopasdfghjkl",
		NodeID:   1,
		NodeType: "Shadowsocks",
	}
	client := adminapi.New(apiConfig)
	userList, err := client.GetUserList()
	if err != nil {
		t.Error(err)
	}

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

func TestGetUserList(t *testing.T) {
	client := CreateClient()

	userList, err := client.GetUserList()
	if err != nil {
		t.Error(err)
	}

	t.Log(userList)
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
			Upload:   114514,
			Download: 114514,
		}
	}
	// client.Debug()
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
