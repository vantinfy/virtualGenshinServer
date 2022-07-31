package module

import (
	"fmt"
	"sync"
)

const (
	TASK_STATE_INIT   = 0 // 任务状态-初始
	TASK_STATE_DOING  = 1 // 任务状态-进行中
	TASK_STATE_FINISH = 2 // 任务状态-完成
)

// Player 玩家
type Player struct {
	ModPlayer     *ModPlayer     // 基础信息
	ModAvatar     *ModAvatar     // 头像
	ModCard       *ModCard       // 名片
	ModUniqueTask *ModUniqueTask // 任务
	ModRole       *ModRole       // 角色
	ModBag        *ModBag        // 背包
	ModWeapon     *ModWeapon     // 武器
	ModMap        *ModMap        // 地图
	ModRelics     *ModRelics     // 圣遗物
	ModCook       *ModCook       // 烹饪技能
	ModHome       *ModHome       // 家园物品
	ModPool       *ModPool       // 卡池
}

// --- 初始化相关 ---

func NewPlayer() *Player {
	player := new(Player)

	// todo 各个模块的new放到各自的文件中
	player.ModPlayer = new(ModPlayer)

	player.ModAvatar = new(ModAvatar)
	player.ModAvatar.AvatarInfo = make(map[int]*AvatarInfo)

	player.ModCard = new(ModCard)
	player.ModCard.CardInfo = make(map[int]*CardInfo)

	player.ModUniqueTask = new(ModUniqueTask)

	player.ModRole = new(ModRole)
	player.ModRole.RoleInfo = make(map[int]*RoleInfo)

	player.ModBag = new(ModBag)
	player.ModBag.BagInfo = make(map[int]*ItemInfo)

	player.ModWeapon = new(ModWeapon)
	player.ModWeapon.WeaponInfo = make(map[int]*Weapon)

	player.ModMap = new(ModMap)
	player.ModMap.InitData()

	player.ModRelics = new(ModRelics)
	player.ModRelics.RelicsInfo = make(map[int64]*Relics) // 获得圣遗物的次数要比武器多得多，因此键用int64

	player.ModCook = new(ModCook)
	player.ModCook.CookInfo = make(map[int]*Cook)

	player.ModHome = new(ModHome)
	player.ModHome.HomeItemInfo = make(map[int]*HomeItemInfo)

	player.ModPool = new(ModPool)
	player.ModPool.UpPoolInfo = &PoolInfo{}
	player.ModPool.UpPoolInfo.Statistics = &Statistic{}

	// --- 数据初始化 ---
	player.ModPlayer.Name = "旅行者"
	player.ModPlayer.PlayerLevel = 1
	player.ModPlayer.WorldLevel = 1
	player.ModPlayer.WorldLevelNow = 1

	return player
}

// Wait 多个玩家协程的时候syncWait
var Wait = sync.WaitGroup{}

// Run 处理客户端请求
func (p *Player) Run() {
	//pwd:/main/bag/add tree:/
	fmt.Println("====== 输入命令 ======")
	cmd := ""
	// 初次登录，选择角色
	for {
		fmt.Println("选择初始角色 0:旅行者妹妹(荧) 1:旅行者哥哥(空) -1:退出")
		_, _ = fmt.Scan(&cmd)
		if cmd == "0" {
			p.ModBag.AddItem(p, 2000001, 1)
			break
		}
		if cmd == "1" {
			p.ModBag.AddItem(p, 2000002, 1)
			break
		}
		if cmd == "-1" {
			Wait.Done()
			return
		}
	}
	for {
		fmt.Println("///main/// -1:展示路径 0:退出 1:查看玩家信息 2:进入背包 3:抽up池 4:查看角色 5:仓检抽卡(天选幸运儿) 6:平衡仓检抽卡")
		_, _ = fmt.Scan(&cmd)
		if cmd == "-1" {
			fmt.Println(`
主界面(当前位置) ─┬─ 个人信息 #1
				├─ 背包 #2
				├─ 抽卡(up池子) #3
				├─ 查看角色 #4
				├─ 仓检抽卡 #5
				├─ 平衡仓检抽卡 #6
				└─ 地图(开发中)
			`)
		} else if cmd == "0" {
			break
		} else if cmd == "1" {
			showInfos(p)
			fmt.Println("****** 分隔线 ******")
			personalInfo(p)
		} else if cmd == "2" {
			bag(p, cmd)
		} else if cmd == "3" {
			times := 0
			fmt.Println("输入要抽卡的次数")
			_, _ = fmt.Scan(&times)
			p.HandleDraw(p, 0, times)
		} else if cmd == "4" {
			p.HandleShowRoles()
		} else if cmd == "5" {
			times := 0
			fmt.Println("输入(仓检)抽卡次数")
			_, _ = fmt.Scan(&times)
			p.HandleDrawWithCheck(p, 0, times)
		} else if cmd == "6" {
			times := 0
			fmt.Println("输入(平衡仓检)抽卡次数")
			_, _ = fmt.Scan(&times)
			p.HandleDrawWithCheckAvg(p, 0, times)
		}
	}
	Wait.Done()
}

func showInfos(p *Player) {
	fmt.Println("*** 旅行者信息如下 ***")
	fmt.Println("昵称", p.ModPlayer.Name)
	fmt.Println("签名", p.ModPlayer.Sign)
	fmt.Println("等级", p.ModPlayer.PlayerLevel)
	fmt.Println("经验", p.ModPlayer.PlayerExp)
	p.ModPlayer.GetBirthDay()
	fmt.Println("世界等级", p.ModPlayer.WorldLevelNow)
	p.ModPlayer.GetAvatar()
	p.ModPlayer.GetCard()
}

func personalInfo(p *Player) {
	outloop := true
	cmd := ""
	for outloop {
		fmt.Println("//show// -1:展示路径 0:返回上一级 1:设置昵称 2:设置签名 3:头像 4:名片 5:生日")
		_, _ = fmt.Scan(&cmd)
		switch cmd {
		case "-1":
			fmt.Println(`
主界面#-1─┬─ 个人信息(当前位置)
		 │    ├─设置昵称 #1
		 │	  ├─设置签名 #2
		 │	  ├─头像 #3
		 │	  ├─名片 #4
		 │	  └─生日 #5
		 ├─ 背包
		 └─ 地图(开发中)
					`)
		case "0":
			outloop = false
		case "1":
			fmt.Println("输入新名字")
			_, _ = fmt.Scan(&cmd)
			p.ReceiveSetName(cmd)
		case "2":
			fmt.Println("输入新签名")
			_, _ = fmt.Scan(&cmd)
			p.ReceiveSetSign(cmd)
		case "3":
			avatarInfo(p)
		case "4":
			cardInfo(p)
		case "5":
			birthDay(p)
		}
	}
}

func avatarInfo(p *Player) {
	loop := true
	cmd := ""
	itemId := 0
	for loop {
		fmt.Println("//头像// -1:展示路径 0:返回上一级 1:查询拥有的头像 2:设置头像")
		_, _ = fmt.Scan(&cmd)
		switch cmd {
		case "-1":
			fmt.Println(`
主界面#-1─┬─ 个人信息
		 │    ├─设置昵称
		 │	  ├─设置签名
		 │	  ├─头像(当前位置)
	 	 │	  │	 ├─查询拥有的头像 #1
		 │	  │	 └─设置头像 #2
		 │	  ├─名片 #4
		 │	  └─生日 #5
		 ├─ 背包
		 └─ 地图(开发中)
					`)
		case "0":
			loop = false
		case "1":
			for i, avatar := range p.ModAvatar.AvatarInfo {
				fmt.Println(i, avatar.AvatarName)
			}
		case "2":
			fmt.Println("输入新头像id")
			_, _ = fmt.Scan(&itemId)
			p.ReceiveSetAvatar(itemId)
		}
	}
}

func cardInfo(p *Player) {
	loop := true
	cmd := ""
	itemId := 0
	for loop {
		fmt.Println("//名片// -1:展示路径 0:返回上一级 1:查询拥有的名片 2:设置名片")
		_, _ = fmt.Scan(&cmd)
		switch cmd {
		case "-1":
			fmt.Println(`
主界面#-1─┬─ 个人信息
		 │    ├─设置昵称
		 │	  ├─设置签名
		 │	  ├─头像
		 │	  ├─名片(当前位置)
	 	 │	  │	 ├─查询拥有的名片 #1
		 │	  │	 └─设置名片 #2
		 │	  └─生日 #5
		 ├─ 背包
		 └─ 地图(开发中)
					`)
		case "0":
			loop = false
		case "1":
			for i, card := range p.ModCard.CardInfo {
				fmt.Println(i, card.CardName)
			}
		case "2":
			fmt.Println("输入新名片id")
			_, _ = fmt.Scan(&itemId)
			p.ReceiveSetCard(itemId)
		}
	}
}

func birthDay(p *Player) {
	loop := true
	cmd := ""
	itemId := 0
	for loop {
		fmt.Println("//生日// -1:展示路径 0:返回上一级 1:查询生日 2:设置生日")
		_, _ = fmt.Scan(&cmd)
		switch cmd {
		case "-1":
			fmt.Println(`
主界面#-1─┬─ 个人信息
		 │    ├─设置昵称
		 │	  ├─设置签名
		 │	  ├─头像
		 │	  ├─名片
		 │	  └─生日(当前位置)
	 	 │	   	 ├─查看生日 #1
		 │	   	 └─设置生日 #2
		 ├─ 背包
		 └─ 地图(开发中)
					`)
		case "0":
			loop = false
		case "1":
			p.ModPlayer.GetBirthDay()
		case "2":
			fmt.Println("输入生日（月数*100+天数，示例：6月28对应628）")
			_, _ = fmt.Scan(&itemId)
			p.SetBirthDay(itemId)
		}
	}
}

func bag(p *Player, cmd string) {
	itemId := 0
	itemAccount := 0
	outloop := true
	for outloop {
		fmt.Println("//背包// -1:展示路径 0:返回上一级 1:添加物品 2:移除物品 3:使用物品")
		_, _ = fmt.Scan(&cmd)
		switch cmd {
		case "-1":
			fmt.Println(`
主界面#0─┬─ 个人信息
        ├─ 背包(当前位置)
 		│	  ├─添加物品 #1
 		│	  ├─移除物品 #2
		│	  └─使用物品 #3
		└─ 地图(开发中)
					`)
		case "0":
			outloop = false
		case "1":
			fmt.Println("输入添加的物品id")
			_, _ = fmt.Scan(&itemId)
			fmt.Println("输入添加的物品数量")
			_, _ = fmt.Scan(&itemAccount)
			p.ModBag.AddItem(p, itemId, int64(itemAccount))
		case "2":
			fmt.Println("输入移除的物品id")
			_, _ = fmt.Scan(&itemId)
			fmt.Println("输入移除的物品数量")
			_, _ = fmt.Scan(&itemAccount)
			p.ModBag.RemoveItem(p, itemId, int64(itemAccount))
		case "3":
			fmt.Println("输入使用的物品id")
			_, _ = fmt.Scan(&itemId)
			fmt.Println("输入使用的物品数量")
			_, _ = fmt.Scan(&itemAccount)
			p.HandleUseItem(itemId, int64(itemAccount))
		}
	}
}

// --- 由用户主动发起的方法 ---

// ReceiveSetAvatar 设置头像
func (p *Player) ReceiveSetAvatar(avatar int) {
	p.ModPlayer.SetAvatar(p, avatar)
}

// ReceiveSetCard 设置名片
func (p *Player) ReceiveSetCard(card int) {
	p.ModPlayer.SetCard(p, card)
}

//ReceiveSetName 修改玩家名字
func (p *Player) ReceiveSetName(name string) {
	p.ModPlayer.SetName(p, name)
}

// ReceiveSetSign 修改签名
func (p *Player) ReceiveSetSign(sign string) {
	p.ModPlayer.SetSign(p, sign)
}

// ReduceWorldLevel 降低世界等级
func (p *Player) ReduceWorldLevel() {
	p.ModPlayer.ReduceWorldLevel()
}

// ReturnWorldLevel 恢复世界等级
func (p *Player) ReturnWorldLevel() {
	p.ModPlayer.ReturnWorldLevel()
}

// SetBirthDay 设置生日
func (p *Player) SetBirthDay(birth int) {
	p.ModPlayer.SetBirthDay(birth)
}

// SetShowCard 设置展示名片
func (p *Player) SetShowCard(showCards []int) {
	p.ModPlayer.SetShowCard(p, showCards)
}

// SetShowTeam 设置展示阵容
func (p *Player) SetShowTeam(roles []int) {
	p.ModPlayer.SetShowTeam(p, roles)
}

// SetHideShowTeam 设置隐藏展示阵容的角色详情
func (p *Player) SetHideShowTeam(hide int) {
	p.ModPlayer.SetHideShowTeam(hide)
}

// HandleUseItem 使用物品
func (p *Player) HandleUseItem(itemId int, account int64) {
	p.ModBag.UseItem(p, itemId, account)
}

// HandleDraw 限定池抽卡
func (p *Player) HandleDraw(player *Player, pool int, times int) {
	p.ModPool.UpPoolDraw(player, times, WithNGold(), WithWhichGold(), WithDropCnt(), WithHistoryCnt())
}

func (p *Player) HandleShowRoles() {
	p.ModRole.ShowRoles()
}

// HandleDrawWithCheck 限定池抽卡(仓检版)
func (p *Player) HandleDrawWithCheck(player *Player, pool int, times int) {
	p.ModPool.UpPoolDrawWithCheck(player, times, WithNGold(), WithWhichGold(), WithDropCnt(), WithHistoryCnt())
}

func (p *Player) HandleDrawWithCheckAvg(player *Player, pool int, times int) {
	p.ModPool.UpPoolDrawWithCheckAvg(player, times, WithNGold(), WithWhichGold(), WithDropCnt(), WithHistoryCnt())
}
