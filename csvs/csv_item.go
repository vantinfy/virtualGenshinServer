package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigItem struct {
	ItemId   int    `json:"ItemId"`
	SortType int    `json:"SortType"`
	ItemName string `json:"ItemName"`
}

var ConfigItemMap map[int]*ConfigItem

const (
	ITEM_TYPE_NORMAL        = 1  // 其它
	ITEM_TYPE_ROLE          = 2  // 角色
	ITEM_TYPE_AVATAR        = 3  // 头像
	ITEM_TYPE_CARD          = 4  // 名片
	ITEM_TYPE_CONSTELLATION = 5  // 命星
	ITEM_TYPE_WEAPON        = 6  // 武器
	ITEM_TYPE_RELICS        = 7  // 圣遗物
	ITEM_TYPE_COOKBOOK      = 8  // 食谱
	ITEM_TYPE_COOK          = 9  // 烹饪技能
	ITEM_TYPE_FOOD          = 10 // 食物
	ITEM_TYPE_HOME_Item     = 11 // 家园物品
)

func init() {
	ConfigItemMap = make(map[int]*ConfigItem)
	err := utils.GetCsvUtilMgr().LoadCsv("Item", ConfigItemMap)
	if err != nil {
		fmt.Println("加载物品配置表失败", err)
	}
}

// GetItemConfig 根据物品id查询物品信息（类型及物品名称）
func GetItemConfig(itemId int) *ConfigItem {
	return ConfigItemMap[itemId]
}

// GetItemName 根据物品id查询物品名
func GetItemName(itemId int) string {
	return ConfigItemMap[itemId].ItemName
}
