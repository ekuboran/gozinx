package ziface

type IMessage interface {
	GetMsgId() uint32
	GetData() []byte
	GetDataLen() uint32
	SetMsgId(id uint32)
	SetData(data []byte)
	SetDataLen(id uint32)
}
