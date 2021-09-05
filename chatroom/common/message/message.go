package message

const (
	LoginMesType            = "loginMes"       //登录信息
	LoginResMesType         = "LoginResMes"    //登录返回的验证信息
	RegisterMesType         = "RegisterMes"    // 注册信息
	RegisterResMesType      = "RegisterResMes" // 注册返回的验证信息
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes" //群发消息
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"date"` //消息内容
}

//定义两个消息，需要再增加
//用户登录消息
type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"UserName"`
}

//返回登录信息
type LoginResMes struct {
	Code    int    `json:"code"`  //500表示还没注册 ，200表示登录成功
	Error   string `json:"error"` //返回错误信息
	UserIds []int  //保存用户id

}

//注册信息
type RegisterMes struct {
	//...
	User User `json:"user"`
}

//和LoginResMes相同，可以共用。但是从扩展性的角度考虑，我们重写一个结构体
type RegisterResMes struct {
	Code  int    `json:"code"`  //400表示该用户已经占用 ，200表示注册成功
	Error string `json:"error"` //返回错误信息
}

//为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"UserId"`
	Status int `json:"status"`
}

//这里定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffLine
	UserBusyStatus
)

type SmsMes struct {
	Content string `json:"content"`
	User //匿名的类型
}
