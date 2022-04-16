package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigCook struct {
	CookId int `json:"CookId"`
	Star   int `json:"Star"`
}

var (
	ConfigCookMap map[int]*ConfigCook
)

func init() {
	ConfigCookMap = make(map[int]*ConfigCook)
	err := utils.GetCsvUtilMgr().LoadCsv("Cook", ConfigCookMap)
	if err != nil {
		fmt.Println("加载食谱配置失败")
		return
	}
}

func GetConfigCook(cookId int) *ConfigCook {
	return ConfigCookMap[cookId]
}
