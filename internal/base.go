package internal

import (
	"encoding/json"
	"github.com/go-well/spider/silk"
	"github.com/super-l/machine-code/machine"
)

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

	RegisterHandler(silk.Connect, func(c *Client, p *silk.Package) {
		p.Type = silk.ConnectAck

		var info machine.MachineData
		err := json.Unmarshal(p.Data, &info)
		if err != nil {
			p.SetError(err.Error())
		}
		//TODO 查询数据库，找到设备
		//info.

		p.Data = nil
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
