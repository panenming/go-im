package message

const (
	KindConnected  = iota + 1
	KindUserJoined // 用户加入
	KindUserLeft   // 用户离开
	KindStroke     // 用户draw
	KindClear      // 清除板面
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type User struct {
	ID    string `json:"id"`
	Color string `json:"color"` // 画线的颜色
}

// 连接的实体
// 这样的连接实体会将每次将所有的连接的客户端信息都发送给前端
type Connected struct {
	Kind  int    `json:"kind"`
	Color string `json:"color"`
	Users []User `json:"users"`
}

func NewConnected(color string, users []User) *Connected {
	return &Connected{
		Kind:  KindConnected,
		Color: color,
		Users: users,
	}
}

// 用户加入消息实体
type UserJoined struct {
	Kind int  `json:"kind"`
	User User `json:"user"`
}

func NewUserJoined(userId string, color string) *UserJoined {
	return &UserJoined{
		Kind: KindUserJoined,
		User: User{ID: userId, Color: color},
	}
}

// 用户离开实体
type UserLeft struct {
	Kind   int    `json:"kind"`
	UserID string `json:"userId"`
}

func NewUserLeft(userId string) *UserLeft {
	return &UserLeft{
		Kind:   KindUserLeft,
		UserID: userId,
	}
}

// draw 实体
type Stroke struct {
	Kind   int     `json:"kind"`
	UserID string  `json:"userId"`
	Points []Point `json:"points"`
	Finish bool    `json:"finish"`
}

// 清除实体
type Clear struct {
	Kind   int    `json:"kind"`
	UserID string `json:"userId"`
}
