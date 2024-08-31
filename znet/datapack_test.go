package znet

/*
datapack模块单元测试
*/

import (
	"gozinx/utils"
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {

	// 模拟服务端
	listener, err := net.Listen("tcp4", "127.0.0.1:7777")
	if err != nil {
		utils.DPrintln("server listen err: ", err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				utils.DPrintln("server Accept err: ", err)
				continue
			}

			go func(conn net.Conn) {
				// 模拟拆包
				dp := NewDataPack()
				for {
					// 第一次读头信息
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						utils.DPrintln("server read head err: ", err)
						return
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						utils.DPrintln("server unpack err: ", err)
						return
					}
					if msgHead.GetDataLen() > 0 {
						// 如果有数据，进行第二次读取
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							utils.DPrintln("server read data err: ", err)
							return
						}
						utils.DPrintln("server recv: msglen=", msg.Len, " msgId=", msg.MsgId, " data=", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp4", "127.0.0.1:7777")
	if err != nil {
		utils.DPrintln("client dial err: ", err)
		return
	}
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一起发送
	// 封装msg1包
	msg1 := &Message{
		MsgId: 1,
		Len:   4,
		Data:  []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		utils.DPrintln("client pack err: ", err)
		return
	}

	// 封装msg2包
	msg2 := &Message{
		MsgId: 2,
		Len:   5,
		Data:  []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		utils.DPrintln("client pack err: ", err)
		return
	}

	// 两个包 粘在一起
	sendData1 = append(sendData1, sendData2...)

	conn.Write(sendData1)

	time.Sleep(3 * time.Second)
}
