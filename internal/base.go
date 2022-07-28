package internal

import "github.com/go-well/spider/silk"

func (c *Client) Publish(topic string, payload []byte) error {
	ll := len(topic)
	buf := make([]byte, ll+1+len(payload))
	copy(buf, []byte(topic))
	buf[ll] = '\n'
	copy(buf[ll+1:], payload)
	_, err := c.Ask(&silk.Package{Type: silk.Message, Data: buf})
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

	RegisterHandler(silk.Publish, func(c *Client, p *silk.Package) {
		//TODO server publish

		p.Type = silk.PublishAck
		p.Data = nil
		_ = c.Send(p)
	})

	RegisterHandler(silk.Subscribe, func(c *Client, p *silk.Package) {
		c.topics.Store(string(p.Data), true)
		p.Type = silk.SubscribeAck
		p.Data = nil
		_ = c.Send(p)
	})

	RegisterHandler(silk.Unsubscribe, func(c *Client, p *silk.Package) {
		c.topics.Delete(string(p.Data))
		p.Type = silk.UnsubscribeAck
		p.Data = nil
		_ = c.Send(p)
	})

}
