package model

import (
	"chatroom/common/message"
	"net"
)

type CurUser struct {
	Conn net.Conn
	message.User
}
