package internal

import (
	"encoding/binary"
	"github.com/go-well/spider/silk"
	"io"
)

func (c *Client) TaskCreate(command string, conn io.ReadWriteCloser) (uint16, error) {
	p, err := c.Ask(&silk.Package{Type: silk.TaskCreate, Data: []byte(command)})
	if err != nil {
		return 0, err
	}

	id := binary.BigEndian.Uint16(p.Data)
	c.tasks.Store(id, conn)

	go func() {
		buf := make([]byte, 512)
		for {
			n, e := conn.Read(buf[2:])
			if e != nil {
				break
			}
			_ = c.Send(&silk.Package{Type: silk.TaskData, Data: buf[:2+n]})
		}
		_ = c.Send(&silk.Package{Type: silk.TaskDataEnd, Data: buf[:2]})
	}()

	return id, nil
}

func (c *Client) TaskRun(command string) (string, error) {
	p, err := c.Ask(&silk.Package{Type: silk.TaskRun, Data: []byte(command)})
	if err != nil {
		return "", err
	}
	return string(p.Data), nil
}

func (c *Client) TaskStart(command string) error {
	_, err := c.Ask(&silk.Package{Type: silk.TaskStart, Data: []byte(command)})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) TaskKill(id uint16) error {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, id)
	_, err := c.Ask(&silk.Package{Type: silk.TaskRun, Data: buf})
	if err != nil {
		return err
	}
	return nil
}

func init() {

	RegisterHandler(silk.TaskData, func(c *Client, p *silk.Package) {
		p.Type = silk.TaskDataAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := c.files.Load(id)
		if !ok {
			p.SetError("task not exists")
			_ = c.Send(p)
			return
		}

		t := f.(io.ReadWriteCloser)
		_, err := t.Write(p.Data[2:])
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})
}
