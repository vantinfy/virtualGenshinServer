package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type HomeItemInfo struct {
	HomeItemId  int   `json:"HomeItemId"`
	HomeItemNum int64 `json:"HomeItemNum"`
	KeyId       int64 `json:"KeyId"`
	Type        int   `json:"Type"`
}

type ModHome struct {
	HomeItemInfo map[int]*HomeItemInfo
}

// AddItem 添加家园物品
func (p *ModHome) AddItem(player *Player, itemId int, account int64) {
	// 判断获得的家园物品信息配置是否真实存在
	config := csvs.GetConfigHomeItem(itemId)
	if config == nil {
		fmt.Println(itemId, "家园物品配置不存在")
		return
	}

	_, ok := p.HomeItemInfo[itemId]
	if ok {
		p.HomeItemInfo[itemId].HomeItemNum += account
	} else {
		p.HomeItemInfo[itemId] = &HomeItemInfo{
			HomeItemId:  itemId,
			HomeItemNum: account,
			Type:        config.Type,
		}
	}

	fmt.Printf("获得家园物品: %vx%v, 当前数量:%v\n", csvs.GetItemName(itemId), account, p.HomeItemInfo[itemId].HomeItemNum)
}
