package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gozinx/utils"
	"gozinx/ziface"
)

type DataPack struct{}

func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

// 获取包头的长度
func (d *DataPack) GetHeadLen() uint32 {
	// 包头长度固定为8字节，数据长度4字节+ID长度4字节
	return 8
}

// 封包
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放字节流的缓冲
	databuff := bytes.NewBuffer([]byte{})

	// 将dataLen写入databuff
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 将Id写入databuff
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将数据写入databuff
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return databuff.Bytes(), nil
}

// 拆包，将包的头信息读出来，之后根据头信息中的数据长度再读一次
func (d *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	// 创建一个从输入读取二进制数据的ioReader
	databuff := bytes.NewReader(data)

	// 解压头信息，获取数据长度和msgID
	msg := &Message{}
	if err := binary.Read(databuff, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}
	if err := binary.Read(databuff, binary.LittleEndian, &msg.MsgId); err != nil {
		return nil, err
	}

	// 判断数据长度是否超过MaxPackageSize
	if msg.Len > utils.GloablObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg, nil
}
