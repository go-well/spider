package spy

import (
	"encoding/binary"
	"encoding/json"
	"github.com/go-well/spider/silk"
	"os"
	"os/exec"
	"strings"
)

func init() {

	RegisterHandler(silk.SystemShell, func(c *Client, p *silk.Package) {
		p.Type = silk.SystemShellAck
		p.SetError("unsupported")
		_ = c.Send(p)
	})

	RegisterHandler(silk.SystemExecute, func(c *Client, p *silk.Package) {
		p.Type = silk.SystemExecuteAck
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

	RegisterHandler(silk.SystemStart, func(c *Client, p *silk.Package) {
		p.Type = silk.SystemStartAck
		names := strings.Split(string(p.Data), " ")
		cmd := exec.Command(names[0], names[1:]...)
		err := cmd.Start()
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.SystemKill, func(c *Client, p *silk.Package) {
		p.Type = silk.SystemKillAck
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

	RegisterHandler(silk.SystemEnvironment, func(c *Client, p *silk.Package) {
		p.Type = silk.SystemEnvironmentAck
		env := os.Environ()
		p.Data, _ = json.Marshal(env)
		_ = c.Send(p)
	})

}
