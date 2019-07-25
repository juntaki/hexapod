package common

import (
	"encoding/json"
	"time"
)

type MessageType int

const (
	MessageTypeArms MessageType = iota
	MessageTypeWalk
	MessageTypeRotate
	MessageTypeDemo
	MessageTypeReset
)

type Message struct {
	Now time.Time
	Sequence int
	MessageType MessageType
	Arms []*Arm
}

type Arm struct {
	Degrees []float64
}

func (m *Message) Message()[]byte {
	ret, _ := json.Marshal(m)
	return ret
}

func UnmarshalMessage(mes []byte) (*Message, error){
	var ret Message
	err := json.Unmarshal(mes, &ret)
	if err != nil {
		return nil, err
	}
	return &ret ,nil
}