package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigMap struct {
	MapId   int    `json:"MapId"`
	MapName string `json:"MapName"`
}

type ConfigMapEvent struct {
	EventId   int `json:"EventId"`
	EventType int `json:"EventType"`
	EventItem int `json:"EventItem"`
}

var (
	ConfigMapMap       map[int]*ConfigMap
	ConfigMapMondstadt map[int]*ConfigMapEvent
)

func init() {
	ConfigMapMap = map[int]*ConfigMap{}
	err := utils.GetCsvUtilMgr().LoadCsv("Map", &ConfigMapMap)
	if err != nil {
		fmt.Println("初始化地图模块错误", err)
	}

	ConfigMapMondstadt = map[int]*ConfigMapEvent{}
	err = utils.GetCsvUtilMgr().LoadCsv("MapMondstadt", &ConfigMapMondstadt)
	if err != nil {
		fmt.Println("初始化蒙德地图错误", err)
	}
}
