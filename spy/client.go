package spy

import (
	"encoding/json"
	"github.com/go-well/spider/silk"
	"github.com/pkg/errors"
	"io"
	"net"
	"os"
	"os/exec"
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

func ReplaceHandler(tp silk.Type, handler Handler) {
	handlers[tp] = Handlers{handler}
}

type task struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

type Options struct {
	Addr      string `yaml:"addr" json:"addr"`
	Once      bool   `yaml:"once,omitempty" json:"once,omitempty"`
	Heartbeat int    `yaml:"heartbeat,omitempty" json:"heartbeat,omitempty"`
}

type Client struct {
	options Options

	conn   net.Conn
	parser silk.Parser

	//处理队列
	packages chan *silk.Package

	//缓存文件
	files     sync.Map
	fileIndex uint16

	tunnels     sync.Map
	tunnelIndex uint16

	tasks     sync.Map
	taskIndex uint16

	heartbeatTicker *time.Timer

	requests     sync.Map
	requestIndex uint16
}

func (c *Client) newFile(file *os.File) uint16 {
	c.fileIndex++
	c.files.Store(c.fileIndex, file)
	return c.fileIndex
}

func (c *Client) newTunnel(conn net.Conn) uint16 {
	c.tunnelIndex++
	c.tunnels.Store(c.tunnelIndex, conn)
	return c.tunnelIndex
}

func (c *Client) newTask(cmd *exec.Cmd, stdin io.WriteCloser) uint16 {
	c.taskIndex++
	c.tasks.Store(c.taskIndex, &task{
		cmd:   cmd,
		stdin: stdin,
	})
	return c.taskIndex
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

func (c *Client) reconnect() {
	time.AfterFunc(time.Minute, func() {
		_ = c.Open()
	})
}

func (c *Client) connect() {
	data, _ := json.Marshal(&regPack)
	_ = c.Send(&silk.Package{
		Type: silk.Connect,
		Data: data,
	})
}

func (c *Client) heartbeat() {
	if c.options.Heartbeat == 0 {
		return
	}

	c.heartbeatTicker = time.AfterFunc(time.Second*time.Duration(c.options.Heartbeat), func() {
		//log.Println("heartbeat ticker active")
		_ = c.Send(&silk.Package{Type: silk.Heartbeat})
		c.heartbeatTicker.Reset(time.Second * time.Duration(c.options.Heartbeat))
	})
}

func (c *Client) run() {
	c.heartbeat()
	c.connect()

	for {
		var buf = make([]byte, 1024)

		n, err := c.conn.Read(buf)
		if err != nil {
			break
		}

		//重置心跳
		c.heartbeatTicker.Reset(time.Second * time.Duration(c.options.Heartbeat))

		packs, err := c.parser.Parse(buf[:n])
		for _, p := range packs {
			c.handle(p)
		}
		if err != nil {
			//print
		}
	}

	//重连
	if !c.options.Once {
		c.reconnect()
	}

	//关闭心跳
	if c.heartbeatTicker != nil {
		c.heartbeatTicker.Stop()
	}
}

func (c *Client) Send(p *silk.Package) error {
	//重置心跳
	c.heartbeatTicker.Reset(time.Second * time.Duration(c.options.Heartbeat))

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
		return nil, errors.New("Timeout")
	}
}

func (c *Client) Close() error {
	//TODO 判断是否已经关闭
	//TODO 关闭文件，关闭通道，关闭task
	return c.conn.Close()
}

func (c *Client) Spawn() error {
	cc := &Client{
		options:  c.options,
		packages: make(chan *silk.Package, 64),
	}
	return cc.Open()
}

func (c *Client) Open() error {
	var err error
	c.conn, err = net.Dial("tcp", c.options.Addr)
	if err != nil {
		//初次未连接成功，也要重连
		if !c.options.Once {
			c.reconnect()
		}

		return err
	}

	go c.run()
	return nil
}

func Connect(options Options) (*Client, error) {
	cli := &Client{
		options:  options,
		packages: make(chan *silk.Package, 64), //TODO 需要从配置调整
	}
	//默认心跳 60s
	if cli.options.Heartbeat == 0 {
		cli.options.Heartbeat = 60
	}
	return cli, cli.Open()
}
