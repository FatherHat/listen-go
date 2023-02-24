package main

/**
Server类
**/

import (
	"fmt"
	"io"
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

	user := NewUser(conn, this)

	//用户上线
	user.OffLine()

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//用户下线
				user.OffLine()
				return
			}
			//io.EOF表示文件的末尾，读到文件的末尾是非法的操作
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
			}

			//提取用户的消息（去除'\n'）,:号表示读取范围
			msg := string(buf[:n-1])

			//用户针对msg处理
			user.DoMessage(msg)
		}
	}()

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
