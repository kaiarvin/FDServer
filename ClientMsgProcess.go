// ClientMsgProcess
package ChatServer

import (
	"fmt"
	"io"
)

func (this *Client) ClientMsgProcess() {
	fmt.Println("In ClientMsgProcess...")
	select {
	case buf := <-this.RecMsg:
		{
			fmt.Println("In <-this.RecMsg...")
			fmt.Println(buf)
			//this.SendToClient(buf)
			fmt.Println("Out <-this.RecMsg...")
		}
	case data := <-this.AckMsg:
		{
			fmt.Println("In <-this.AckMsg...")
			fmt.Println("data:", data)
			buf := []byte(data)
			n, err := this.Client_Socket.Write(buf)
			if err != nil {
				if err != io.EOF {
					return
				}
			}
			fmt.Println("Ack len: ", n)
			if n != 0 {
				this.ChatServer.AckDataSize += uint64(n)
			}
			fmt.Println("In <-this.AckMsg...")
		}
	}
}

func (this *Client) SendToSelfClient(buf string) {
	this.AckMsg <- buf
	this.ClientMsgProcess()
}
