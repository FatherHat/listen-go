package main

import "net"

/*
*
用户类
*
*/
type User struct {
	Name   string      //名称
	Addr   string      //ip地址
	C      chan string //消息管道
	conn   net.Conn    //客户端链接
	server *Server     //用户所在的服务器
}

// new一个用户对象
func NewUser(conn net.Conn, server *Server) *User {
	//获取客户端ip地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//启动监听当前user的消息管道goroutine
	go user.ListenMessage()

	return user
}

// 用户上线
func (this *User) OnLine() {

	// 用户上线，将用户加入到nolineMap中
	this.server.mapLock.Lock() //加锁
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock() //释放锁

	// 广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}

// 用户下线
func (this *User) OffLine() {
	// 用户下线，将用户从nolineMap
	this.server.mapLock.Lock() //加锁移除
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock() //释放锁

	// 广播当前用户上线消息
	this.server.BroadCast(this, "已下线")
}

// 处理消息业务
func (this *User) DoMessage(msg string) {

	//广播消息
	this.server.BroadCast(this, msg)

}

// 监听当前user的消息channel
func (this *User) ListenMessage() {
	//持续监听
	for {
		msg := <-this.C                     //从管道中获取消息
		this.conn.Write([]byte(msg + "\n")) //发送byte类型数组
	}
}
