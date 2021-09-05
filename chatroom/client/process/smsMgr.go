package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

func outPutGroupMes(mes *message.Message) {//这里mes.data里存的是smsMes
	//显示即可
	//1.反序列化mes.Data
	var smsMes message.SmsMes
	err:=json.Unmarshal([]byte(mes.Data),&smsMes)
	if err!=nil{
		fmt.Println("反序列化失败2",err)
		return
	}

	//显示消息
	info:=fmt.Sprintf("用户ID\t：%d对大家说：\t%s",
		smsMes.UserId,smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
