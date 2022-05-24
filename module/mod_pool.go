package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type ModPool struct {
	UpPoolInfo *PoolInfo
}

type PoolInfo struct {
	PoolId        int
	FiveDrawTimes int // 该池子抽取了多少次还未出5星，出货后重置
	FourDrawTimes int // 该池子抽取了多少次还未出4星，出货后重置
}

func (p *ModPool) UpPoolDraw(times int) {
	cards := map[int]int{}
	cnt5 := 0 // 5星统计
	cnt4 := 0 // 4星统计
	for i := 0; i < times; i++ {
		dropGroup := csvs.ConfigDropMap[1000]
		if dropGroup == nil {
			return
		}
		// 概率修正
		if p.UpPoolInfo.FiveDrawTimes > csvs.FIVE_STAR_TIMES_LIMIT || p.UpPoolInfo.FourDrawTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			alterGroup := &csvs.DropGroup{}
			alterGroup.DropId = dropGroup.DropId
			alterGroup.WeightAll = dropGroup.WeightAll

			// 本次抽卡增加概率
			increaseFiveStar := (p.UpPoolInfo.FiveDrawTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_TRIGGER_INCREASE
			if increaseFiveStar < 0 {
				increaseFiveStar = 0
			}
			increaseFourStar := (p.UpPoolInfo.FourDrawTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_TRIGGER_INCREASE
			if increaseFourStar < 0 {
				increaseFourStar = 0
			}
			for _, drop := range dropGroup.ConfigDrops {
				tmpDrop := &csvs.ConfigDrop{}
				tmpDrop.DropId = drop.DropId
				tmpDrop.IsEnd = drop.IsEnd
				tmpDrop.Result = drop.Result
				// 4、5星概率增加，3星概率下降
				if drop.Result == 10001 {
					tmpDrop.Weight = drop.Weight + increaseFiveStar
				} else if drop.Result == 10002 {
					tmpDrop.Weight = drop.Weight + increaseFourStar
				} else if drop.Result == 10003 {
					tmpDrop.Weight = drop.Weight - increaseFiveStar - increaseFourStar
				}
				// 使用修改过概率的配置
				alterGroup.ConfigDrops = append(alterGroup.ConfigDrops, tmpDrop)
			}
			dropGroup = alterGroup
		}

		card := csvs.GetRandDropRecursion(dropGroup)
		if card != nil {
			// 考虑了保底歪了武器
			role := csvs.GetRoleConfig(card.Result)
			if role != nil { // 出了角色
				if role.Star == 5 { // 出货5星角色后重置次数
					p.UpPoolInfo.FiveDrawTimes = 0
					p.UpPoolInfo.FourDrawTimes++
					cnt5++
				} else if role.Star == 4 { // 出货4星角色后重置次数
					p.UpPoolInfo.FourDrawTimes = 0
					p.UpPoolInfo.FiveDrawTimes++
					cnt4++
				}
			} else { // 出了武器[因为角色只有4、5星，role==nil则掉落一定为武器（3、4、5星）]
				weapon := csvs.GetConfigWeapon(card.Result)
				if weapon == nil {
					fmt.Println("抽卡结果掉落：武器配置信息错误")
					return
				} else if weapon.Star == 5 { // 五星武器
					p.UpPoolInfo.FiveDrawTimes = 0
					p.UpPoolInfo.FourDrawTimes++
					cnt5++
				} else if weapon.Star == 4 { // 四星武器
					p.UpPoolInfo.FourDrawTimes = 0
					p.UpPoolInfo.FiveDrawTimes++
					cnt4++
				} else { // 三星武器
					p.UpPoolInfo.FiveDrawTimes++
					p.UpPoolInfo.FourDrawTimes++
				}
			}
			cards[card.Result]++
		}
	}

	for k, v := range cards {
		fmt.Printf("%s\t次数:%d\n", csvs.GetItemName(k), v)
	}
	fmt.Printf("累计5星:%d，累计四星:%d\n", cnt5, cnt4)
}
