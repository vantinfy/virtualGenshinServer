package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type CardInfo struct {
	CardId   int    `json:"CardId"`
	CardName string `json:"CardName"`
}

type ModCard struct {
	CardInfo map[int]*CardInfo
}

// HasCard 查询是否拥有cardId对应的名片
func (p *ModCard) HasCard(cardId int) bool {
	_, ok := p.CardInfo[cardId]
	return ok
}

// AddItem 添加名片
func (p *ModCard) AddItem(cardId int, friendliness int) {
	if p.HasCard(cardId) {
		fmt.Println("名片已经拥有", cardId)
		return
	}

	// 判断获得的名片是否真实存在
	config := csvs.GetCardConfig(cardId)
	if config == nil {
		fmt.Println(cardId, "名片不存在")
		return
	}

	if friendliness < config.Friendliness {
		fmt.Println("好感度不足", cardId, "名片获取失败")
		return
	}

	p.CardInfo[cardId] = &CardInfo{
		CardId:   cardId,
		CardName: csvs.GetItemName(cardId),
	}
	fmt.Println("获得新名片", p.CardInfo[cardId].CardName)
}

// CheckGetCard 获得角色后连带检查好感度
func (p *ModCard) CheckGetCard(roleId int, friendliness int) {
	configCard := csvs.GetCardConfigByRoleId(roleId)
	if configCard == nil {
		fmt.Println("名片不存在")
		return
	}
	// 增加名片
	p.AddItem(configCard.CardId, friendliness)
}
