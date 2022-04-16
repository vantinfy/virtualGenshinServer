package csvs

import (
	"VirtualGenshinServer/utils"
	"fmt"
)

type ConfigUniqueTask struct {
	TaskId    int `json:"TaskId"`    // 任务id
	SortType  int `json:"SortType"`  // 是否是突破任务
	OpenLevel int `json:"OpenLevel"` // 任务开放等级
	TaskType  int `json:"TaskType"`  // 任务类型
	Condition int `json:"Condition"` // 拓展
}

var ConfigUniqueTaskMap map[int]*ConfigUniqueTask

func init() {
	ConfigUniqueTaskMap = make(map[int]*ConfigUniqueTask)
	err := utils.GetCsvUtilMgr().LoadCsv("UniqueTask", ConfigUniqueTaskMap)
	if err != nil {
		fmt.Println("加载突破任务表失败")
	}
}
