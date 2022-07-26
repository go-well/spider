package spy

import (
	"github.com/go-well/spider/silk"
)

func init() {

	RegisterHandler(silk.TunnelCreate, func(c *Client, p *silk.Package) {
		p.Type++
		p.SetError("unsupported")

		_ = c.Send(p)
	})

}
