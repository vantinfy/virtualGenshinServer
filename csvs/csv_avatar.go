package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

// ConfigAvatar 注意csv文件首行字段名要与json tag一致
type ConfigAvatar struct {
	AvatarId int `json:"AvatarId"`
	Check    int `json:"Check"`
}

var (
	ConfigAvatarMap         map[int]*ConfigAvatar // 键是头像id
	ConfigAvatarMapByRoleId map[int]*ConfigAvatar // 键是角色id
)

func init() {
	ConfigAvatarMap = make(map[int]*ConfigAvatar)
	ConfigAvatarMapByRoleId = make(map[int]*ConfigAvatar)

	// load base csv
	err := utils.GetCsvUtilMgr().LoadCsv("Avatar", ConfigAvatarMap)
	if err != nil {
		fmt.Println("头像配置文件加载失败", err)
	}
	for _, v := range ConfigAvatarMap {
		ConfigAvatarMapByRoleId[v.Check] = v
	}
}

// GetAvatarConfig 通过id获取头像配置
func GetAvatarConfig(itemId int) *ConfigAvatar {
	return ConfigAvatarMap[itemId]
}

// GetAvatarConfigByRoleId 通过角色id获取对应头像配置
func GetAvatarConfigByRoleId(roleId int) *ConfigAvatar {
	return ConfigAvatarMapByRoleId[roleId]
}
