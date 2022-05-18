package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type ModPool struct {
	UpPoolInfo *PoolInfo
}

type PoolInfo struct {
	PoolId    int
	DrawTimes int // 该池子抽取了多少次还未出5星，出货后重置
}

func (p *ModPool) UpPoolDraw(times int) {
	cards := map[int]int{}
	for i := 0; i < times; i++ {
		dropGroup := csvs.ConfigDropMap[1000]
		if dropGroup == nil {
			return
		}
		// 概率修正
		if p.UpPoolInfo.DrawTimes > csvs.FIVE_STAR_TIMES_LIMIT {
			alterGroup := &csvs.DropGroup{}
			alterGroup.DropId = dropGroup.DropId
			alterGroup.WeightAll = dropGroup.WeightAll

			// 本次抽卡增加概率
			increase := (p.UpPoolInfo.DrawTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_TRIGGER_INCREASE
			for _, drop := range dropGroup.ConfigDrops {
				tmpDrop := &csvs.ConfigDrop{}
				tmpDrop.DropId = drop.DropId
				tmpDrop.IsEnd = drop.IsEnd
				tmpDrop.Result = drop.Result
				// 5星概率增加，3星概率下降，4星不变
				if drop.Result == 10001 {
					tmpDrop.Weight = drop.Weight + increase
				} else if drop.Result == 10003 {
					tmpDrop.Weight = drop.Weight - increase
				} else {
					tmpDrop.Weight = drop.Weight
				}
				// 使用修改过概率的配置
				alterGroup.ConfigDrops = append(alterGroup.ConfigDrops, tmpDrop)
			}
			dropGroup = alterGroup
		}

		card := csvs.GetRandDropRecursion(dropGroup)
		if card != nil {
			role := csvs.GetRoleConfig(card.Result)
			if role != nil && role.Star == 5 {
				p.UpPoolInfo.DrawTimes = 0 // 出货5星后重置次数
			} else {
				p.UpPoolInfo.DrawTimes++
			}
			cards[card.Result]++
		}
	}

	for k, v := range cards {
		fmt.Printf("%s\t次数:%d\n", csvs.GetItemName(k), v)
	}
}
