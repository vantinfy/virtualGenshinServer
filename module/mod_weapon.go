package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type Weapon struct {
	WeaponId   int    `json:"WeaponId"`
	KeyId      int    `json:"Key"`
	WeaponName string `json:"WeaponName"`
}

type ModWeapon struct {
	WeaponInfo map[int]*Weapon
	MaxKey     int // 记录历史得到的武器总量，不会随着武器的销毁减少
}

// AddItem 添加武器
func (p *ModWeapon) AddItem(weaponId int, account int64) {
	if len(p.WeaponInfo)+int(account) > csvs.WEAPON_BAG_MAX_ACCOUNT {
		fmt.Println("武器背包容量即将超过上限")
		return
	}

	// 是否存在对应的武器配置
	config := csvs.GetConfigWeapon(weaponId)
	if config == nil {
		fmt.Println("非法武器")
		return
	}

	weaponName := csvs.GetItemName(weaponId)

	for i := 0; i < int(account); i++ {
		// 新获得武器id递增
		p.MaxKey++
		key := p.MaxKey
		p.WeaponInfo[key] = &Weapon{
			WeaponId:   weaponId,
			KeyId:      key,
			WeaponName: weaponName,
		}
	}

	fmt.Println("获得新武器", weaponName, account, "件")
}
