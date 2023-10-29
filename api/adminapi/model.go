package adminapi

import (
	"encoding/json"
)

type serverConfig struct {
	RCode      int    `json:"code"`
	RType      string `json:"type"`
	RMessage   string `json:"message"`
	ResultData struct {
		NodeId      int     `json:"nodeid"`
		NodeType    string  `json:"nodetype"`
		ServerPort  int     `json:"port"`
		SpeedLimit  float64 `json:"speedlimit"`
		ClientLimit int     `json:"clientlimit"`
		NodeSetting struct {
			Cipher      string           `json:"cipher"`
			Obfs        string           `json:"obfs"`
			Path        string           `json:"path"`
			ServerKey   string           `json:"serverkey"`
			Network     string           `json:"network"`
			Headers     *json.RawMessage `json:"headers"`
			ServiceName string           `json:"serviceName"`
			Header      *json.RawMessage `json:"header"`
			Tls         int              `json:"tls"`
			Host        string           `json:"host"`
			Method      string           `json:"method"`
		} `json:"setting"`
		Routes []route `json:"routes"`
	} `json:"result"`
	RExtras string `json:"extras"`
	RTime   string `json:"time"`
}

type route struct {
	Id          int      `json:"id"`
	Match       []string `json:"match"`
	Action      string   `json:"action"`
	ActionValue string   `json:"action_value"`
}

type user struct {
	Id          int    `json:"id"`
	PassWord    string `json:"password"`
	SpeedLimit  uint64 `json:"nodespeedlimit"`
	DeviceLimit int    `json:"nodeconnector"`
}

type OnlineUser struct {
	UID int    `json:"UserNo"`
	IP  string `json:"IP"`
}
type PostData struct {
	Users   interface{} `json:"users"`
	Onlines interface{} `json:"onlines"`
}

type UserTraffic struct {
	UID      int   `json:"userno"`
	Upload   int64 `json:"upload"`
	Download int64 `json:"download"`
}

type NodeStatus struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	Net    string `json:"net"`
	Disk   string `json:"disk"`
	Uptime int    `json:"uptime"`
}
