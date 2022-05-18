package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigDrop struct {
	DropId int `json:"DropId"`
	Weight int `json:"Weight"`
	Result int `json:"Result"`
	IsEnd  int `json:"IsEnd"`
}

var ConfigDropSlice []*ConfigDrop

func init() {
	ConfigDropSlice = make([]*ConfigDrop, 0)
	err := utils.GetCsvUtilMgr().LoadCsv("Drop", &ConfigDropSlice)
	if err != nil {
		fmt.Println("加载掉落表配置失败")
		return
	}
}
