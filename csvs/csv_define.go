package csvs

const (
	EVENT_START  = 0
	EVENT_FINISH = 9
	EVENT_END    = 10
)

const (
	MAP_REFRESH_DAY  = 1
	MAP_REFRESH_WEEK = 2
	MAP_REFRESH_SELF = 3

	MAP_REFRESH_DAY_TIME  = 43200 // 刷新间隔 单位秒
	MAP_REFRESH_WEEK_TIME = 604800
	MAP_REFRESH_SELF_TIME = 180
)

const (
	LOGIC_FALSE = 0 // 逻辑0，目前用于是否隐藏展示阵容详情
	LOGIC_TRUE  = 1 // 逻辑1

	SHOW_CARD_LIMIT = 9 // 展示名片数量上限
	SHOW_TEAM_LIMIT = 8 // 展示阵容角色数量上限

	REDUCE_WORLD_LEVEL_START     = 5         // 主动降低世界等级的最低要求
	REDUCE_WORLD_LEVEL_MAX       = 1         // 最多能降低多少级
	REDUCE_WORLD_LEVEL_COOL_TIME = 24 * 3600 // 等级变动后的冷却时间 单位秒

	ADD_ROLE_TIMES_NORMAL_MIN = 2 // 重复获得角色材料转换
	ADD_ROLE_TIMES_NORMAL_MAX = 7 // 重复获得角色最大值后的材料转换

	WEAPON_BAG_MAX_ACCOUNT = 2000 // 武器背包最大容量
	RELICS_BAG_MAX_ACCOUNT = 1500 // 圣遗物背包最大容量

	FIVE_STAR_TIMES_LIMIT                  = 73   // 抽卡次数修正
	FIVE_STAR_TIMES_LIMIT_TRIGGER_INCREASE = 600  // 超过修正临界值后如果还没有抽到5星，则每次增加概率600/10000
	FOUR_STAR_TIMES_LIMIT                  = 8    // 抽卡次数修正
	FOUR_STAR_TIMES_LIMIT_TRIGGER_INCREASE = 5100 // 超过修正临界值后如果还没有抽到4星，则每次增加概率5100/10000
)
