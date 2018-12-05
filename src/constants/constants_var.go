package constants

//npc buffs
const (
	BUFF_STREET_VISIT_CAP = "10"

	BUFF_RESTAURANT_OFFER = "20"
	BUFF_ENTERTAIN_OFFER  = "21"

	BUFF_SPECIFIED_RESTAURANT_OFFER = "30"
	BUFF_SPECIFIED_ENTERTAIN_OFFER  = "31"

	BUFF_RESTAURANT_PROFIT = "40"
	BUFF_ENTERTAIN_PROFIT  = "41"
	BUFF_MY_SHOP_PROFIT    = "42"
	BUFF_ALL_PROFIT        = "43"

	BUFF_SPECIFIED_RESTAURANT_PROFIT = "50"
	BUFF_SPECIFIED_ENTERTAIN_PROFIT  = "51"

	BUFF_RESTAURANT_CLEAN_NEED = "60"
	BUFF_ENTERTAIN_CLEAN_NEED  = "61"
	BUFF_ALL_CLEAN_NEED        = "62"

	BUFF_SPECIFIED_RESTAURANT_CLEAN_NEED = "70"
	BUFF_SPECIFIED_ENTERTAIN_CLEAN_NEED  = "71"

	BUFF_CLEAN_OBJECT_101_POWER = "80"
	BUFF_CLEAN_OBJECT_201_POWER = "81"
	BUFF_ALL_CLEAN_POWER        = "82"

	BUFF_CLEAN_OBJECT_101_SALARY = "90"
	BUFF_CLEAN_OBJECT_201_SALARY = "91"
	BUFF_ALL_CLEAN_SALARY        = "92"

	BUFF_UNLOCK_BUILDING = "00"
)

var NPC_BUFF_MAP = map[string]bool{
	BUFF_STREET_VISIT_CAP:                true, //整条街客流容量 百分比buff
	BUFF_RESTAURANT_OFFER:                true, //餐饮类建筑接待能力 百分比buff
	BUFF_ENTERTAIN_OFFER:                 true, //娱乐类建筑接待能力 百分比buff
	BUFF_SPECIFIED_RESTAURANT_OFFER:      true, //指定餐饮类建筑接待能力 百分比buff
	BUFF_SPECIFIED_ENTERTAIN_OFFER:       true, //指定娱乐类建筑接待能力 百分比buff
	BUFF_RESTAURANT_PROFIT:               true, //餐饮类建筑盈利能力 百分比buff
	BUFF_ENTERTAIN_PROFIT:                true, //娱乐类建筑盈利能力 百分比buff
	BUFF_MY_SHOP_PROFIT:                  true, //我的店盈利能力 百分比buff
	BUFF_ALL_PROFIT:                      true, //所有店盈利能力 百分比buff
	BUFF_SPECIFIED_RESTAURANT_PROFIT:     true, //指定餐饮类建筑盈利能力 百分比buff
	BUFF_SPECIFIED_ENTERTAIN_PROFIT:      true, //指定娱乐类建筑盈利能力 百分比buff
	BUFF_RESTAURANT_CLEAN_NEED:           true, //餐饮类建筑卫生需求 百分比buff
	BUFF_ENTERTAIN_CLEAN_NEED:            true, //娱乐类建筑卫生需求 百分比buff
	BUFF_ALL_CLEAN_NEED:                  true, //整条街的卫生需求 百分比buff
	BUFF_SPECIFIED_RESTAURANT_CLEAN_NEED: true, //指定餐饮建筑卫生需求 百分比buff
	BUFF_SPECIFIED_ENTERTAIN_CLEAN_NEED:  true, //指定娱乐建筑卫生需求 百分比buff
	BUFF_CLEAN_OBJECT_101_POWER:          true, //改变清洁机器人的清洁能力 百分比buff
	BUFF_CLEAN_OBJECT_201_POWER:          true, //改变清洁飞碟的清洁能力 百分比buff
	BUFF_ALL_CLEAN_POWER:                 true, //改变整条街的清洁能力 百分比buff
	BUFF_CLEAN_OBJECT_101_SALARY:         true, //改变清洁机器人的工资 百分比buff
	BUFF_CLEAN_OBJECT_201_SALARY:         true, //改变清洁飞碟的工资 百分比buff
	BUFF_ALL_CLEAN_SALARY:                true, //改变整条街的清洁工资 百分比buff
	BUFF_UNLOCK_BUILDING:                 true, //解锁某个特定建筑
}
