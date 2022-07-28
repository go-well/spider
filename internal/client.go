package internal

import (
	"errors"
	"github.com/go-well/spider/silk"
	"net"
	"sync"
	"time"
)

type Handler func(c *Client, p *silk.Package)
type Handlers []Handler

var handlers = map[silk.Type]Handlers{}

func RegisterHandler(tp silk.Type, handler Handler) {
	hs, ok := handlers[tp]
	if !ok {
		hs = Handlers{handler}
		handlers[tp] = hs
	} else {
		handlers[tp] = append(hs, handler)
	}
}

type Client struct {
	conn   net.Conn
	parser silk.Parser

	files   sync.Map
	tunnels sync.Map
	tasks   sync.Map

	requests     sync.Map
	requestIndex uint16
}

func (c *Client) handle(p *silk.Package) {
	//先处理会话
	if p.Id > 0 {
		if ch, ok := c.requests.LoadAndDelete(p.Id); ok {
			cc := ch.(chan *silk.Package)
			cc <- p
		}
	}

	if hs, ok := handlers[p.Type]; ok {
		for _, h := range hs {
			h(c, p)
		}
	}
}

func (c *Client) run() {

	for {
		var buf = make([]byte, 1024)

		n, err := c.conn.Read(buf)
		if err != nil {
			break
		}

		packs, err := c.parser.Parse(buf[:n])
		for _, p := range packs {
			c.handle(p)
		}
		if err != nil {
			//print
		}
	}

}

func (c *Client) Send(p *silk.Package) error {
	buf := p.Encode()
	_, err := c.conn.Write(buf)
	return err
}

func (c *Client) Ask(p *silk.Package) (*silk.Package, error) {
	//分配ID
	c.requestIndex++
	if c.requestIndex == 0 {
		c.requestIndex++
	}
	id := c.requestIndex
	p.Id = id
	err := c.Send(p)
	if err != nil {
		return nil, err
	}

	//等待结果和超时
	ch := make(chan *silk.Package)
	c.requests.Store(id, ch)

	select {
	case p := <-ch:
		if p.Fail {
			return nil, errors.New(string(p.Data))
		}
		return p, nil
	case <-time.After(time.Minute):
		return nil, errors.New("timeout")
	}
}

func (c *Client) Close() error {
	//TODO 判断是否已经关闭
	//TODO 关闭文件，关闭通道，关闭task
	return c.conn.Close()
}
