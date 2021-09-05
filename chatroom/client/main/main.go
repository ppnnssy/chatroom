package main

import (
	"chatroom/client/process"
	"fmt"
	"os"
)

//定义两个全局变量ID和密码
var userId int
var userPwd string
var userName string

func main() {
	//接收用户的选择
	var key int
	//判断是否继续显示菜单
	var loop =true
	for loop{
		fmt.Println("----------------欢迎登录多人聊天系统-------------------")
		fmt.Println("				   1.登录聊天系统")
		fmt.Println("				   2.注册用户")
		fmt.Println("				   3.退出登录")
		fmt.Println("请选择（1-3）：")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户ID：")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入密码：")
			fmt.Scanf("%s\n", &userPwd)
			//完成登录
			//1.创建一个Userprocess的实例
			up:=&process.UserProcess{}
			up.Login(userId,userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id：")
			fmt.Scanf("%d\n",&userId)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n",&userPwd)
			fmt.Println("请输入用户昵称：")
			fmt.Scanf("%s",&userName)
			//调用UserProcess实例，完成注册请求
			up:=&process.UserProcess{
			}
			up.Register(userId,userPwd,userName)

		case 3:
			fmt.Println("退出系统")
			//loop = false
			os.Exit(0)
		default:
			fmt.Println("输入有误，重新输入：")

		}

	}

}

