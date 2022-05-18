package csvs

import (
	"fmt"
	"math/rand"
	"time"
)

type DropGroup struct {
	DropId      int
	WeightAll   int
	ConfigDrops []*ConfigDrop
}

var ConfigDropMap map[int]*DropGroup

func CheckLoadCsv() {
	// 设置好随机种子
	rand.Seed(time.Now().UnixNano())
	// 是否加载完成，二次处理
	makeDropGroupMap()
	fmt.Println("csv配置加载完成")
}

func makeDropGroupMap() {
	ConfigDropMap = make(map[int]*DropGroup)
	for _, drop := range ConfigDropSlice {
		//
		dropGroup, ok := ConfigDropMap[drop.DropId]
		if !ok {
			dropGroup = new(DropGroup)
			dropGroup.DropId = drop.DropId
			ConfigDropMap[drop.DropId] = dropGroup
		}
		dropGroup.WeightAll += drop.Weight
		dropGroup.ConfigDrops = append(dropGroup.ConfigDrops, drop)
	}
	// 测试一百连抽
	//TimesDraw(100)
	return
}

func TimesDraw(times int) {
	dropGroup := ConfigDropMap[1000]
	if dropGroup == nil {
		return
	}
	cnt := 0
	for {
		// 使用递归算法后config如果不为nil，则config.IsEnd==LOGIC_TRUE
		config := GetRandDropRecursion(dropGroup)
		if config == nil {
			return
		}
		//if config.IsEnd == LOGIC_TRUE
		fmt.Println(GetItemName(config.Result))
		cnt++
		if cnt > times {
			break
		}
	}
}

// GetRandDropRecursion 递归计算掉落物品（eg.抽卡）
func GetRandDropRecursion(dropGroup *DropGroup) *ConfigDrop {
	randNum := rand.Intn(dropGroup.WeightAll)
	randNow := 0
	for _, drop := range dropGroup.ConfigDrops {
		randNow += drop.Weight
		if randNum < randNow {
			if drop.IsEnd == LOGIC_TRUE {
				// 递归出口，当获取的掉落配置isEnd=1，表示已经摇到了具体某一个item（角色/武器）而不是另一个掉落配置
				return drop
			}
			dropGroup = ConfigDropMap[drop.Result]
			if dropGroup == nil {
				return nil
			}
			return GetRandDropRecursion(dropGroup)
		}
	}
	return nil
}

func RandDropTest() {
	dropGroup := ConfigDropMap[1000]
	if dropGroup == nil {
		return
	}
	for {
		config := GetRandDrop(dropGroup)
		if config == nil {
			return
		}
		if config.IsEnd == LOGIC_TRUE {
			fmt.Println(GetItemName(config.Result))
			break
		}
		dropGroup = ConfigDropMap[config.Result]
		if dropGroup == nil {
			return
		}
	}
}

func GetRandDrop(dropGroup *DropGroup) *ConfigDrop {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(dropGroup.WeightAll)
	randNow := 0
	for _, drop := range dropGroup.ConfigDrops {
		randNow += drop.Weight
		if randNum < randNow {
			return drop
		}
	}
	return nil
}
