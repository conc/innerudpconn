package innerserver

import (
	"log"
	"net"
)

type TransactBussiness func(byte, []byte) []byte

type InnerServer struct {
	ListenAddr      net.UDPAddr
	TransactProcess TransactBussiness
}

func (c *InnerServer) Start() {
	conn, err := net.ListenUDP("udp4", &(c.ListenAddr))
	if err != nil {
		log.Println("err ListenUDP:", err)
		panic(err)
	}
	c.dealConn(conn)
}

func (c *InnerServer) dealConn(conn *net.UDPConn) {
	defer conn.Close()
	var from *net.UDPAddr
	var data [1460]byte
	var err error
	var datalen int

	for {
		datalen, from, err = conn.ReadFromUDP(data[:])
		if err != nil {
			log.Println("Error conn.Read:", err)
			continue
		}

		res := c.dealReceiveData(data[0:datalen])
		resByte := connStuToBytes(res)
		_, err = conn.WriteToUDP(resByte, from)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func (c *InnerServer) dealReceiveData(data []byte) *connStu {
	var res connStu
	req, err := bytesToConnStu(data)
	if err != nil {
		return &res
	}

	res.Data = c.TransactProcess(req.RequestType, req.Data)
	res.RequestType = req.RequestType
	res.RequestId = req.RequestId
	res.DataLen = uint64(len(res.Data))
	return &res
}
