package utils

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)
//将这些方法关联到结构体中
type Transfer struct {
	//分析它应该有哪些字段
	Conn net.Conn
	Buf [8096]byte //传输时使用的缓冲

}

func (this *Transfer)WritePkg( data []byte) (err error) {
	//先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[0:4], pkgLen) //把pkglen放进切片中

	//发送长度
	n, err := this.Conn.Write(bytes[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes[0:4])失败：", err)
		return
	}

	//发送data本身
	n, err =this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes[0:4])失败：", err)
		return
	}
	return
}

func (this *Transfer)ReadPkg() (mes message.Message, err error) {
	//将读取数据包封装成一个函数readPkg
	fmt.Println("读取客户端发送的数据")

	//conn.read在conn没有关闭的情况下才会阻塞
	//如果客户端关闭了conn就不会阻塞

	//这里读取消息头，内容是消息长度
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}

	//记录消息长度
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])
	//根据长度  读取内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}

	//把消息反序列化
	err = json.Unmarshal(this.Buf[:pkgLen], &mes) //重要的坑，&mes使用指针
	if err != nil {
		fmt.Println("反序列化失败")
		return
	}
	return
}