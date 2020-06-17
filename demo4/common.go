package demo4

import (
	"encoding/binary"
	"io"
	"sync"
)

type connRW struct {
	mut sync.Mutex
}

var ConnRW = &connRW{}

func (c *connRW) WriteTo(w io.Writer, body []byte) (int64, error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	var total int64 = -1
	var buf [4]byte
	bufs := buf[:]
	binary.BigEndian.PutUint32(bufs, uint32(len(body)))
	n, err := w.Write(bufs)
	total += int64(n)
	n, err = w.Write(body)
	total += int64(n)
	if err != nil {
		return total, err
	}
	return total, nil
}

// 4 字节的数据长度+具体数据
func (c *connRW) ReadResponse(r io.Reader) ([]byte, error) {
	var msgSize int32
	// message size
	err := binary.Read(r, binary.BigEndian, &msgSize)
	if err != nil {
		return nil, err
	}
	// message binary data
	buf := make([]byte, msgSize)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
