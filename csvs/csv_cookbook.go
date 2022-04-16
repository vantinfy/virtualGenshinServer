package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigCookBook struct {
	CookBookId int `json:"CookBookId"`
	Reward     int `json:"Reward"`
}

var (
	ConfigCookBookMap map[int]*ConfigCookBook
)

func init() {
	ConfigCookBookMap = make(map[int]*ConfigCookBook)
	err := utils.GetCsvUtilMgr().LoadCsv("CookBook", ConfigCookBookMap)
	if err != nil {
		fmt.Println("加载食谱配置失败")
		return
	}
}

func GetConfigCookBook(cookbookId int) *ConfigCookBook {
	return ConfigCookBookMap[cookbookId]
}
