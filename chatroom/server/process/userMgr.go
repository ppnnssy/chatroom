package process2

import "fmt"

//因为在很多地方都会使用到，并且服务端只有一个UserMgr，所以定义成全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对userMgr的初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//完成对onlineUser的添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

//完成对onlineUser的删除
func (this *UserMgr) DeOnlineUser(up *UserProcess) {
	delete(this.onlineUsers, up.UserId)
}

//完成对onlineUser的查询和返回
func (this *UserMgr)GetOnlineUserById(userId int)(up *UserProcess ,err error) {
	//从map中取出一个值，带检测
	up,ok := this.onlineUsers[userId]
	if !ok{ //说明当前用户不在线
		err=fmt.Errorf("用户id不存在或不在线")
		return
	}
	return
}


//返回所有在线的用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}
