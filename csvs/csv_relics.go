package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigRelics struct {
	RelicsId int `json:"RelicsId"`
	Type     int `json:"Type"`
	Pos      int `json:"Pos"`
	Star     int `json:"Star"`
}

var (
	ConfigRelicsMap map[int]*ConfigRelics
)

func init() {
	ConfigRelicsMap = make(map[int]*ConfigRelics)
	err := utils.GetCsvUtilMgr().LoadCsv("Relics", ConfigRelicsMap)
	if err != nil {
		fmt.Println("加载圣遗物配置失败")
		return
	}
}

func GetConfigRelics(relicsId int) *ConfigRelics {
	return ConfigRelicsMap[relicsId]
}
