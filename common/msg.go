package common

import (
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/api_tcp_demo/protos"
)

func WriteTo(w io.Writer, msg *protos.Msg) (int64, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return -1, err
	}
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
func ReadResponse(r io.Reader) (*protos.Msg, error) {
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
	smg := &protos.Msg{}
	return smg, json.Unmarshal(buf, smg)
}
