package main

import "net"

type User struct {
	Name string
	Ip string
	C chan int
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	//获取当前客户的连接地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:userAddr,
		Ip:userAddr,
		C:make(chan stirng),
		conn:conn,
	}
	return user
}
