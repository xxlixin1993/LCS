package protocol

import (
	"time"
	"io"
	"github.com/xxlixin1993/LCS/logging"
	"github.com/xxlixin1993/LCS/utils"
)

const (
	kC0S0Version = 3
	kC1S1Zero = 0
)
var (
	timeout = 5 * time.Second
)


func (conn *Conn) HandshakeServer() (err error) {
	var allData [(1 + 1536*2) * 2]byte

	C0C1C2 := allData[:1536*2+1]
	C0 := C0C1C2[:1]
	C1 := C0C1C2[1 : 1536+1]
	C0C1 := C0C1C2[:1536+1]
	C2 := C0C1C2[1536+1:]

	S0S1S2 := allData[1536*2+1:]
	S0 := S0S1S2[:1]
	S1 := S0S1S2[1 : 1536+1]
	S2 := S0S1S2[1536+1:]

	// 读取C0 C1
	conn.Conn.SetDeadline(time.Now().Add(timeout))
	if _, err = io.ReadFull(conn.rw, C0C1); err != nil {
		return err
	}

	conn.Conn.SetDeadline(time.Now().Add(timeout))
	if C0[0] != kC0S0Version {
		logging.ErrorF("rtmp: handshake version=%d invalid", C0[0])
		return nil
	}

	// 读取C1 time
	clientTime := U32BigEndian(C1[0:4])
	logging.InfoF("client time is %d", clientTime)
	/*
	// zero 保留值可做特殊处理
	zero := U32BigEndian(C1[4:8])
	if zero != 0 {
		// 保留值zero 如果是默认0 则可以使用默认格式 C2和S2包几乎分别是S1和C1的复制 否则需要自定义校验生成S1 S2
	} else {
		copy(S1, C2)
		copy(S2, C1)
	}*/
	serverTimeC1 := uint32(utils.GetTimestamp())

	// 设置S0
	S0[0] = kC0S0Version
	// 设置S1
	copy(S1[0:4], C1[0:4])
	PutU32BigEndian(S1[0:4], serverTimeC1)
	PutU32BigEndian(S1[4:8], kC1S1Zero)
	copy(S1[8:1536], C1[8:1536])
	// 设置S2
	PutU32BigEndian(S2[0:4], serverTimeC1)
	PutU32BigEndian(S2[4:8], uint32(utils.GetTimestamp()))
	copy(S2[8:1536], C1[8:1536])


	// 发送 S0S1S2
	conn.Conn.SetDeadline(time.Now().Add(timeout))
	if _, err = conn.rw.Write(S0S1S2); err != nil {
		return err
	}
	conn.Conn.SetDeadline(time.Now().Add(timeout))
	if err = conn.rw.Flush(); err != nil {
		return err
	}

	// 接收 C2
	conn.Conn.SetDeadline(time.Now().Add(timeout))
	if _, err = io.ReadFull(conn.rw, C2); err != nil {
		return err
	}
	conn.Conn.SetDeadline(time.Time{})
	return
}


// RTMP 默认都是使用 Big-Endian 进行写入和读取，除非强调对某个字段使用 Little-Endian 字节序。
// Big-Endian 将高序字节存储在起始地址（高位编址）
// 获取
func U32BigEndian(b []byte) (i uint32) {
	i = uint32(b[0])
	i <<= 8
	i |= uint32(b[1])
	i <<= 8
	i |= uint32(b[2])
	i <<= 8
	i |= uint32(b[3])
	return
}

// 设置 Big-Endian
func PutU32BigEndian(b []byte, v uint32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}
