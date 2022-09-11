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
	EventId     int    `json:"EventId"`
	EventType   int    `json:"EventType"`
	RefreshType int    `json:"RefreshType"`
	Name        string `json:"Name"`
	EventDrop   int    `json:"EventDrop"`
	MapId       int    `json:"MapId"`
}

var (
	ConfigMapMap      map[int]*ConfigMap
	ConfigMapEventMap map[int]*ConfigMapEvent
)

func init() {
	ConfigMapMap = map[int]*ConfigMap{}
	err := utils.GetCsvUtilMgr().LoadCsv("Map", &ConfigMapMap)
	if err != nil {
		fmt.Println("初始化地图模块错误", err)
	}

	ConfigMapEventMap = map[int]*ConfigMapEvent{}
	err = utils.GetCsvUtilMgr().LoadCsv("MapEvent", &ConfigMapEventMap)
	if err != nil {
		fmt.Println("初始化蒙德地图错误", err)
	}
}

func GetMapName(mapId int) string {
	_, ok := ConfigMapMap[mapId]
	if !ok {
		return ""
	}
	return ConfigMapMap[mapId].MapName
}

func GetEventName(eventId int) string {
	_, ok := ConfigMapEventMap[eventId]
	if !ok {
		return ""
	}
	return ConfigMapEventMap[eventId].Name
}

func GetEventConfig(eventId int) *ConfigMapEvent {
	return ConfigMapEventMap[eventId]
}
