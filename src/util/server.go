//服务端封装
//支持以下功能：在线管理，心跳同步，
package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//
type Server struct {
	listen          net.Listener
	online          map[string]*net.Conn //在线
	businesshandler func(p *PacketData)  //业务事件处理
	lock            sync.Mutex           //锁
}

//
func NewServer(addr string, handler func(p *PacketData)) (s *Server, err error) {
	s = new(Server)
	s.listen, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}
	s.online = make(map[string]*net.Conn)
	s.businesshandler = handler
	msgdump := make(chan *PacketData, 20)
	go s.processBusinessHandler(msgdump)
	go func() { //接受外部连接
		for {
			conn, err := s.listen.Accept()
			if err != nil {
				continue
			} else {
				go s.handler(msgdump, &conn)
			}
		}
	}()
	return
}

//
func (s *Server) Close() {
	s.listen.Close()
}

//处理连接
func (s *Server) handler(pc chan<- *PacketData, conn *net.Conn) {
	//请求客户端返回个人ID
	(*conn).Write(EnPacket(ACT_REQUEST_CLIENTID, []byte{}))
	var clientId string
	//接收消息
	for {
		head := make([]byte, PACKETHEADLEN)
		r, _ := io.ReadFull(*conn, head)
		if r != PACKETHEADLEN {
			//读取头长度不对
			s.closeClient(conn)
			return
		} else if !bytes.Equal(DEFAULT_HEAD, head[:HEADLEN]) {
			//头标志不对
			s.closeClient(conn)
			return
		} else {
			size := byte2int(head[HEADLEN:])
			data := make([]byte, size)
			r, _ := io.ReadFull(*conn, data)
			if r != size {
				//取到的数据长度不对
				s.closeClient(conn)
				return
			} else {
				//数据正常了，识别下ACT
				act, data := Parse(data)
				switch act {
				case ACT_RESPONSE_CLIENTID: //加入
					clientId = string(data)
					s.join(clientId, conn)
				case ACT_REQUEST_HEARTBEAT: //心跳回应
					s.processheartbeat(clientId, data, conn)
				default:
					pc <- &PacketData{Act: act, Data: data, SourceId: clientId}
				}
			}
		}

	}
}

//
func (s *Server) processBusinessHandler(pc <-chan *PacketData) {
	for {
		select {
		case p := <-pc:
			go s.businesshandler(p)
		}
	}
}

//
func (s *Server) join(id string, conn *net.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if v, ok := s.online[id]; ok {
		(*v).Close()
	}
	fmt.Println(id, "加入")
	s.online[id] = conn
}

//心跳检查，发送心跳包
func (s *Server) checkheartbeat() {
	now := time.Now().Unix()
	for _, v := range s.online {
		(*v).Write(EnPacket(ACT_REQUEST_HEARTBEAT, int2byte(int(now))))
	}
	time.AfterFunc(5*time.Second, s.checkheartbeat)
}

//处理心跳回应
func (s *Server) processheartbeat(clientId string, data []byte, conn *net.Conn) {
	fmt.Println(clientId, "回应心跳")
	now := time.Now().Unix()
	cnow := int64(byte2int(data))
	if cnow > now-10 { //超时10秒开始清理
		s.closeClient(conn)
	}
}

//对制定客户端发消息
func (s *Server) Write(clientId string, act int, data []byte) error {
	if v, ok := s.online[clientId]; ok {
		(*v).Write(EnPacket(act, data))
		return nil
	} else {
		return errors.New("客户端未连接")
	}
}

//关闭客户端
func (s *Server) closeClient(conn *net.Conn) {
	s.lock.Lock()
	for k, v := range s.online {
		if v == conn {
			delete(s.online, k)
			break
		}
	}
	s.lock.Unlock()
	(*conn).Close()
}
