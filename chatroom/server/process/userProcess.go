package process2

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//表明是哪个用户的字段
	UserId int
}

//通知所有用户在线的方法
//userId int要通知其他的人，我上线了
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历onlineUsers,然后一个一个的发送NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		//过滤掉自己
		if id == userId {
			continue
		}
		//开始通知，单独写一个方法
		up.NotifyMeOnline(userId)
	}
}

//服务器向客户端推送我的上线信息
func (this *UserProcess) NotifyMeOnline(userId int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("序列化出错")
		return
	}
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("序列化出错")
		return
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("发送消息出错")
		return
	}
}

//处理注册信息
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//核心代码
	//1.先从message中取出data，并反序列化成LoginMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("1反序列化失败")
		return
	}

	//先声明一个resMes，用来存储返回给客户端的登录是否成功的信息
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//到数据库中完成注册
	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 400
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 404
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
		//这里用户登录成功，把这个用户放入到UserMgr中
		fmt.Println("注册成功")
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("序列化失败")
		return
	}

	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化失败")
		return
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	return

}

//登录函数
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//核心代码
	//1.先从message中取出data，并反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("反序列化失败")
		return
	}

	//先声明一个resMes，用来存储返回给客户端的登录是否成功的信息
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//再声明一个LoginResMes，用来判定登录是否成功
	var loginResMes message.LoginResMes

	//到数据库中完成登录验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 300
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 400
			loginResMes.Error = "未知错误"
		}

		//先测试成功，然后根据返回的err显示错误信息
	} else {
		loginResMes.Code = 200
		this.UserId = loginMes.UserId
		//添加到UserMgr结构体
		userMgr.AddOnlineUser(this)

		//通知其他用户我上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线用户的id放到loginResMes.UserIds中
		//遍历UserMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UserIds = append(loginResMes.UserIds, id)
		}

		fmt.Println(user, "登录成功")

	}

	//3.将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("序列化失败")
		return err
	}
	fmt.Println("序列化成功1")
	resMes.Data = string(data) //将序列化好的data存储到resMes.Data中
	//将结构体resMes序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化失败")
		return err
	}
	fmt.Println("序列化成功2")
	//发送给客户端
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return

}
