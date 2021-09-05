package process

import (
	"chatroom/common/message"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//显示登录成功后的界面
func ShowMenu() {
	fmt.Println("---------------恭喜***登录成功-------------")
	fmt.Println("---------------1.显示在线用户列表-------------")
	fmt.Println("---------------2.发送消息-------------")
	fmt.Println("---------------3.信息列表-------------")
	fmt.Println("---------------4.退出系统-------------")
	fmt.Println("请选择（1-4）：")
	var key int
	var content string
	//因为总会使用到发送消息，所以把smsProcess定义在外面
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		outPutOnlineUsers()
	case 2:
		fmt.Println("请输入你想群发的消息")
		fmt.Scanf("%s\n",&content)
		smsProcess.SendGroupMes(content)

	case 3:
		fmt.Println("信息列表")

	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("输入有误，请重新输入")
	}

}

//和服务器端保持通信
func serverProcessMes(conn net.Conn) {
	//创建一个transfer实例，不停读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("服务器端出错了", err)
			return
		}
		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线了
			//1.取出NotifyUserStatusMes并反序列化
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)

			//2.把这个用户的状态信息保存到客户端的map[int]User中
			updateUserStatus(&notifyUserStatusMes)

		//1.取出NotifyUserStatusMes
		//2.把这个用户的状态信息保存到客户端的map中
		case message.SmsMesType:
			outPutGroupMes(&mes)

		default:
			fmt.Println("服务器端返回了暂时不能处理的消息类型")
		}

		//读取到消息没出错，就是下一步逻辑。暂时写成输出
		//fmt.Println("系统发来消息：", mes)
	}
}
