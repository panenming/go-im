package message

// 消息type
const (
	SYSTEM = iota // 系统消息
	TXT           // 文本消息
	VOICE         // 语音消息
	PIC           // 图片消息
	VIDEO         // 短视频
	FILE          // 文件消息
	TXTURL        // 文本链接
	PICURL        // 图文链接
	ATTXT         // @消息
)

// 系统消息和TXT消息 Message Data中没有值

// 语音消息
type VOICEDATA struct {
	Id string `json:"id"` // 语音文件id
	T  int    `json:"t"`  // 语音时长 单位 秒s
}

// 图片消息
type PICDATA struct {
	Id string `json:"id"` // 图片文件id
	W  int    `json:"w"`  // 图片宽 单位pix
	H  int    `json:"h"`  // 图片高 单位pix
}

// 短视频消息
type VIDEODATA struct {
	Id string `json:"id"` // 短视频的id
	W  int    `json:"w"`  // 视频宽 单位pix
	H  int    `json:"h"`  // 视频高 单位pix
	T  int    `json:"t"`  // 视频长度 单位秒s
}

// 文件消息
type FILEDATA struct {
	S int32  `json:"s"` // 文件大小 单位byte
	X string `json:"x"` // 文件类型 待定
	N string `json:"n"` // 文件原名
}

// 文本链接
type TXTURLDATA struct {
	Url string `json:"url"` // 文件跳转链接
}

// 图文链接
type PICURLDATA struct {
	Id      string `json:"id"`      // 图片文件id
	Summary string `json:"summary"` // 摘要
	Url     string `json:"url"`     // 跳转链接
}

// @消息
type ATTXTDATA struct {
	UIds []string `json:"ids"` // @的人员id
}

// 其他类型消息待定
