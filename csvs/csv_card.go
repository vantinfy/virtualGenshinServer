package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigCard struct {
	CardId       int `json:"CardId"`
	Friendliness int `json:"Friendliness"`
	Check        int `json:"Check"`
}

var (
	ConfigCardMap         map[int]*ConfigCard
	ConfigCardMapByRoleId map[int]*ConfigCard
)

func init() {
	ConfigCardMap = make(map[int]*ConfigCard)
	ConfigCardMapByRoleId = make(map[int]*ConfigCard)

	// load base csv
	err := utils.GetCsvUtilMgr().LoadCsv("Card", ConfigCardMap)
	if err != nil {
		fmt.Println("加载名片配置失败")
		return
	}
	for _, card := range ConfigCardMap {
		ConfigCardMapByRoleId[card.Check] = card
	}
}

// GetCardConfig 获得名片配置信息
func GetCardConfig(carId int) *ConfigCard {
	return ConfigCardMap[carId]
}

// GetCardConfigByRoleId 通过角色id获取对应名片配置
func GetCardConfigByRoleId(roleId int) *ConfigCard {
	return ConfigCardMapByRoleId[roleId]
}
