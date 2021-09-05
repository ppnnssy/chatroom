package main

import (
	"fmt"
	"net"
)

//
func process(conn net.Conn)  {
	//延迟关闭conn
	defer conn.Close()
	//循环读取数据
	for  {
		//创建一个新的切片
		buf:=make([]byte,1024)
		//对方不发信息，就一直在这等，协程阻塞。有可能超时
		//fmt.Println("服务器在等待客户端"+conn.RemoteAddr().String()+"发送信息")
		n,err:=conn.Read(buf) //从conn读取信息
		if err!=nil{
			fmt.Println("读取数据出错：",err)
			return
		}
		fmt.Printf("%s",string(buf[:n])+"\n")
	}
}




func main() {

	fmt.Println("服务器开始监听。。。")
	//tcp表示使用的网络协议是tcp
	//0.0.0.0:8888表示在本地监听8888端口
	listen,err:=net.Listen("tcp","0.0.0.0:8889")
	defer listen.Close()//延迟关闭
	if err !=nil{
		fmt.Println("出现错误：",err)
		return
	}

	//循环等待客户端来连接
	for  {
		conn,err:=listen.Accept()
		fmt.Println("等待客户端来连接。。。")
		if err!=nil{
			fmt.Println("listen.Accept出错",err)
		}else {
			fmt.Printf("Accept()suc con = %v,客户端IP：%v\n",conn,conn.RemoteAddr().String())
		}
		//这里起个协程，为客户端服务
		go process(conn)

		}

}
