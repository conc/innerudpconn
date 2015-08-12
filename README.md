# innerudpconn

#How to use:

go get github.com/conc/innerudpconn/server   
go get github.com/conc/innerudpconn/client   


#Examples:

```golang
client.go:
package main

import (
	"github.com/conc/innerudpconn/client"
	"log"
)

func main() {
	var clientReq innerclient.InnerClient
	clientReq.PoolSize = 10
	clientReq.ServerAddr = "127.0.0.1:4333"
	clientReq.ErrRetryTimes = 2
	clientReq.Init()

	for i := 0; i < 1000000; i++ {
		ret, err := clientReq.Request([]byte("hhhhhhhhhh"), 99)
		log.Println(i, string(ret))
		if err != nil {
			log.Println(err)
		}
	}

	return
}

server.go:
package main

import (
	"github.com/conc/innerudpconn/server"
	"net"
    "log"
)

func main() {
	var server innerserver.InnerServer
	server.ListenAddr = net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 4333,
	}
	server.TransactProcess = dealTestReq

	server.Start()
	return
}

func dealTestReq(reqType byte, reqData []byte) []byte {

	switch reqType {
	case 0:
	case 1:
	case 2:
	case 99:
		log.Println(string(reqData))
		return []byte("----")
	default:
	}

	return []byte("...")
}
```


