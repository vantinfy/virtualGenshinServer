package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
)

type ModAvatar struct {
	AvatarInfo map[int]*AvatarInfo
}

type AvatarInfo struct {
	AvatarId   int
	AvatarName string
}

// HasAvatar 玩家是否获得这个角色/头像
func (p *ModAvatar) HasAvatar(avatar int) bool {
	_, ok := p.AvatarInfo[avatar]
	return ok
}

// AddItem 添加头像
func (p *ModAvatar) AddItem(itemId int) {
	if _, ok := p.AvatarInfo[itemId]; ok {
		fmt.Println("头像已经拥有", itemId)
		return
	}

	config := csvs.GetAvatarConfig(itemId)
	if config == nil {
		fmt.Println("头像不存在", itemId)
		return
	}

	p.AvatarInfo[itemId] = &AvatarInfo{
		AvatarId:   itemId,
		AvatarName: csvs.GetItemConfig(itemId).ItemName,
	}
	fmt.Println("获得头像", p.AvatarInfo[itemId].AvatarName)
}

// CheckGetAvatar 获得角色后连带获得头像时的检查
func (p *ModAvatar) CheckGetAvatar(roleId int) {
	configAvatar := csvs.GetAvatarConfigByRoleId(roleId)
	if configAvatar == nil {
		fmt.Println("头像不存在")
		return
	}
	// 增加头像
	p.AddItem(configAvatar.AvatarId)
}
