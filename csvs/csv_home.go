package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigHomeItem struct {
	HomeItemId int `json:"HomeItemId"`
	Type       int `json:"Type"`
}

var (
	ConfigHomeMap map[int]*ConfigHomeItem
)

func init() {
	ConfigHomeMap = make(map[int]*ConfigHomeItem)
	err := utils.GetCsvUtilMgr().LoadCsv("Home", ConfigHomeMap)
	if err != nil {
		fmt.Println("加载家园配置失败")
		return
	}
}

func GetConfigHomeItem(homeItemId int) *ConfigHomeItem {
	return ConfigHomeMap[homeItemId]
}
