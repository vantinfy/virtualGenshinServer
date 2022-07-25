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
	FiveDrawTimes int        // 该池子抽取了多少次还未出5星，出货后重置
	FourDrawTimes int        // 该池子抽取了多少次还未出4星，出货后重置
	IsUpRole      int        // 保底五星是否为限定角色
	Statistics    *Statistic // 抽卡数据统计
}

// Statistic 抽卡数据统计
type Statistic struct {
	// 记录具体启用下面的哪些统计
	Switch map[string]int

	NGold        map[int]int // 十连n黄 map
	WhichGold    map[int]int // 第多少发出金
	PurpleCnt    int         // 4星统计
	GoldCnt      int         // 五星统计
	HistoryCnt   int         // 历史抽卡数统计
	GoldEveryTen int         // 十连内金色统计
}

type Options func(statistics *Statistic)

// WithNGold 启用[十连N金]
func WithNGold() Options {
	return func(statistics *Statistic) {
		// 这些方法在initStatistic中执行 只要在switch map初始化完成后是没有空指针问题的
		statistics.Switch["NGold"] = 1
	}
}

// WithWhichGold 启用[第多少发出金统计]
func WithWhichGold() Options {
	return func(statistics *Statistic) {
		statistics.Switch["WhichGold"] = 1
	}
}

// WithDropCnt 启用[4 5星总数掉落统计]
func WithDropCnt() Options {
	return func(statistics *Statistic) {
		statistics.Switch["PurpleCnt"] = 1
		statistics.Switch["GoldCnt"] = 1
	}
}

// WithHistoryCnt 启用[历史抽卡数统计]
func WithHistoryCnt() Options {
	return func(statistics *Statistic) {
		statistics.Switch["HistoryCnt"] = 1
	}
}

func (p *ModPool) InitStatistic(opts ...Options) map[string]int {
	// 初始化
	if p.UpPoolInfo.Statistics.Switch == nil {
		p.UpPoolInfo.Statistics.Switch = make(map[string]int)
	}
	if p.UpPoolInfo.Statistics.WhichGold == nil {
		p.UpPoolInfo.Statistics.WhichGold = make(map[int]int)
	}
	if p.UpPoolInfo.Statistics.NGold == nil {
		p.UpPoolInfo.Statistics.NGold = make(map[int]int)
	}
	// 启用哪些统计
	for _, opt := range opts {
		opt(p.UpPoolInfo.Statistics)
	}
	return p.UpPoolInfo.Statistics.Switch
}

func (p *ModPool) PrintStatistic(switches map[string]int) {
	if value, ok := switches["NGold"]; ok && value == 1 {
		fmt.Println("[十连N金]: ", p.UpPoolInfo.Statistics.NGold)
	}
	if value, ok := switches["WhichGold"]; ok && value == 1 {
		fmt.Println("[第N发出金]: ", p.UpPoolInfo.Statistics.WhichGold)
	}
	if value, ok := switches["GoldCnt"]; ok && value == 1 {
		fmt.Println("[5星掉落总计]: ", p.UpPoolInfo.Statistics.GoldCnt)
	}
	if value, ok := switches["PurpleCnt"]; ok && value == 1 {
		fmt.Println("[4星掉落总计]: ", p.UpPoolInfo.Statistics.PurpleCnt)
	}
	if value, ok := switches["HistoryCnt"]; ok && value == 1 {
		fmt.Println("[历史抽卡数总计]: ", p.UpPoolInfo.Statistics.HistoryCnt)
	}
}

func (p *ModPool) UpPoolDraw(player *Player, times int, opts ...Options) {
	cards := map[int]int{}
	switches := p.InitStatistic(opts...) // 开关 记录统计哪些信息
	for i := 0; i < times; i++ {
		p.UpPoolInfo.Statistics.HistoryCnt++ // 历史抽卡数+1
		// 抽卡前先自增 现实的抽卡次数是1...10+，而不是0...9+
		p.UpPoolInfo.FiveDrawTimes++
		p.UpPoolInfo.FourDrawTimes++
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
					p.UpPoolInfo.Statistics.WhichGold[p.UpPoolInfo.FiveDrawTimes]++
					p.UpPoolInfo.FiveDrawTimes = 0
					p.UpPoolInfo.Statistics.GoldEveryTen++
					p.UpPoolInfo.Statistics.GoldCnt++
					// up池处理大小保底问题
					card = p.handle5star(card)
				} else if role.Star == 4 { // 出货4星角色后重置次数
					p.UpPoolInfo.FourDrawTimes = 0
					p.UpPoolInfo.Statistics.PurpleCnt++
				}
			} else { // 出了武器[因为角色只有4、5星，role==nil则掉落一定为武器（3、4、5星）]
				weapon := csvs.GetConfigWeapon(card.Result)
				if weapon == nil {
					fmt.Println("抽卡结果掉落：武器配置信息错误")
					return
				} else if weapon.Star == 5 { // 五星武器情况 限定池是不会走到这里来的，保底五星是角色（应该）
					p.UpPoolInfo.Statistics.WhichGold[p.UpPoolInfo.FiveDrawTimes]++
					p.UpPoolInfo.FiveDrawTimes = 0
					p.UpPoolInfo.Statistics.GoldEveryTen++
					p.UpPoolInfo.Statistics.GoldCnt++
				} else if weapon.Star == 4 { // 四星武器
					p.UpPoolInfo.FourDrawTimes = 0
					p.UpPoolInfo.Statistics.PurpleCnt++
				} // else { // 三星武器 }
			}
			// 每累计十抽过后记录该十抽内的金色总量
			if p.UpPoolInfo.Statistics.HistoryCnt%10 == 0 {
				p.UpPoolInfo.Statistics.NGold[p.UpPoolInfo.Statistics.GoldEveryTen]++
				p.UpPoolInfo.Statistics.GoldEveryTen = 0 // 抽中N金
			}
			cards[card.Result]++
			// 抽到的东西放入背包
			player.ModBag.AddItem(player, card.Result, 1)
		}
	}

	// 根据统计开关打印信息
	p.PrintStatistic(switches)

	for k, v := range cards {
		fmt.Printf("%s\t次数:%d\n", csvs.GetItemName(k), v)
	}
	fmt.Printf("累计5星:%d，累计四星:%d\n", p.UpPoolInfo.Statistics.GoldCnt, p.UpPoolInfo.Statistics.PurpleCnt)
}

// todo 对应用户抽限定池1/2，这里map也要对应改变
func (p *ModPool) handle5star(card *csvs.ConfigDrop) *csvs.ConfigDrop {
	// 如果上次5星不是限定，则本次5星结果强制替换为限定池中的角色
	if p.UpPoolInfo.IsUpRole == csvs.LOGIC_TRUE {
		dropGroup := csvs.ConfigDropMap[100012]
		if dropGroup == nil {
			fmt.Println("大保底抽卡配置数据异常")
			return card
		}
		card = csvs.GetRandDropRecursion(dropGroup)
		if card == nil {
			fmt.Println("大保底抽卡结果数据异常")
			return card
		}
	}
	if card.DropId == 100012 { // 结果是5星限定角色，重置大保底
		p.UpPoolInfo.IsUpRole = csvs.LOGIC_FALSE
	} else { // 歪了...
		p.UpPoolInfo.IsUpRole = csvs.LOGIC_TRUE
	}
	return card
}
