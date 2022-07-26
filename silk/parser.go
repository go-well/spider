package silk

import (
	"github.com/go-well/spider/lib"
)

type OnPackage func(p *Package)

type Parser struct {
	remain []byte
}

func (p *Parser) Parse(data []byte) ([]*Package, error) {
	//TODO 处理状态
	packs := make([]*Package, 0)

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

		//回调中放入队列处理
		packs = append(packs, pack)
		//p.onPackage(pack)
	}

	//处理剩余长度
	remain := len(data)
	if remain > 0 {
		p.remain = data
	} else {
		p.remain = nil
	}

	return packs, nil
}
