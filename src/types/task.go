package types

type RawTask struct {
	TaskID int    `json:"taskID"`
	Desc   string `json:"desc"`
	Type   int    `json:"type"`
}

type TaskExtra struct {
	TaskID   int    `bson:"taskID" json:"taskID"`     // 任务id
	Progress int    `bson:"progress" json:"progress"` // 任务进度
	Desc     string `bson:"desc" json:"-"`            // 任务描述
	Target   int    `bson:"target" json:"-"`          // 目标
	Type     int    `bson:"type" json:"type"`         // 任务类型
	// RewardType int    `bson:"rewardType" json:"rewardType"` // 奖励类型
	// Reward     int    `bson:"reward" json:"reward"`         // 奖励的具体信息，根据奖励类型的不同而不同
	RewardGot bool   `bson:"rewardGot" json:"rewardGot"` // 奖励是否已经被获得
	Info      int    `bson:"info" json:"info"`           // 任务的一些额外信息
	DeviceId  string `json:"deviceId, omitempty"`
}
