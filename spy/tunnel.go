package spy

import (
	"encoding/binary"
	"github.com/go-well/spider/silk"
	"net"
	"strings"
	"sync"
)

var tunnels sync.Map
var tunnelIndex uint16 = 1

//TODO 封装成Tunnel类
func receiveTunnel(c *Client, conn net.Conn, id uint16) {
	buf := make([]byte, 512)
	binary.BigEndian.PutUint16(buf, id)
	for {
		n, e := conn.Read(buf[2:])
		if e != nil {
			_ = c.Send(&silk.Package{
				Type: silk.TunnelError,
			})
			break
		}
		data := buf[:2+n]
		pack := &silk.Package{
			Type: silk.TunnelData,
			Data: data,
		}
		_ = c.Send(pack)
	}
	_ = c.Send(&silk.Package{
		Type: silk.TunnelEnd,
	})

	tunnels.Delete(id)
}

func init() {

	RegisterHandler(silk.TunnelCreate, func(c *Client, p *silk.Package) {
		p.Type = silk.TunnelCreateAck

		name := string(p.Data)
		names := strings.Split(name, ",")

		conn, err := net.Dial(names[0], names[1])
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}

		go receiveTunnel(c, conn, tunnelIndex)

		//缓存
		tunnels.Store(tunnelIndex, conn)
		p.Data = make([]byte, 2)
		binary.BigEndian.PutUint16(p.Data, tunnelIndex)
		tunnelIndex++

		_ = c.Send(p)
	})

	RegisterHandler(silk.TunnelClose, func(c *Client, p *silk.Package) {
		p.Type = silk.TunnelCloseAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := tunnels.LoadAndDelete(id)
		if !ok {
			p.SetError("tunnel not exists")
			_ = c.Send(p)
			return
		}

		conn := f.(net.Conn)
		err := conn.Close()
		if err != nil {
			p.SetError(err.Error())
		}

		_ = c.Send(p)
	})

	RegisterHandler(silk.TunnelData, func(c *Client, p *silk.Package) {
		p.Type = silk.TunnelDataAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := tunnels.Load(id)
		if !ok {
			p.SetError("tunnel not exists")
			_ = c.Send(p)
			return
		}

		conn := f.(net.Conn)
		_, err := conn.Write(p.Data[2:])
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.TunnelEnd, func(c *Client, p *silk.Package) {
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := tunnels.LoadAndDelete(id)
		if ok {
			conn := f.(net.Conn)
			_ = conn.Close()
		}
	})

}
