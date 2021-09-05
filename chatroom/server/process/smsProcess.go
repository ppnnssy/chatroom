package process2

import (
	"chatroom/common/message"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
	//...
}

//转发消息
func (this *SmsProcess)SendGroupMes(mes *message.Message)  {
	//遍历服务器端的onlineUsers map[int]*UserProcess
	var smsMes message.SmsMes
	err:=json.Unmarshal([]byte(mes.Data),&smsMes)
	if err!=nil{
		fmt.Println("反序列化失败1",err)
		return
	}

	data,err:=json.Marshal(mes)

	for id,up:=range userMgr.onlineUsers{
		if id==smsMes.UserId{
			continue
		}
		this.SendMesToEachOnlineUser(data,up.Conn)
	}

}

func (this *SmsProcess)SendMesToEachOnlineUser(data []byte,conn net.Conn)  {
	tf:=&utils.Transfer{
		Conn: conn,
	}
	err:=tf.WritePkg(data)
	if err!=nil{
		fmt.Println("转发消息失败：",err)
	}
}