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

	// 根据表自动创建地图
	for _, configMap := range csvs.ConfigMapMap {
		_, ok := m.MapInfo[configMap.MapId]
		if !ok {
			m.MapInfo[configMap.MapId] = m.NewMapInfo(configMap.MapId)
		}
	}

	// 初始化地图数据
	for _, v := range csvs.ConfigMapMondstadt {
		_, ok := m.MapInfo[v.MapId]
		if !ok {
			continue
		}

		_, ok = m.MapInfo[v.MapId].EventInfo[v.EventId]
		if !ok {
			m.MapInfo[v.MapId].EventInfo[v.EventId] = new(Event)
			m.MapInfo[v.MapId].EventInfo[v.EventId].EventId = v.EventId
			m.MapInfo[v.MapId].EventInfo[v.EventId].State = csvs.LOGIC_FALSE
		}
	}
	//for _, m2 := range m.MapInfo {
	//	fmt.Println("map", csvs.GetMapName(m2.MapId))
	//	for _, i2 := range m.MapInfo[m2.MapId].EventInfo {
	//		fmt.Println("event", csvs.GetEventName(i2.EventId))
	//	}
	//}
}

func (m *ModMap) NewMapInfo(mapId int) *Map {
	return &Map{
		MapId:     mapId,
		EventInfo: map[int]*Event{},
	}
}
