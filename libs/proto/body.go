package proto

import (
	"encoding/json"
)

const (
	HTTP_GET = iota
	HTTP_POST
	HTTP_PUT
	HTTP_DELETE
)

// 从客户端发起的握手实体
type HandShake struct {
	DId        string `json:"dId"`        // 客户端的唯一id
	Token      string `json:"token"`      // 登录后获取的token
	DviceToken string `json:"dviceToekn"` // 如果是ios，需要推送证书token
}

// 心跳请求实体 心跳实体为nil

// 发送文本消息的实体
type Message struct {
	Id     string      `json:"id"`     // 消息id
	Fr     string      `json:"fr"`     // 来自谁
	To     string      `json:"to"`     // 发送给哪个群
	Type   string      `json:"type"`   // 消息类型
	Txt    string      `json:"txt"`    // 消息提醒文本
	Ct     int32       `json:"ct"`     // 创建时间
	Extra  interface{} `json:"extra"`  // 和发送者相关的属性
	Data   interface{} `json:"data"`   // 消息其他字段
	Status int         `json:"status"` // 消息有效性 0 有效 1 无效
}

type Extra struct {
	PhotoId string `json:"photoId"` // 人员头像
	NName   string `json:"nName"`   // 人员在群里的昵称
}

// 调用网关实体
type Rpc struct {
	Path   string      `json:"path"`   // 请求地址
	Method int         `json:"method"` // 请求类型 GET POST PUT DELETE
	Param  interface{} `json:"param"`  // 请求的参数
}

// 透传类消息
type Action struct {
	Id   string      `json:"id"`   // action id
	Type int         `json:"type"` // action类型
	GId  string      `json:"gId"`  // 操作的群组id
	Ct   int32       `json:"ct"`   // action 生成时间
	Data interface{} `json:"data"` // 对象实体
}

// 所有的请求返回统一的消息体（心跳除外）
type Resp struct {
	Err   string      `json:"err"`   // 错误信息，如果请求正确此值为nil
	ErrNo int         `json:"errNo"` // 错误编号 0 表示正确，其他待定
	Data  interface{} `json:"data"`  // 请求正确的时候返回的数据，收集端需要转化为[]byte
}
