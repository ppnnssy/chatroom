package main

import (
	"chatroom/common/message"
	"chatroom/server/process"
	"chatroom/server/utils"
	"fmt"
	"io"
	"net"
)

//先创建一个processor的结构体
type Processor struct {
	Conn net.Conn
}

//总控函数
func (this *Processor) process2 ()(err error){
	//循环接收客户端发送的信息
	for {
		//创建一个Transfer实例完成读包的任务
		tf:=&utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端关闭了连接，服务器也退出这个连接")
				return err
			} else {
				fmt.Println("readPkg err:", err)
				return err
			}
		}
		err=this.serverProcessMes(&mes)
		if err!=nil{
			return err
		}
	}
}

//通过用户输入，选择需要执行的程序
func (this *Processor)serverProcessMes(mes *message.Message) (err error) {

	//看看是否能够接收到客户端发送的消息
	fmt.Println("客户端发来群发消息",mes)
	switch mes.Type {

	//每个case都要创建一个信息up实例，对应每一个上线的用户
	case message.LoginMesType:
		//处理的登录逻辑
		//创建一个UserProcess实例
		up:=&process2.UserProcess{
		Conn: this.Conn,
		}

		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		up:=&process2.UserProcess{
		Conn: this.Conn,
		}
		//缺少一个处理注册信息的函数
		err=up.ServerProcessRegister(mes)

	case message.SmsMesType:
		//处理群发消息
		smsProcess:=&process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)

	default:
		fmt.Println("消息类型不存在无法处理")

	}
	return
}

