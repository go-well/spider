package internal

import (
	"net"
)

func Serve() error {

	listener, err := net.Listen("tcp", ":1206")
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		//创建连接
		_ = &Client{conn: conn}
	}
}
