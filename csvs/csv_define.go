package csvs

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
)
