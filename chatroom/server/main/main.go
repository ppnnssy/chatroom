package main

import (
	"chatroom/server/model"
	"fmt"
	"net"
	"time"
)


func main() {
	//当服务器启动时初始化redis的连接池
	initPool("localhost:6379",16,0,300*time.Second)
	//初始化一个UserDao
	initUserDao()
	//提示信息
	fmt.Println("服务器[新的结构]在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("监听出错", err)
		return
	}

	for {
		fmt.Println("等待客户端连接服务器。。。")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("连接出错：", err)
		}

		//一旦连接成功，则启动一个协程和客户端保持通信。。。
		go process(conn)
	}
}

//编写一个函数，完成对UserDao的初始化任务
func initUserDao() {
	model.MyUserDao=model.NewUserDao(pool) //这里的pool是redis.go中定义的全局变量，需要先初始化pool
}


func process(conn net.Conn) {
	//需要延迟关闭conn
	defer conn.Close()

	//创建一个总控实例进行调用
	processor := &Processor{
		Conn: conn,
	}

	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务端通讯的协程出错：", err)
		return
	}
}
