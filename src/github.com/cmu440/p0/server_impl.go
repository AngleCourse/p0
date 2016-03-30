// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"errors"
	"fmt"
	"net"
)

const (
	ServerPort = 9999
)

type multiEchoServer struct {
	// TODO: implement this!
	isStarted   bool      //Signifies whether the server is started or not
	isClosed    bool      // Signifies whether the server is closed or not
	numConn     uint      // Used to count the number of active clients
	exitChannel chan bool // Used to end the dispatcher
	listenConn  net.Listener
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	return &multiEchoServer{
		isStarted:   false,
		isClosed:    false,
		numConn:     0,
		exitChannel: make(chan bool),
	}
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!
	var err error
	if mes.isClosed {
		return errors.New("Error: Server has been closed.")
	}
	if mes.listenConn, err = net.Listen("tcp",
		fmt.Sprintf(":%d", ServerPort)); err != nil {
		return err
	}
	go mes.dispatch()
	mes.isStarted = true
	return nil
}

// After the server started listenning on the ServerPort,
// this dispatcher will accept a connection and fork another
// go routine to handle this connection. When recieved the
// "exit" from server, it stops all of its connections
// immediately.
func (mes *multiEchoServer) dispatch() {

}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	return -1
}

// TODO: add additional methods/functions below!
