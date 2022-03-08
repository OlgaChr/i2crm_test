package framesocket

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/gorilla/websocket"
)

type FrameSocket struct {
	wsConn *websocket.Conn
}

func NewFrameSocket(ws *websocket.Conn) *FrameSocket {
	return &FrameSocket{
		wsConn: ws,
	}
}

func (fs *FrameSocket) Read() ([]byte, error) {
	messageType, msg, err := fs.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if messageType != websocket.BinaryMessage {
		return nil, fmt.Errorf("reading - get not binary type")
	}
	if len(msg) < 3 {
		return nil, fmt.Errorf("reading - too small message length")
	}
	// get body from message
	l := msg[:3]
	length := binary.BigEndian.Uint32(append(make([]byte, 1), l...))
	return msg[3 : length+3], err
}

func (fs *FrameSocket) Write(data []byte) error {
	l := uint32(len(data))
	if l > 2^24 {
		return fmt.Errorf("writing - too big input data")
	}

	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, uint32(l))
	bs := buff.Bytes()
	header := bs[len(bs)-3:]

	sendData := append(header, data...)
	return fs.wsConn.WriteMessage(websocket.BinaryMessage, sendData)
}

func (fs *FrameSocket) Close() error {
	return fs.wsConn.Close()
}
