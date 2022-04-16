package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type Cook struct {
	CookId   int    `json:"CookId"`
	Star     int    `json:"Star"`
	CookName string `json:"CookName"`
	// 熟练度
}

type ModCook struct {
	CookInfo map[int]*Cook
}

// AddItem 学习烹饪技能
func (p *ModCook) AddItem(cookId int) {

	if _, ok := p.CookInfo[cookId]; ok {
		fmt.Println("已经学会该烹饪技能")
		return
	}

	// 是否存在对应的食谱配置
	config := csvs.GetConfigCook(cookId)
	if config == nil {
		fmt.Println("非法食谱")
		return
	}

	cookName := csvs.GetItemName(cookId)

	p.CookInfo[cookId] = &Cook{
		CookId:   cookId,
		Star:     config.Star,
		CookName: cookName,
	}

	fmt.Println("习得新烹饪技能", cookName)
}
