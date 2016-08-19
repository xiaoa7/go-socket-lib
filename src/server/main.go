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
	})
	defer s.Close()
	lock := make(chan bool)
	<-lock
}
