package module

import (
	"VirtualGenshinServer/csvs"
	"VirtualGenshinServer/manage"
	"fmt"
	"time"
)

type ShowRole struct {
	RoleId    int // 展示角色id
	RoleLevel int // 展示角色的等级
}

type ModPlayer struct {
	// todo 使用pb定义
	Uid            int64       // uid
	Name           string      // 名字
	Sign           string      // 签名
	Avatar         int         // 头像
	Card           int         // 名片
	PlayerLevel    int         // 等级
	PlayerExp      int         // 经验
	WorldLevel     int         // 世界等级
	WorldLevelNow  int         // 当前真实世界等级
	WorldLevelCool int64       // 世界等级主动降低后冷却时间
	Birthday       int         // 生日
	ShowTeam       []*ShowRole // 展示阵容(头像+等级)
	HideShowTeam   int         // 是否隐藏阵容角色详情
	ShowCard       []int       // 展示名片
	// others
	Prohibit int // 封禁状态
	IsGM     int // GM内部号
}

// SetAvatar 设置头像
func (p *ModPlayer) SetAvatar(player *Player, avatar int) {
	if player.ModAvatar.HasAvatar(avatar) {
		p.Avatar = avatar
		fmt.Println("设置头像成功", avatar)
	} else {
		fmt.Println("尚未拥有该头像", avatar)
	}
}

// GetAvatar 当前使用的头像
func (p *ModPlayer) GetAvatar() {
	avatarConfig := csvs.GetItemConfig(p.Avatar)
	if avatarConfig == nil {
		fmt.Println("头像", "未设置")
	} else {
		fmt.Println("头像", avatarConfig.ItemName)
	}
}

// SetCard 设置名片
func (p *ModPlayer) SetCard(player *Player, card int) {
	if player.ModCard.HasCard(card) {
		p.Card = card
		fmt.Println("设置名片成功", card)
	} else {
		fmt.Println("尚未拥有该名片", card)
	}
}

// GetCard 当前使用的名片
func (p *ModPlayer) GetCard() {
	cardConfig := csvs.GetItemConfig(p.Card)
	if cardConfig == nil {
		fmt.Println("名片", "未设置")
	} else {
		fmt.Println("名片", cardConfig.ItemName)
	}
}

// SetName 修改名字
func (p *ModPlayer) SetName(player *Player, name string) {
	banWords := manage.GetBanWords()
	if banWords.IsBanWord(name) {
		return
	}
	p.Name = name
	fmt.Println("设置新名字成功", p.Name)
}

// SetSign 修改签名
func (p *ModPlayer) SetSign(player *Player, sign string) {
	banWords := manage.GetBanWords()
	if banWords.IsBanWord(sign) {
		return
	}
	p.Sign = sign
	fmt.Println("设置新签名成功", p.Sign)
}

// AddExp 获得经验
func (p *ModPlayer) AddExp(player *Player, exp int) {
	p.PlayerExp += exp

	for {
		levelConfig := csvs.GetNowLevelConfig(p.PlayerLevel)
		if levelConfig == nil {
			fmt.Println("获取等级配置信息错误")
			break
		}
		if levelConfig.PlayerExp == 0 {
			// 满级
			break
		}
		// 突破任务未完成
		if levelConfig.ChapterId > 0 && !player.ModUniqueTask.IsTaskFinish(levelConfig.ChapterId) {
			break
		}
		// 升级
		if p.PlayerExp >= levelConfig.PlayerExp {
			p.PlayerLevel++
			//fmt.Println("升级了", p.PlayerLevel)
			p.PlayerExp -= levelConfig.PlayerExp
		} else {
			break
		}
	}
	fmt.Println("当前等级", p.PlayerLevel, "当前经验", p.PlayerExp)
}

// ReduceWorldLevel 降低世界等级
func (p *ModPlayer) ReduceWorldLevel() {
	if p.WorldLevel < csvs.REDUCE_WORLD_LEVEL_START {
		fmt.Println("降低世界等级失败", p.WorldLevel, "最低要求等级", csvs.REDUCE_WORLD_LEVEL_START)
		return
	}

	if p.WorldLevel-p.WorldLevelNow >= csvs.REDUCE_WORLD_LEVEL_MAX {
		fmt.Println("降低世界等级失败", p.WorldLevel)
		return
	}

	if time.Now().Unix() < p.WorldLevelCool {
		fmt.Println("降低世界等级失败，冷却中")
		return
	}

	// todo 完成突破任务时worldLevelNow需要跟worldLevel对上
	p.WorldLevelNow = p.WorldLevel - 1
	p.WorldLevelCool = time.Now().Unix() + csvs.REDUCE_WORLD_LEVEL_COOL_TIME
	fmt.Println("降低世界等级成功", p.WorldLevel, p.WorldLevelNow)
}

// ReturnWorldLevel 恢复世界等级
func (p *ModPlayer) ReturnWorldLevel() {
	if p.WorldLevel-p.WorldLevelNow < csvs.REDUCE_WORLD_LEVEL_MAX {
		fmt.Println("恢复世界等级失败", p.WorldLevel)
		return
	}

	if time.Now().Unix() < p.WorldLevelCool {
		fmt.Println("恢复世界等级失败，冷却中")
		return
	}

	p.WorldLevelNow++
	p.WorldLevelCool = time.Now().Unix() + csvs.REDUCE_WORLD_LEVEL_COOL_TIME
	fmt.Println("恢复世界等级成功", p.WorldLevel, p.WorldLevelNow)
}

// SetBirthDay 设置生日
func (p *ModPlayer) SetBirthDay(birth int) {
	// 设置生日方法只能调用一次
	if p.Birthday > 0 {
		fmt.Println("生日已经设置过啦，不能修改")
		return
	}

	month := birth / 100
	day := birth % 100
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		if day < 0 || day > 31 {
			fmt.Printf("生日设置错误,%v月没有%v日\n", month, day)
			return
		}
	case 4, 6, 9, 11:
		if day < 0 || day > 30 {
			fmt.Printf("生日设置错误,%v月没有%v日\n", month, day)
			return
		}
	case 2:
		if day < 0 || day > 29 {
			fmt.Printf("生日设置错误,%v月没有%v日\n", month, day)
			return
		}
	default:
		fmt.Println("不存在的月份", month)
		return
	}

	p.Birthday = birth
	fmt.Printf("生日设置成功: %v月%v日\n", month, day)

	if p.IsBirthDay() {
		fmt.Println("生日快乐^o^")
	} else {
		fmt.Println("期待你生日的到来")
	}
}

// IsBirthDay 是否玩家生日
func (p *ModPlayer) IsBirthDay() bool {
	month := time.Now().Month()
	day := time.Now().Day()
	if int(month) == p.Birthday/100 && day == p.Birthday%100 {
		return true
	}
	return false
}

// GetBirthDay 查询生日
func (p *ModPlayer) GetBirthDay() {
	if p.Birthday == 0 {
		fmt.Println("生日暂未设置")
	} else {
		fmt.Printf("生日是: %v月%v日\n", p.Birthday/100, p.Birthday%100)
	}
}

// SetShowCard 设置展示名片
func (p *ModPlayer) SetShowCard(player *Player, cards []int) {
	// 最多展示9个名片
	if len(cards) > csvs.SHOW_CARD_LIMIT {
		fmt.Println("名片展示数量异常")
		return
	}

	// 校验要展示的名片中是否有重复的
	cardExist := make(map[int]int)
	// 新的展示名片列表
	latestShowCard := make([]int, 0)

	for _, card := range cards {
		// 名片重复展示
		if _, ok := cardExist[card]; ok {
			fmt.Println("名片展示重复")
			return
		}
		if !player.ModCard.HasCard(card) {
			fmt.Println("设置了玩家尚未拥有的名片", card)
			return
		}
		latestShowCard = append(latestShowCard, card)
		cardExist[card] = 1
	}

	p.ShowCard = latestShowCard
	fmt.Println("设置新展示名片成功", p.ShowCard)
}

// SetShowTeam 设置展示阵容
func (p *ModPlayer) SetShowTeam(player *Player, roles []int) {
	// 最多展示8个角色
	if len(roles) > csvs.SHOW_TEAM_LIMIT {
		fmt.Println("阵容展示角色数量异常", len(roles))
		return
	}

	// 校验展示角色是否重复
	roleExist := make(map[int]int)
	// 最新的展示角色列表
	latestShowTeam := make([]*ShowRole, 0)

	for _, role := range roles {
		if _, ok := roleExist[role]; ok {
			fmt.Println("设置展示阵容角色重复")
			return
		}
		// 校验是否拥有这个角色
		if !player.ModRole.HasRole(role) {
			fmt.Println("设置了尚未拥有的角色", role)
			return
		}
		latestShowTeam = append(latestShowTeam, &ShowRole{
			RoleId:    role,
			RoleLevel: player.ModRole.GetRoleLevel(role),
		})
		roleExist[role] = 1
	}
	p.ShowTeam = latestShowTeam
	fmt.Println("设置展示阵容成功", p.ShowTeam)
}

// SetHideShowTeam 设置隐藏展示阵容的角色详情
func (p *ModPlayer) SetHideShowTeam(hide int) {
	if hide != csvs.LOGIC_TRUE && hide != csvs.LOGIC_FALSE {
		fmt.Println("设置隐藏展示阵容参数非法", hide)
		return
	}
	p.HideShowTeam = hide
}

// SetProhibit 账号封禁
func (p *ModPlayer) SetProhibit(prohibit int) {
	p.Prohibit = prohibit
}

// IsProhibit 账号是否处于封禁状态
func (p *ModPlayer) IsProhibit() bool {
	return int64(p.Prohibit) < time.Now().Unix()
}

// SetGm 设置gm号
func (p *ModPlayer) SetGm(gm int) {
	p.IsGM = gm
}
