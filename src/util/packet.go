//协议封装
package util

import (
	"bytes"
	"encoding/binary"
)

const (
	HEADLEN       = 4
	PACKETHEADLEN = 8 //头4位+数据长度4位
)

var DEFAULT_HEAD = []byte{0x1, 0x2, 0x3, 0x4} //数据包头

//
type PacketData struct {
	Act      int
	Data     []byte
	SourceId string //服务端使用
}

//
func NewPacketDate(act int, data []byte) *PacketData {
	return &PacketData{Act: act, Data: data}
}

//
func byte2int(src []byte) (ret int) {
	var tmp int32
	b_buf := bytes.NewBuffer(src)
	binary.Read(b_buf, binary.BigEndian, &tmp)
	ret = int(tmp)
	return
}

//
func int2byte(src int) []byte {
	var tmp int32
	tmp = int32(src)
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, tmp)
	return b_buf.Bytes()
}

//head+length+data
func (p *PacketData) Bytes() (ret []byte) {
	ret = append(DEFAULT_HEAD)
	ret = append(ret, int2byte(len(p.Data)+4)...)
	ret = append(ret, int2byte(p.Act)...)
	ret = append(ret, p.Data...)
	return
}

//
func EnPacket(act /*指令*/ int, bs /*实际数据*/ []byte) (ret []byte) {
	ret = append(DEFAULT_HEAD)
	ret = append(ret, int2byte(len(bs)+4)...)
	ret = append(ret, int2byte(act)...)
	ret = append(ret, bs...)
	return
}

//解析数据
func Parse(bs []byte) (act int, ret []byte) {
	act = byte2int(bs[:4])
	ret = bs[4:]
	return
}
