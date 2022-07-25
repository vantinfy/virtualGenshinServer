package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type ModRole struct {
	RoleInfo map[int]*RoleInfo
}

type RoleInfo struct {
	RoleId   int // 角色id
	GetTimes int // 该角色获得次数
	Star     int // 角色星级
	// extension
	RoleName  string // 角色名字
	RoleLevel int    // 角色等级
	RoleExp   int    // 角色经验
	WeaponId  int    // 角色武器
	Relics    []int  // 圣遗物
	Fate      int    // 命之座
	Talent    int    // 天赋
}

// HasRole 查询是否拥有角色
func (r *ModRole) HasRole(roleId int) bool {
	_, ok := r.RoleInfo[roleId]
	return ok
}

// GetRoleLevel 获得角色等级
func (r *ModRole) GetRoleLevel(roleId int) int {
	return 80
}

// GetRoleExp 获得角色经验
func (r *ModRole) GetRoleExp(roleId int) int {
	return 80
}

// AddItem 获得角色
func (r *ModRole) AddItem(player *Player, roleId int, account int64) {
	config := csvs.GetRoleConfig(roleId)
	if config == nil {
		fmt.Println("非法角色", roleId)
		return
	}
	// 可能需要加个校验，如果初始选择了妹妹则不能获得哥哥，反之同理（不过目前暂时不做
	for i := 0; i < int(account); i++ {
		if role, ok := r.RoleInfo[roleId]; ok {
			// 已经获得过
			r.RoleInfo[roleId].GetTimes++
			if role.GetTimes >= csvs.ADD_ROLE_TIMES_NORMAL_MIN && role.GetTimes <= csvs.ADD_ROLE_TIMES_NORMAL_MAX {
				// 命星与材料
				player.ModBag.AddItemToBag(config.Stuff, config.StuffNum)
				player.ModBag.AddItemToBag(config.StuffItem, config.StuffItemNum)
			} else {
				// 最大材料转换
				player.ModBag.AddItemToBag(config.MaxStuffItem, config.MaxStuffItemNum)
			}
			//fmt.Printf("命座转换: %#v\n", config)
		} else {
			// 第一次获得
			role := new(RoleInfo)
			role.RoleId = roleId
			role.RoleName = config.ItemName
			role.GetTimes = 1
			role.Star = config.Star
			r.RoleInfo[roleId] = role
			// 增加头像
			player.ModAvatar.CheckGetAvatar(roleId)
			player.ModCard.CheckGetCard(roleId, 10)
		}
	}
	fmt.Printf("第%v次，共获得了%v个角色%v\n", r.RoleInfo[roleId].GetTimes, account, r.RoleInfo[roleId].RoleName)
}

// ShowRoles 查看拥有哪些角色
func (r *ModRole) ShowRoles() {
	fmt.Println("角色如下")
	for _, info := range r.RoleInfo {
		info.ShowInfos()
	}
}

func (r *RoleInfo) ShowInfos() {
	fmt.Printf("%s累计获得次数: %d\n", csvs.GetItemName(r.RoleId), r.GetTimes)
}
