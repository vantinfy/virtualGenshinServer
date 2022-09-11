package module

import (
	"VirtualGenshinServer/csvs"
	"fmt"
	"time"
)

type Map struct {
	MapId     int `json:"MapId"`
	EventInfo map[int]*Event
}

type Event struct {
	EventId       int
	State         int
	NextResetTime int64 // 下次刷新时间
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
	for _, v := range csvs.ConfigMapEventMap {
		_, ok := m.MapInfo[v.MapId]
		if !ok {
			continue
		}

		_, ok = m.MapInfo[v.MapId].EventInfo[v.EventId]
		if !ok {
			m.MapInfo[v.MapId].EventInfo[v.EventId] = new(Event)
			m.MapInfo[v.MapId].EventInfo[v.EventId].EventId = v.EventId
			m.MapInfo[v.MapId].EventInfo[v.EventId].State = csvs.EVENT_START
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

func (m *ModMap) GetEventList(config *csvs.ConfigMap) {
	_, ok := m.MapInfo[config.MapId]
	if !ok {
		return
	}
	for _, event := range m.MapInfo[config.MapId].EventInfo {
		m.CheckRefresh(event)
		lastTime := event.NextResetTime - time.Now().Unix()
		noticeTime := ""
		if lastTime < 0 {
			noticeTime = "已刷新"
		} else {
			noticeTime = fmt.Sprintf("%d秒后刷新", lastTime)
		}
		fmt.Println(fmt.Sprintf("事件id: %d 事件名: %s 事件状态: %d, %s", event.EventId, csvs.GetEventName(event.EventId), event.State, noticeTime))
	}
}

func (m *ModMap) SetEventState(mapId, eventId, state int) {
	_, ok := m.MapInfo[mapId]
	if !ok {
		fmt.Println("地图不存在")
		return
	}
	_, ok = m.MapInfo[mapId].EventInfo[eventId]
	if !ok {
		fmt.Println("事件不存在")
		return
	}
	if m.MapInfo[mapId].EventInfo[eventId].State >= state {
		fmt.Println("设置状态异常")
		return
	}
	m.MapInfo[mapId].EventInfo[eventId].State = state

	if state > 0 {
		config := csvs.GetEventConfig(eventId)
		if config == nil {
			return
		}
		switch config.RefreshType {
		case csvs.MAP_REFRESH_SELF:
			m.MapInfo[mapId].EventInfo[eventId].NextResetTime = time.Now().Unix() + csvs.MAP_REFRESH_SELF_TIME
		}
	}
}

func (m *ModMap) RefreshDay() {
	for _, mapInfo := range m.MapInfo {
		for _, event := range m.MapInfo[mapInfo.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[event.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_DAY {
				continue
			}
			event.State = csvs.EVENT_START
		}
	}
}

func (m *ModMap) RefreshWeek() {
	for _, mapInfo := range m.MapInfo {
		for _, event := range m.MapInfo[mapInfo.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[event.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_WEEK {
				continue
			}
			event.State = csvs.EVENT_START
		}
	}
}

func (m *ModMap) RefreshSelf() {
	for _, mapInfo := range m.MapInfo {
		for _, event := range m.MapInfo[mapInfo.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[event.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_SELF {
				continue
			}
			// 这里要判定刷新时间
			if time.Now().Unix() >= event.NextResetTime {
				event.State = csvs.EVENT_START
			}
		}
	}
}

func (m *ModMap) CheckRefresh(event *Event) {
	// 还没到刷新时间的直接返回
	if event.NextResetTime > time.Now().Unix() {
		return
	}

	eventConfig := csvs.GetEventConfig(event.EventId)
	if eventConfig == nil {
		return
	}
	event.State = csvs.EVENT_START
	switch eventConfig.RefreshType {
	case csvs.MAP_REFRESH_DAY:
		count := time.Now().Unix() / csvs.MAP_REFRESH_DAY_TIME
		count++
		event.NextResetTime = count * csvs.MAP_REFRESH_DAY_TIME
	case csvs.MAP_REFRESH_WEEK:
		count := time.Now().Unix() / csvs.MAP_REFRESH_WEEK_TIME
		count++
		event.NextResetTime = count * csvs.MAP_REFRESH_WEEK_TIME
	case csvs.MAP_REFRESH_SELF:

	}
}
