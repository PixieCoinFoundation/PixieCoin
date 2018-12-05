package tasks

import (
	"appcfg"
	// "errors"
	"constants"
	. "pixie_contract/api_specification"
	"sync"
)

import (
	jl "jsonloader"
	// . "logger"
	// . "player_mem"
	. "types"
)

var (
	AllTasks     []*Task
	AllTaskMap   map[int]*Task
	allTasksLock sync.RWMutex

	allRawTasks   []*RawTask
	allRawTaskMap map[int]*RawTask

	dayTaskMap       map[int]int
	allTaskTargetMap map[int]int
	allTaskNameMap   map[int]string
)

func init() {

	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}
	loadTask()
}

func loadTask() {
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}

	allRawTasks = make([]*RawTask, 0)
	allRawTaskMap = make(map[int]*RawTask)
	AllTaskMap = make(map[int]*Task)
	dayTaskMap = make(map[int]int)
	allTaskTargetMap = make(map[int]int)
	allTaskNameMap = make(map[int]string)

	allTasksLock.Lock()
	defer allTasksLock.Unlock()

	fn := "scripts/data/data_task.json"
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		fn = "scripts/data/data_task_kor.json"
	}

	if err := jl.LoadFile(fn, &AllTasks); err != nil {
		panic(err)
	}

	if err := jl.LoadFile(fn, &allRawTasks); err != nil {
		panic(err)
	}

	for _, at := range AllTasks {
		allTaskTargetMap[at.TaskID] = at.Target
		if at.Type == 1 {
			dayTaskMap[at.TaskID] = 1
		}
		AllTaskMap[at.TaskID] = at
	}

	for _, at := range allRawTasks {
		allTaskNameMap[at.TaskID] = at.Desc
		allRawTaskMap[at.TaskID] = at
	}
}

func GetTaskNameType(id int) (string, int) {
	t := allRawTaskMap[id]
	return t.Desc, t.Type
}

func GetTaskName(id int) string {
	return allTaskNameMap[id]
}

func IsDayTask(taskID int) bool {
	return dayTaskMap[taskID] > 0
}

func GetTaskTarget(taskID int) int {
	return allTaskTargetMap[taskID]
}

// func ReloadTasks() (succeed bool) {
// 	loadTask()
// 	return true
// }

func GetTask(taskID int) (template *Task) {
	allTasksLock.RLock()
	defer allTasksLock.RUnlock()
	// for _, v := range AllTasks {
	// 	if v.TaskID == taskID {
	// 		template = v
	// 	}
	// }
	return AllTaskMap[taskID]
}

// func IncreTask(player *GFPlayer, incre int, taskID int) (isOver bool, changed bool, currentProgress int, err error) {
// 	// 从tasks模板中获取task的奖励类型和奖励信息
// 	var template *Task
// 	allTasksLock.RLock()
// 	for _, v := range AllTasks {
// 		if v.TaskID == taskID {
// 			template = v
// 		}
// 	}
// 	allTasksLock.RUnlock()

// 	return player.IncreTask(taskID, incre, template)
// }

// func UpdateTask(player *GFPlayer, taskID int, progress int) (isOver bool, dayProgress int, err error) {
// 	// 从tasks模板中获取task的奖励类型和奖励信息
// 	var template *Task
// 	allTasksLock.RLock()
// 	for _, v := range AllTasks {
// 		if v.TaskID == taskID {
// 			template = v
// 		}
// 	}
// 	allTasksLock.RUnlock()

// 	return player.UpdateTask(taskID, progress, template)
// }

// func GetTaskReward(player *GFPlayer, taskID int, deviceId string) (rewardType int, rewardInfo int, err error) {
// 	// 从tasks模板中获取task的奖励类型和奖励信息
// 	var template *Task
// 	allTasksLock.RLock()
// 	for _, v := range AllTasks {
// 		if v.TaskID == taskID {
// 			template = v
// 		}
// 	}
// 	allTasksLock.RUnlock()
// 	if template != nil {
// 		// 存在此任务，则获得此任务的奖励
// 		return player.GetTaskReward(template, deviceId)
// 	} else {
// 		err = errors.New("no such task found")
// 		Err(err.Error())

// 		return
// 	}
// 	return
// }
