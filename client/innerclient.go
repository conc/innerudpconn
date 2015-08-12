package innerclient

import (
	"errors"
	"log"
	"net"
	"sync/atomic"
	"time"
)

type InnerClient struct {
	PoolSize      int
	connPool      chan *net.Conn
	ServerAddr    string
	ErrRetryTimes int
}

var reqNum uint64

func (c *InnerClient) Init() {
	log.Println("Loading conn...............")
	c.connPool = make(chan *net.Conn, c.PoolSize)
	for i := 0; i < c.PoolSize; i++ {
		connTemp, err := c.createConn()
		if err != nil {
			log.Println("create conn error:", err)
			panic(err)
		}
		c.connPool <- connTemp
	}

	log.Println("Loading conn over..........")
	return
}

func (c *InnerClient) Request(reqData []byte, reqType byte) (res []byte, err error) {

	req := connStu{RequestType: reqType, Data: reqData}
	req.DataLen = uint64(len(reqData))
	req.RequestId = atomic.AddUint64(&reqNum, 1)

	var ret *connStu
	for i := 0; i < c.ErrRetryTimes+1; i++ {
		ret, err = c.sendReceive(&req)
		if err == nil {
			res = ret.Data
			return
		} else {
			continue
		}
	}
	return
}

func (c *InnerClient) sendReceive(req *connStu) (res *connStu, err error) {
	var client *net.Conn
	select {
	case client = <-c.connPool:
	case <-time.After(2 * time.Second):
		err = errors.New("Get nothing from connPool")
		return
	}

	err = c.send(client, req)
	if err != nil {
		log.Println(err)
		return
	}

	res, err = c.receive(client)
	if err != nil {
		log.Println(err)
		return
	}

	if res.RequestId != req.RequestId {
		err = errors.New("receive data error(requsetid error)!")
		return
	}
	return
}

func (c *InnerClient) send(client *net.Conn, req *connStu) (err error) {
	reqByte := connStuToBytes(req)
	(*client).SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err = (*client).Write(reqByte)
	if err != nil {
		go c.dealErrorConn(client)
		return
	}
	return
}

func (c *InnerClient) receive(client *net.Conn) (ret *connStu, err error) {
	var dataBuf [1460]byte
	var readLen int
	(*client).SetReadDeadline(time.Now().Add(5 * time.Second))
	readLen, err = (*client).Read(dataBuf[:])
	if err != nil {
		go c.dealErrorConn(client)
		return
	}
	ret, err = bytesToConnStu(dataBuf[0:readLen])
	c.connPool <- client
	return
}

func (c *InnerClient) createConn() (*net.Conn, error) {
	conn, err := net.Dial("udp", c.ServerAddr)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func (c *InnerClient) dealErrorConn(client *net.Conn) {
	(*client).Close()
	var err error
	for {
		if client, err = c.createConn(); err != nil {
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	c.connPool <- client
	return
}
