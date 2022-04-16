package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigWeapon struct {
	WeaponId   int `json:"WeaponId"`
	WeaponType int `json:"WeaponType"`
	Star       int `json:"Star"`
}

var (
	ConfigWeaponMap map[int]*ConfigWeapon
)

func init() {
	ConfigWeaponMap = make(map[int]*ConfigWeapon)
	err := utils.GetCsvUtilMgr().LoadCsv("Weapon", ConfigWeaponMap)
	if err != nil {
		fmt.Println("加载武器配置失败")
		return
	}
}

func GetConfigWeapon(weaponId int) *ConfigWeapon {
	return ConfigWeaponMap[weaponId]
}
