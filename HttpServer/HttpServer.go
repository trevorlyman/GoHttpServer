package HttpServer

import (
	"net"
	"fmt"
	"os"
	//"strconv"
	"./Request"
	"./Response"
)

type HttpServer struct {
	running bool
}

func (s HttpServer) Run(port string) {
	port = ":" + port
	fmt.Println("Starting Http Server on Port " + port)

	tcpData, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	socket, err := net.ListenTCP("tcp", tcpData)

	for {
		conn, err := socket.Accept()
		if err != nil {
			fmt.Println("Error: " + err.Error())
			continue
		}

		go s.handleConn(conn)

	}

}

func (s HttpServer) handleConn(conn net.Conn) {
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if (err != nil) {
		fmt.Println("Error: " + err.Error())
		conn.Close()
		return
	}
	_ = n


	request := Request.NewRequest(buf)

	response := Response.NewResponse(request)
	response.Send(conn)

	conn.Close()
}
