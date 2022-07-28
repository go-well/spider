package spy

import (
	"encoding/binary"
	"github.com/go-well/spider/silk"
	"os/exec"
	"strings"
)

func init() {

	RegisterHandler(silk.TaskCreate, func(c *Client, p *silk.Package) {
		p.Type = silk.TaskCreateAck

		names := strings.Split(string(p.Data), " ")
		cmd := exec.Command(names[0], names[1:]...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}
		stdin, err := cmd.StdinPipe()
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}

		err = cmd.Run()
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}

		//缓存
		id := c.newTask(cmd, stdin)
		buf := make([]byte, 512)
		binary.BigEndian.PutUint16(buf, id)

		p.Data = buf[:2]
		_ = c.Send(p)

		go func() {
			for {
				n, e := stdout.Read(buf[2:])
				if e != nil {
					break
				}
				_ = c.Send(&silk.Package{
					Type: silk.TaskData,
					Data: buf[:2+n],
				})
			}
			_ = c.Send(&silk.Package{
				Type: silk.TaskDataEnd,
				Data: buf[:2],
			})
		}()
	})

	RegisterHandler(silk.TaskData, func(c *Client, p *silk.Package) {
		p.Type = silk.TaskDataAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := c.files.Load(id)
		if !ok {
			p.SetError("task not exists")
			_ = c.Send(p)
			return
		}

		t := f.(*task)
		_, err := t.stdin.Write(p.Data[2:])
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.TaskRun, func(c *Client, p *silk.Package) {
		p.Type = silk.TaskRunAck
		names := strings.Split(string(p.Data), " ")
		cmd := exec.Command(names[0], names[1:]...)
		err := cmd.Run()
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}
		p.Data, err = cmd.CombinedOutput()
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.TaskStart, func(c *Client, p *silk.Package) {
		p.Type = silk.TaskStartAck
		names := strings.Split(string(p.Data), " ")
		cmd := exec.Command(names[0], names[1:]...)
		err := cmd.Start()
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.TaskKill, func(c *Client, p *silk.Package) {
		p.Type = silk.TaskKillAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := c.files.LoadAndDelete(id)
		if !ok {
			p.SetError("task not exists")
			_ = c.Send(p)
			return
		}

		t := f.(*task)
		err := t.cmd.Process.Kill()
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})
}
