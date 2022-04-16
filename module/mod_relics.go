package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type Relics struct {
	RelicsId   int    `json:"RelicsId"`
	Type       int    `json:"Type"`
	Pos        int    `json:"Pos"`
	Star       int    `json:"Star"`
	RelicsName string `json:"RelicsName"`
	KeyId      int64  `json:"Key"`
	// 圣遗物词条，攻击防御生命（百分比/固定值）充能精通双爆属性伤（7+物理）治疗加成
}

type ModRelics struct {
	RelicsInfo map[int64]*Relics
	MaxKey     int64 // 记录历史得到的圣遗物总量，不会随着武器的销毁减少
}

// AddItem 添加圣遗物
func (p *ModRelics) AddItem(relicsId int, account int64) {
	if len(p.RelicsInfo)+int(account) > csvs.RELICS_BAG_MAX_ACCOUNT {
		fmt.Println("圣遗物背包容量即将超过上限")
		return
	}

	// 是否存在对应的圣遗物配置
	config := csvs.GetConfigRelics(relicsId)
	if config == nil {
		fmt.Println("非法圣遗物")
		return
	}

	relicsName := csvs.GetItemName(relicsId)

	for i := int64(0); i < account; i++ {
		// 新获得圣遗物id递增
		p.MaxKey++
		key := p.MaxKey
		p.RelicsInfo[key] = &Relics{
			RelicsId:   relicsId,
			Type:       config.Type,
			Pos:        config.Pos,
			Star:       config.Star,
			KeyId:      key,
			RelicsName: relicsName,
		}
	}

	fmt.Println("获得新圣遗物", relicsName, account, "件")
}
