package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigRole struct {
	RoleId          int    `json:"RoleId"`          // 角色id
	ItemName        string `json:"ItemName"`        // 角色名
	Star            int    `json:"Star"`            // 星级
	Stuff           int    `json:"Stuff"`           // 角色对应的命座材料
	StuffNum        int64  `json:"StuffNum"`        // 持有的角色命座材料数量
	StuffItem       int    `json:"StuffItem"`       // 重复获得角色时的材料返还
	StuffItemNum    int64  `json:"StuffItemNum"`    // 材料返还数量
	MaxStuffItem    int    `json:"MaxStuffItem"`    // 超过最大获得次数后转换材料
	MaxStuffItemNum int64  `json:"MaxStuffItemNum"` // 超过最大获得次数后转换材料数量
}

var RoleMap map[int]*ConfigRole

func init() {
	RoleMap = make(map[int]*ConfigRole)
	err := utils.GetCsvUtilMgr().LoadCsv("Role", RoleMap)
	if err != nil {
		fmt.Println("加载角色配置失败")
	}
}

func GetRoleConfig(roleId int) *ConfigRole {
	return RoleMap[roleId]
}
