package process

import (
	"chatroom/common/message"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

//发送群聊消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	//1.创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType

	//创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	//序列化
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("群发消息序列化失败1")
		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("群发消息序列化失败2")
		return
	}

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("群发消息失败")
	}
	return
}
