package spy

import (
	"encoding/binary"
	"github.com/go-well/spider/silk"
	"net"
	"strings"
)

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
		Type: silk.TunnelDataEnd,
		Data: buf[:2],
	})

	c.tunnels.Delete(id)
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

		//缓存
		p.Data = make([]byte, 2)
		id := c.newTunnel(conn)
		binary.BigEndian.PutUint16(p.Data, id)
		go receiveTunnel(c, conn, id)

		_ = c.Send(p)
	})

	RegisterHandler(silk.TunnelClose, func(c *Client, p *silk.Package) {
		p.Type = silk.TunnelCloseAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := c.tunnels.LoadAndDelete(id)
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
		f, ok := c.tunnels.Load(id)
		if !ok {
			p.SetError("tunnel not exists")
			_ = c.Send(p)
			return
		}

		conn := f.(net.Conn)
		_, err := conn.Write(p.Data[2:])
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
		}
		//正常不回复
	})

	RegisterHandler(silk.TunnelDataEnd, func(c *Client, p *silk.Package) {
		p.Type = silk.TunnelDataEndAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := c.tunnels.LoadAndDelete(id)
		if !ok {
			p.SetError("tunnel not exists")
			_ = c.Send(p)
			return
		}
		conn := f.(net.Conn)
		err := conn.Close()
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
		}
		//正常不回复
	})

}
