package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type ModBag struct {
	BagInfo map[int]*ItemInfo
}

type ItemInfo struct {
	ItemId  int
	ItemNum int64
}

// AddItem 添加物品
func (b *ModBag) AddItem(player *Player, itemId int, account int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	switch itemConfig.SortType {
	//case csvs.ITEM_TYPE_NORMAL: // 杂物
	//	b.AddItemToBag(itemId, account)
	case csvs.ITEM_TYPE_ROLE: // 角色
		player.ModRole.AddItem(player, itemId, account)
	case csvs.ITEM_TYPE_AVATAR: // 角色头像
		player.ModAvatar.AddItem(itemId)
	case csvs.ITEM_TYPE_CARD: // 名片
		player.ModCard.AddItem(itemId, 10)
	case csvs.ITEM_TYPE_WEAPON: // 武器
		player.ModWeapon.AddItem(itemId, account)
	case csvs.ITEM_TYPE_RELICS: // 圣遗物
		player.ModRelics.AddItem(itemId, account)
	//case csvs.ITEM_TYPE_COOKBOOK: // 食谱
	//	fmt.Println(itemConfig.ItemName, "物品增加了", account)
	case csvs.ITEM_TYPE_COOK: // 烹饪
		player.ModCook.AddItem(itemId)
	case csvs.ITEM_TYPE_HOME_Item: // 家园物品相关
		player.ModHome.AddItem(player, itemId, account)
	default:
		b.AddItemToBag(itemId, account)
		// 同normal
		//fmt.Println("尚有未归类物品没有实现添加功能", itemId)
	}
}

// RemoveItem 移除物品
func (b *ModBag) RemoveItem(player *Player, itemId int, account int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	switch itemConfig.SortType {
	case csvs.ITEM_TYPE_WEAPON: // 武器
	case csvs.ITEM_TYPE_RELICS: // 圣遗物
	default:
		b.RemoveItemFromBag(itemId, account)
	}
}

// UseItem 使用物品
func (b *ModBag) UseItem(player *Player, itemId int, account int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	if b.Has(itemId) < account {
		fmt.Println("背包物品剩余数量不足", b.Has(itemId), account)
		return
	}

	switch itemConfig.SortType {
	case csvs.ITEM_TYPE_COOKBOOK: // 食谱
		b.UseCookBook(player, itemId, account)
	case csvs.ITEM_TYPE_FOOD: // 食物
		fmt.Println("食用了", csvs.GetItemName(itemId), account, "个")
	default:
		fmt.Println("该物品无法使用", itemId, csvs.GetItemName(itemId))
	}
}

// Has 查询背包有多少itemId的物品，如果没有，返回0
func (b *ModBag) Has(itemId int) int64 {
	item, ok := b.BagInfo[itemId]
	if ok {
		return item.ItemNum
	}
	return 0
}

// AddItemToBag 添加物品到背包
func (b *ModBag) AddItemToBag(itemId int, account int64) {
	config := csvs.GetItemConfig(itemId)
	if config == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	item, ok := b.BagInfo[itemId]
	if ok {
		// 如果物品已经拥有过 todo 有些物品有数量上限，比如浓缩树脂最多持有5个，同一食谱只能有一个
		if config.SortType == csvs.ITEM_TYPE_COOKBOOK {
			fmt.Println("该食谱已经存在", config.ItemName)
			return
		}
		item.ItemNum += account
	} else {
		if config.SortType == csvs.ITEM_TYPE_COOKBOOK {
			fmt.Println("该食谱最多拥有1个", config.ItemName)
			account = 1
		}
		b.BagInfo[itemId] = &ItemInfo{
			ItemId:  itemId,
			ItemNum: account,
		}
	}
	fmt.Printf("获得了%v个%v	现有%v个\n", account, config.ItemName, b.BagInfo[itemId].ItemNum)
}

// RemoveItemFromBagGM 从背包中移除物品（不校验数量）
func (b *ModBag) RemoveItemFromBagGM(itemId int, account int64) {
	config := csvs.GetItemConfig(itemId)
	if config == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	item, ok := b.BagInfo[itemId]
	if ok {
		item.ItemNum -= account
	} else {
		b.BagInfo[itemId] = &ItemInfo{
			ItemId:  itemId,
			ItemNum: -account,
		}
	}
	fmt.Printf("使用了%v个%v物品，剩余%v\n", account, itemId, b.BagInfo[itemId].ItemNum)
}

// RemoveItemFromBag 从背包中移除物品
func (b *ModBag) RemoveItemFromBag(itemId int, account int64) {
	config := csvs.GetItemConfig(itemId)
	if config == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	if b.Has(itemId) < account {
		fmt.Println(itemId, "可用数量不足")
		return
	}

	b.BagInfo[itemId].ItemNum -= account
	fmt.Printf("使用了%v个%v，剩余%v\n", account, csvs.GetItemName(itemId), b.BagInfo[itemId].ItemNum)
}

// UseCookBook 使用食谱学习烹饪技能
func (b *ModBag) UseCookBook(player *Player, cookBookId int, account int64) {
	config := csvs.GetConfigCookBook(cookBookId)
	if config == nil {
		fmt.Println("不存在对应的食谱")
		return
	}

	//b.Bag[cookBookId].ItemNum--
	b.RemoveItem(player, cookBookId, account)
	//player.ModCook.AddItem(config.Reward)
	b.AddItem(player, config.Reward, account)
	fmt.Println("现存物品", csvs.GetItemName(cookBookId), b.Has(cookBookId), "个")
}
