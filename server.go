package main

/**
Server类
**/

import (
	"fmt"
	"net"
	"sync"
)

// Server结构体
type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User //在线用户列表（因为在线用户表具有全局属性，需要加锁）
	mapLock   sync.RWMutex     //读写锁
	Message   chan string      //消息广播的channle
}

// new一个server对象
func NewServer(ip string, port int) *Server {

	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message消息管道，一旦有消息就发送给全部在线用户
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//将消息发送给全部在线用户
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg //将消息发送到用户消息管道
		}
		this.mapLock.Unlock()
	}
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

// Handler业务方法
func (this *Server) handler(conn net.Conn) {

	user := NewUser(conn)

	//用户上线，将用户加入到nolineMap中
	this.mapLock.Lock() //加锁
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock() //释放锁

	//广播当前用户上线消息
	this.BroadCast(user, "已上线")

	//阻塞hanler
	select {}
}

// 启动Server
func (this *Server) Start() {

	//创建socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// 关闭 listen socket
	defer listener.Close()

	//启动监听Message属性的goroutine
	go this.ListenMessager()

	//阻塞等待接收
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accpet err:", err)
			continue
		}
		//do handler，返回的链接发送到业务函数
		go this.handler(conn)
	}

}
