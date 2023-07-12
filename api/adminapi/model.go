package adminapi

import "encoding/json"

type Response struct {
	Status  string          `json:"type"`
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"result"`
	Message string          `json:"message"`
}

type V2rayNodeInfo struct {
	ID            int    `json:"id"`
	IsUDP         bool   `json:"is_udp"`
	SpeedLimit    uint64 `json:"speed_limit"`
	ClientLimit   int    `json:"client_limit"`
	PushPort      int    `json:"push_port"`
	Secret        string `json:"secret"`
	Key           string `json:"key"`
	Cert          string `json:"pem"`
	V2License     string `json:"v2_license"`
	V2AlterID     uint16 `json:"v2_alter_id"`
	V2Port        uint32 `json:"v2_port"`
	V2Method      string `json:"v2_method"`
	V2Net         string `json:"v2_net"`
	V2Type        string `json:"v2_type"`
	V2Host        string `json:"v2_host"`
	V2Path        string `json:"v2_path"`
	V2TLS         bool   `json:"v2_tls"`
	V2Cdn         bool   `json:"v2_cdn"`
	V2TLSProvider string `json:"v2_tls_provider"`
	RedirectUrl   string `json:"redirect_url"`
}

type ShadowsocksNodeInfo struct {
	ID          int    `json:"id"`
	SpeedLimit  uint64 `json:"speedlimit"`
	ClientLimit int    `json:"clientlimit"`
	Port        uint32 `json:"port"`
	Setting     struct {
		Method string `json:"method"`
	} `json:"setting"`
}

type TrojanNodeInfo struct {
	ID          int    `json:"id"`
	IsUDP       bool   `json:"is_udp"`
	SpeedLimit  uint64 `json:"speed_limit"`
	ClientLimit int    `json:"client_limit"`
	PushPort    int    `json:"push_port"`
	TrojanPort  uint32 `json:"trojan_port"`
}

// NodeStatus Node status report
type NodeStatus struct {
	CPU    string `json:"cpu"`
	Mem    string `json:"mem"`
	Net    string `json:"net"`
	Disk   string `json:"disk"`
	Uptime int    `json:"uptime"`
}

type NodeOnline struct {
	UID int    `json:"UserNo"`
	IP  string `json:"ip"`
}

type UseUserTraffic struct {
	NodeId int        `json:"nodeid"`
	Data   []UsedData `json:"data"`
}

type UsedData struct {
	Upload   int64 `json:"upload"`
	Download int64 `json:"download"`
	UserId   int   `json:"userid"`
}

type VMessUser struct {
	UID        int    `json:"uid"`
	VmessUID   string `json:"vmess_uid"`
	SpeedLimit uint64 `json:"speed_limit"`
}

type TrojanUser struct {
	UID        int    `json:"uid"`
	Password   string `json:"password"`
	SpeedLimit uint64 `json:"speed_limit"`
}

type SSUser struct {
	UID        int    `json:"id"`
	Password   string `json:"password"`
	SpeedLimit uint64 `json:"nodespeedlimit"`
}

type UserTraffic struct {
	UID      int   `json:"uid"`
	Upload   int64 `json:"upload"`
	Download int64 `json:"download"`
}

type NodeRule struct {
	Mode  string         `json:"mode"`
	Rules []NodeRuleItem `json:"rules"`
}

type NodeRuleItem struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

type IllegalReport struct {
	UID    int    `json:"uid"`
	RuleID int    `json:"rule_id"`
	Reason string `json:"reason"`
}

type Certificate struct {
	Key string `json:"key"`
	Pem string `json:"pem"`
}
