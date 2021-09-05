package process

import (
	"chatroom/client/model"
	"chatroom/common/message"
	"fmt"
)

//客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
//客户端很多地方会使用到CurUser，做成全局变量
//用户登录成功后需要完成对CurUser的初始化
var CurUser model.CurUser


//在客户端显示当前在线的用户
func outPutOnlineUsers() {
	//遍历onlineUsers
	fmt.Println("当前在线用户列表")
	for id,_:=range onlineUsers{
		fmt.Println("用户id：\t",id)
	}
}

//编写一个方法，处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}

	user.UserStatus=notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId]=user

	outPutOnlineUsers()
}
