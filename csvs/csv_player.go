package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigPlayerLevel struct {
	PlayerLevel int `json:"PlayerLevel"`
	PlayerExp   int `json:"PlayerExp"`
	WorldLevel  int `json:"WorldLevel"`
	ChapterId   int `json:"ChapterId"`
}

var ConfingPlayerLevelSlice []*ConfigPlayerLevel

func init() {
	err := utils.GetCsvUtilMgr().LoadCsv("PlayerLevel", &ConfingPlayerLevelSlice)
	if err != nil {
		fmt.Println("加载角色等级经验配置文件失败", err)
	}
}

// GetNowLevelConfig 获取等级配置信息
func GetNowLevelConfig(level int) *ConfigPlayerLevel {
	if level < 0 || level >= len(ConfingPlayerLevelSlice) {
		return nil
	}
	return ConfingPlayerLevelSlice[level-1]
}
