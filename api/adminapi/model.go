package adminapi

import (
	"encoding/json"
)

type ShadowsocksNodeInfo struct {
	ID          int   `json:"id"`
	SpeedLimit  uint64 	`json:"speedlimit"`
	ClientLimit int    `json:"clientlimit"`
	Port        uint32 `json:"port"`
	Setting struct {
		Method string `json:"method"` 
	} `json:"setting"`
}


type UserResponse struct {
	ID          int     `json:"id"`
	Passwd      string  `json:"password"`
	SpeedLimit  float64 `json:"nodespeedlimit"`
	DeviceLimit int     `json:"nodeconnector"`
}

type UserTraffic struct {
	UserId      int   `json:"userid"`
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