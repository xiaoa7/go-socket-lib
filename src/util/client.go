//客户端的常用操作封装
package util

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct {
	id   string
	conn net.Conn
}

//创建
func NewClient(serveraddr string, myid string, datahandler func(p *PacketData)) *Client {
	c := &Client{id: myid}
	c.conn, _ = Dial("tcp", serveraddr)
	msgdump := make(chan *PacketData, 20)
	go c.read(msgdump)
	go c.handler(msgdump, datahandler)
	return c
}

//写入处理
func (c *Client) Write(act int, data []byte) error {
	_, err := c.conn.Write(EnPacket(act, data))
	return err
}

//
func (c *Client) read(pc chan<- *PacketData) {
	for {
		head := make([]byte, PACKETHEADLEN)
		r, _ := io.ReadFull(c.conn, head)
		if r != PACKETHEADLEN {
			//读取头长度不对
			fmt.Print("?")
			time.Sleep(2 * time.Second)
			continue
		} else if !bytes.Equal(DEFAULT_HEAD, head[:HEADLEN]) {
			fmt.Print(">")
			time.Sleep(2 * time.Second)
			//头标志不对
			continue
		} else {
			size := byte2int(head[HEADLEN:])
			data := make([]byte, size)
			r, _ := io.ReadFull(c.conn, data)
			if r != size {
				fmt.Print("/")
				time.Sleep(2 * time.Second)
				//取到的数据长度不对
				continue
			} else {
				act, bdata := Parse(data)
				switch act {
				case ACT_REQUEST_CLIENTID:
					c.conn.Write(EnPacket(ACT_RESPONSE_CLIENTID, []byte(c.id)))
				case ACT_REQUEST_HEARTBEAT:
					fmt.Print(".")
					c.conn.Write(EnPacket(ACT_RESPONSE_HEARTBEAT, bdata))
				default:
					pc <- NewPacketDate(act, bdata)
				}

			}
		}

	}
}

//
func (c *Client) handler(pc <-chan *PacketData, fn func(p *PacketData)) {
	for {
		select {
		case p := <-pc:
			go fn(p)
		}
	}
}
