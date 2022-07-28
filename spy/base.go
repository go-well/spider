package spy

import (
	"github.com/go-well/spider/silk"
)

func (c *Client) Publish(topic string, payload []byte) error {
	ll := len(topic)
	buf := make([]byte, ll+1+len(payload))
	copy(buf, []byte(topic))
	buf[ll] = '\n'
	copy(buf[ll+1:], payload)
	_, err := c.Ask(&silk.Package{Type: silk.Publish, Data: buf})
	return err
}

func (c *Client) Subscribe(topic string) error {
	_, err := c.Ask(&silk.Package{Type: silk.Subscribe, Data: []byte(topic)})
	return err
}

func (c *Client) Unsubscribe(topic string) error {
	_, err := c.Ask(&silk.Package{Type: silk.Unsubscribe, Data: []byte(topic)})
	return err
}

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
