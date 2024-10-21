package utils

import (
	"bytes"
	"encoding/gob"
)

type request_type byte

const (
	AGENT_REGISTRATION request_type = iota // iota = 0
	TASK_REQUEST                           // iota = 1
	METRICS_GATHERING                      // iota = 2
	ERROR                                  // iota = 3
)

type Ack struct {
	Acknowledged     bool
	Sender_id        byte // [0, 255]
	Sender_is_server bool
	Request_id       byte         // [0, 255]
	Request_type     request_type //(AGENT_REGISTRATION, TASK_DELEGATION, DATA_COLLECTION)
}

type ack_builder struct {
	Ack Ack
}

func New_ack_builder() *ack_builder {
	return &ack_builder{
		Ack: Ack{
			Acknowledged:     false,
			Sender_id:        0,
			Sender_is_server: false,
			Request_id:       0,
			Request_type:     ERROR},
	}
}

func (a *ack_builder) Has_ackowledged() *ack_builder {
	a.Ack.Acknowledged = true
	return a
}

func (a *ack_builder) Set_sender_id(id byte) *ack_builder {
	a.Ack.Sender_id = id
	return a
}

func (a *ack_builder) Is_server() *ack_builder {
	a.Ack.Sender_is_server = true
	return a
}

func (a *ack_builder) Set_request_id(id byte) *ack_builder {
	a.Ack.Request_id = id
	return a
}

func (a *ack_builder) Set_request_type(request request_type) *ack_builder {
	a.Ack.Request_type = request
	return a
}

func (a *ack_builder) Build() Ack {
	return a.Ack
}

func Decode_ack(message []byte) (Ack, error) {
	var ack Ack
	buffer := bytes.NewBuffer(message)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&ack)
	return ack, err
}

func Encode_ack(ack Ack) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ack)
	if err != nil { // porquê if neste caso de Encode mas não no Decode? Não é possível dar return de um array de bytes que não foi "carregado" com nada?
		return nil, err
	}
	return buffer.Bytes(), nil
}
