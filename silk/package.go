package silk

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
)

const MAGIC = "Spider"

/**
数据包
Spider - 6byte
Id - 2byte
Type - 1byte
Length - 1~3byte
Data - 0~65535byte
*/

type Package struct {
	Id   uint16
	Type Type
	Fail bool
	Data []byte
}

func uVarIntSize(x uint64) {
	binary.PutUvarint(nil, x)
}

func (p *Package) Encode() []byte {
	buf := make([]byte, 12+len(p.Data))
	copy(buf, MAGIC)
	binary.BigEndian.PutUint16(buf[6:], p.Id)
	buf[8] = p.Type
	if p.Fail {
		buf[8] &= 0x80
	}
	n := binary.PutUvarint(buf[9:], uint64(len(p.Data)))
	copy(buf[9+n:], p.Data)
	return buf[:9+n+len(p.Data)]
}

func (p *Package) Decode(buf []byte) (uint64, error) {

	//寻找魔术头
	var cursor uint64 = 0
	for bytes.Compare([]byte(MAGIC), buf) != 0 {
		if len(buf) < 10 {
			return cursor, errors.New("数据包长度不能小于10")
		}
		buf = buf[1:]
		cursor++
	}

	//数据包长度
	length, n := binary.Uvarint(buf[9:])
	size := 9 + uint64(n) + length
	if size < uint64(len(buf)) {
		return cursor, errors.New("数据包长度不够")
	}

	//解析
	p.Id = binary.BigEndian.Uint16(buf[6:])
	p.Type = buf[8]
	//p.Data = lib.Dup(buf[12:size])
	p.Data = buf[9+n : size]

	//返回解析的长度
	return cursor + size, nil
}
