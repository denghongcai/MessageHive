// 通信协议模块
package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/denghongcai/MessageHive/modules/message"
	"github.com/golang/protobuf/proto"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

// Packet最大包长度定义
const (
	MAX_PACKET_LENGTH uint32 = 1 << 20 // 1M
)

// Packet结构
type Packet struct {
	length uint32
	Msg    *message.Container
}

// Packet打包方法
func (pkt *Packet) Pack() []byte {
	msgBytes, _ := proto.Marshal(pkt.Msg)
	pkt.length = uint32(len(msgBytes))
	log.Debug("Packet length: %d bytes", pkt.length)
	return append(UInt32ToBytes(pkt.length), msgBytes...)
}

// Packet解包方法
func (pkt *Packet) UnPack(b *[]byte) (bool, error) {
	if len(*b) < 4 {
		return false, nil
	}
	pkt.length = binary.BigEndian.Uint32((*b)[:4])
	if pkt.length > MAX_PACKET_LENGTH {
		return false, errors.New(fmt.Sprintf("Max packet length exceeded, current %d bytes", pkt.length))
	}
	if len(*b) < (4 + int(pkt.length)) {
		return false, nil
	}
	data := (*b)[4 : 4+pkt.length]
	pkt.Msg = new(message.Container)
	err := proto.Unmarshal(data, pkt.Msg)
	if err != nil {
		log.Debug("Msg: %q", data)
		log.Debug("Msg length: %d", pkt.length)
		return false, err
	}
	*b = (*b)[4+pkt.length:]
	return true, nil
}

// UInt32到Bytes大端序转换
func UInt32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}
