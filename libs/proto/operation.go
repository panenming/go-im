package proto

const (
	// 握手
	OP_HANDSHAKE       = int32(0)
	OP_HANDSHAKE_REPLY = int32(1)
	// 心跳
	OP_HEARTBEAT       = int32(2)
	OP_HEARTBEAT_REPLY = int32(3)
	// 发送文本消息
	OP_SEND_SMS       = int32(4)
	OP_SEND_SMS_REPLY = int32(5)
	// 长连接做网关调用下属微服务
	OP_GATEWAY       = int32(6)
	OP_GATEWAY_REPLY = int32(7)
	// proto 相关
	OP_PROTO_READY  = int32(100)
	OP_PROTO_FINISH = int32(101)
	// 其他待定义
)
