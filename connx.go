package redix

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type Connx struct {
	readWho string // "src" , "dest

	srcConn  redis.Conn
	destConn redis.Conn
}

func NewConnx(srcConn redis.Conn, destConn redis.Conn, readWho string) *Connx {
	return &Connx{
		readWho:  readWho,
		srcConn:  srcConn,
		destConn: destConn,
	}
}

func (c *Connx) Close() error {
	e1 := c.srcConn.Close()
	e2 := c.destConn.Close()

	if e1 == nil && e2 == nil {
		return nil
	}
	return fmt.Errorf("close err, dest err: %v src err: %v", e1, e2)
}

func (c *Connx) Err() error {
	return fmt.Errorf("dest err: %v, src err: %v", c.destConn.Err(), c.srcConn.Err())
}

func (c *Connx) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if isRead(commandName) {
		switch c.readWho {
		case "src":
			return c.srcConn.Do(commandName, args...)
		case "dest":
			return c.destConn.Do(commandName, args...)
		default:
			return c.srcConn.Do(commandName, args...)
		}
	}
	rs1, e := c.srcConn.Do(commandName, args...)

	rs2, e := c.destConn.Do(commandName, args...)

	switch c.readWho {
	case "src":
		return rs1, e
	case "dest":
		return rs2, e
	default:
		return rs1, e
	}
}

func (c *Connx) Send(commandName string, args ...interface{}) error {
	e1 := c.destConn.Send(commandName, args...)
	e2 := c.srcConn.Send(commandName, args...)

	if e1 == nil && e2 == nil {
		return nil
	}

	return fmt.Errorf("send err: dest err:%v src err: %v", e1, e2)
}

func (c *Connx) Flush() error {
	e1 := c.destConn.Flush()
	e2 := c.srcConn.Flush()

	if e1 == nil && e2 == nil {
		return nil
	}
	return fmt.Errorf("send err: dest err:%v src err: %v", e1, e2)
}

func (c *Connx) Receive() (reply interface{}, err error) {
	switch c.readWho {
	case "src":
		return c.srcConn.Receive()
	case "dest":
		return c.destConn.Receive()
	default:
		return c.srcConn.Receive()
	}
}
