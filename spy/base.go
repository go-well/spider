package spy

import (
	"github.com/go-well/spider/silk"
)

func init() {

	RegisterHandler(silk.Close, func(c *Client, p *silk.Package) {
		_ = c.Close()
	})

	RegisterHandler(silk.Ping, func(c *Client, p *silk.Package) {
		p.Type = silk.Pong
		_ = c.Send(p)
	})

	RegisterHandler(silk.Spawn, func(c *Client, p *silk.Package) {
		p.Type = silk.SpawnAck
		err := c.Spawn()
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

}
