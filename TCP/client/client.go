package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err1 := net.Dial("tcp", "192.168.1.5:8889")
	if err1 != nil {
		fmt.Println("出现错误：", err1)
	}
	fmt.Println("conn成功=", conn)
	//功能1：客户端发送单行数据，然后退出


	for {
		reader := bufio.NewReader(os.Stdin) //os.Stdin 代表终端的标准输入
		//从终端读取输入并发送给服务器
		line, err2 := reader.ReadString('\n')
		if err2 != nil {
			fmt.Println("出错了，err2=", err2)
		}
		//如果用户输入的是exit就退出
		line=strings.Trim(line,"\r\n")
		if line=="exit"{
			break
		}

		//发送给服务器
		_,err3 := conn.Write([]byte(line))
		if err3 != nil {
			fmt.Println("出现错误，err3:", err3)
		}
		//fmt.Printf("给服务端发送了%d字节的数据,并退出。", n)
	}
}
