package spy

import (
	"github.com/go-well/spider/silk"
	"io"
	"net"
	"os"
	"os/exec"
	"sync"
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

type Client struct {
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

func (c *Client) Close() error {
	//TODO 判断是否已经关闭
	//TODO 关闭文件，关闭通道，关闭task
	return c.conn.Close()
}

func (c *Client) Spawn() error {

	return nil
}

func newClient(conn net.Conn) *Client {
	cli := &Client{
		conn:     conn,
		packages: make(chan *silk.Package, 64), //TODO 需要从配置调整
	}
	go cli.run()
	return cli
}

func Open() (*Client, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:1206")
	if err != nil {
		return nil, err
	}

	return newClient(conn), nil
}
