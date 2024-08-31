package znet

type Message struct {
	MsgId uint32
	Data  []byte
	Len   uint32
}

func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) GetDataLen() uint32 {
	return m.Len
}

func (m *Message) SetMsgId(id uint32) {
	m.MsgId = id
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(len uint32) {
	m.Len = len
}

func NewMessage(msgId uint32, data []byte) *Message {
	msg := &Message{
		MsgId: msgId,
		Data:  data,
		Len:   uint32(len(data)),
	}
	return msg
}
