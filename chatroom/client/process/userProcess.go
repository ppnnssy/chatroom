package process

import (
	"chatroom/common/message"
	"chatroom/server/utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
	//暂时不需要任何东西
}


//登录函数
func(this *UserProcess)Login(userId int, userPwd string) (err error) {
	//fmt.Printf("userId:%d,userPwd:%s\n", userId, userPwd)
	//return nil

	//1.连接服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}
	//延迟关闭
	defer conn.Close()
	//2.准备通过conn发送消息给服务器
	var mes message.Message         //创建一个Message结构体的对象
	mes.Type = message.LoginMesType //Type是LoginMesType = "loginMes"

	//3.创建一个loginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4.对loginMes进行序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}

	//5.把data赋值给mes.Data字段
	mes.Data = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.mes err:", err)
		return
	}

	//7.现在可以发送消息了,data就是需要发送的消息
	//7.1 先把data的长度发送给服务器
	//先获取data的长度，再转成表示长度的切片。因为write只能发送切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[0:4], pkgLen) //把pkglen放进切片中

	//发送长度
	n, err := conn.Write(bytes[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes[0:4])失败：", err)
		return
	}
	fmt.Println("客户端发送消息长度成功,长度是：", len(data), string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data)失败：", err)
		return
	}
	fmt.Println("发送消息本身成功")



	tf:=&utils.Transfer{
		Conn: conn,
	}
	//处理服务端返回的消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg（）出错：", err)
	}

	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	fmt.Println(err)

	if loginResMes.Code==200{
		fmt.Println("登录成功")
		//初始化CurUser
		CurUser.Conn=conn
		CurUser.UserId=userId
		CurUser.UserStatus=message.UserOnline

		//显示在线用户列表
		fmt.Println("当前在线用户列表如下")
		for _,v:=range loginResMes.UserIds{
			if v==userId{
				continue
			}
			fmt.Println("用户id：\t",v)
			user:=&message.User{
				UserId: v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v]=user  //onlineUsers是一个map类型全局变量。var onlineUsers map[int]*message.User
		}
		fmt.Print("\n\n")

		//这里应该启动一个协程
		//该协程保持和服务端的通讯，如果服务器有消息推送给客户端，则接收并显示在客户端的终端
		go serverProcessMes(conn)

		//1.显示登录成功后的菜单
		for{
			ShowMenu()
		}
	}else {
		fmt.Println(loginResMes.Error)
		fmt.Println("登录失败")
	}
	return
}




//注册函数
func (this *UserProcess)Register (userId int,userPwd string,userName string)(err error) {
	//1.连接服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}
	//延迟关闭
	defer conn.Close()
	//2.准备通过conn发送消息给服务器
	var mes message.Message         //创建一个Message结构体的对象
	mes.Type = message.RegisterMesType //Type是RegisterResMesType="RegisterResMes"
	//创建一个LoginMes的结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId=userId
	registerMes.User.UserPwd=userPwd
	registerMes.User.UserName=userName

	//将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}

	//5.把data赋值给mes.Data字段
	mes.Data = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.mes err:", err)
		return
	}


	//利用Write函数发送消息
	tf:=&utils.Transfer{
		Conn: conn,
	}
	err=tf.WritePkg(data)
	if err!=nil{
		fmt.Println("注册发送信息错误：",err)
	}
	//处理服务器返回的登录成功与否的消息
	//这里出现了阻塞，服务端没有返回消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("注册返回信息错误", err)
	}



	//以下是处理返回的信息
	//将返回的信息序列化
	var RegisterResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &RegisterResMes)
	if RegisterResMes.Code==200 {
		fmt.Println("注册成功，可以重新登录试一下")
		os.Exit(0)
	}else {
		fmt.Println(RegisterResMes.Error)
		os.Exit(0)
	}

	return
	}