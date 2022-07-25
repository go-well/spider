package internal

import (
	"net"
)

type Client struct {
	conn   net.Conn
	remain []byte
}
