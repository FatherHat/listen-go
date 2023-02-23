package main

func main() {
	//new一个服务对象
	server := NewServer("127.0.0.1", 8888)
	//启动服务对象，监听一个端口
	server.Start()
}
