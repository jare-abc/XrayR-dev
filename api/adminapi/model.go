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
			ServerKey   string           `json:"method"`
			Network     string           `json:"network"`
			Headers     *json.RawMessage `json:"headers"`
			ServiceName string           `json:"serviceName"`
			Header      *json.RawMessage `json:"header"`
			Tls         int              `json:"tls"`
			Host        string           `json:"host"`
			ServerName  string           `json:"server_name"`
		} `json:"setting"`
		Routes []route `json:"routes"`
	} `json:"result"`
	RExtras string `json:"extras"`
	RTime   string `json:"time"`
}

/* type serverConfig struct {
	shadowsocks
	v2ray
	trojan

	ServerPort int `json:"port"`
	BaseConfig struct {
		PushInterval int `json:"push_interval"`
		PullInterval int `json:"pull_interval"`
	} `json:"base_config"`
	Routes []route `json:"routes"`
} */

type route struct {
	Id          int      `json:"id"`
	Match       []string `json:"match"`
	Action      string   `json:"action"`
	ActionValue string   `json:"action_value"`
}

type user struct {
	Id         int    `json:"id"`
	Uuid       string `json:"uuid"`
	SpeedLimit int    `json:"speed_limit"`
}
