package main

import (
	"fmt"
	"time"
	"util"
)

//
func main() {
	c := util.NewClient(":8080", "myidisabc", func(p *util.PacketData) {
		fmt.Println(p.Act, string(p.Data))
	})
	time.Sleep(2 * time.Second)
	c.Write(10, []byte("服务器你好"))
	c.Write(10, []byte("服务器你好"))
	c.Write(10, []byte("服务器你好"))
	lock := make(chan bool)
	<-lock
}
