package main

import (
	"fmt"
	"net"
)

// 创建一个结构体
type Server struct {
	Ip   string
	Port int
}

// new一个server对象，然后返回
func NewServer(ip string, port int) *Server {

	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (this *Server) handler(conn net.Conn) {
	//...当前连接的业务
	fmt.Println("链接建立成功")
}

// 启动服务器的接口(套接字
func (this *Server) Start() {

	//1、创建socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//2、阻塞等待接收
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accpet err:", err)
			continue
		}
		//3、do handler，返回的链接发送到业务函数
		go this.handler(conn)
	}

	//4、关闭socket
	defer listener.Close()
}
