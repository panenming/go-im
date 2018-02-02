package action

// 透传类消息 用于多端同步
const (
	READ = iota // 消息已读
	BACK        // 消息撤回
	UPTO        // 消息读数
	KICK        // 踢出群组
)

type ReadAction struct {
	MsgId string `json:"msgId"` // 消息id
}

type BackAction struct {
	MsgId string `json:"msgId"` // 消息id
}

type UptoAction struct {
	MsgId string `json:"msgId"` // 更新到的消息读数
}

type KickAction struct {
	Msg string `json:"msg"` // 被踢出群的提醒消息文本
}
