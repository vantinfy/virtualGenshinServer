package module

import (
	"VirtualGenshinServer/csvs"
)

type Map struct {
	MapId     int `json:"MapId"`
	EventInfo map[int]*Event
}

type Event struct {
	EventId int
	State   int
}

type ModMap struct {
	MapInfo map[int]*Map
}

func (m *ModMap) InitData() {
	if m.MapInfo == nil {
		m.MapInfo = map[int]*Map{}
	}

	// MapInfo[1](蒙德地图) 如果不存在就new
	_, ok := m.MapInfo[1]
	if !ok {
		m.MapInfo[1] = new(Map)
		m.MapInfo[1].EventInfo = map[int]*Event{}
	}

	// 初始化地图数据
	for _, v := range csvs.ConfigMapMondstadt {
		_, ok := m.MapInfo[1].EventInfo[v.EventId]
		if !ok {
			m.MapInfo[1].EventInfo[v.EventId] = new(Event)
			m.MapInfo[1].EventInfo[v.EventId].EventId = v.EventId
			m.MapInfo[1].EventInfo[v.EventId].State = csvs.LOGIC_FALSE
		}
	}
}
