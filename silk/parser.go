package silk

import (
	"github.com/zgwit/spider/lib"
	"net"
)

type Handler func(p *Package)

type Parser struct {
	conn    net.Conn
	remain  []byte
	handler Handler
}

func NewParser(conn net.Conn, handler Handler) *Parser {
	return &Parser{conn: conn, handler: handler}
}

func (p *Parser) Parse() {
	//TODO 处理状态

	//var remain = 0
	var data []byte
	var buf = make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}

		//取有效数据
		data = buf[:n]

		//拼接历史数据
		if p.remain != nil {
			//b := make([]byte, len(p.remain) + len(data))
			//copy(b, p.remain)
			//copy(b[:len(p.remain)], data)
			//data = b
			copy(p.remain[:len(p.remain)], data)
			data = p.remain
		} else {
			data = lib.Dup(data)
		}

		//解析数据，有剩余长度，且大于12
		for len(data) >= 12 {
			pack := &Package{}
			i, e := pack.Decode(data)
			data = data[:i]
			//有异常退出，可能是长度不够
			if e != nil {
				break
			}
			p.handler(pack)
		}

		//处理剩余长度
		remain := len(data)
		if remain > 0 {
			p.remain = data
		} else {
			p.remain = nil
		}
	}
}
