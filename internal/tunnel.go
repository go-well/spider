package internal

import (
	"encoding/binary"
	"github.com/go-well/spider/silk"
	"io"
)

//TODO 封装成Tunnel类
func receiveTunnel(c *Client, conn io.ReadWriteCloser, id uint16) {
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

func (c *Client) TunnelBind(conn io.ReadWriteCloser) error {
	p, err := c.Ask(&silk.Package{Type: silk.TunnelCreate})
	if err != nil {
		return err
	}

	id := binary.BigEndian.Uint16(p.Data)
	c.tunnels.Store(id, conn)

	go receiveTunnel(c, conn, id)

	return nil
}

func init() {

	RegisterHandler(silk.TunnelClose, func(c *Client, p *silk.Package) {
		p.Type = silk.TunnelCloseAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := c.tunnels.LoadAndDelete(id)
		if !ok {
			p.SetError("tunnel not exists")
			_ = c.Send(p)
			return
		}

		conn := f.(io.ReadWriteCloser)
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

		conn := f.(io.ReadWriteCloser)
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
		conn := f.(io.ReadWriteCloser)
		err := conn.Close()
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
		}
		//正常不回复
	})

}
