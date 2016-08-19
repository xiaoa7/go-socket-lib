package main

import (
	"fmt"
	"util"
)

var s *util.Server

//
func main() {
	s, _ = util.NewServer(":8080", func(p *util.PacketData) {
		fmt.Println(p.Act, string(p.Data))
		s.Write(p.SourceId, 20, []byte("服务器的回应内容"))
	}, func(id string) {
		fmt.Println("用户", id, "掉线了")
	})
	//看所有用户
	//s.GetAllClient()
	//判断用户是否在线
	//s.IsOnline(id)
	defer s.Close()
	lock := make(chan bool)
	<-lock
}
