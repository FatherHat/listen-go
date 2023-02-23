package main

import "net"

/*
*
用户类
*
*/
type User struct {
	Name string      //名称
	Addr string      //ip地址
	C    chan string //消息管道
	conn net.Conn    //客户端链接
}

// new一个用户对象
func NewUser(conn net.Conn) *User {
	//获取客户端ip地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	//启动监听当前user的消息管道goroutine
	go user.ListenMessage()

	return user
}

// 监听当前user的消息channel
func (this *User) ListenMessage() {
	//持续监听
	for {
		msg := <-this.C                     //从管道中获取消息
		this.conn.Write([]byte(msg + "\n")) //发送byte类型数组
	}
}
