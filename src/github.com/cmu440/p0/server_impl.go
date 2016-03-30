// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"errors"
	"fmt"
	"net"
	"sync"
    "io"
)

const (
	NumClientBuffer = 100
)

type multiEchoServer struct {
	// TODO: implement this!
	isStarted   bool      //Signifies whether the server is started or not
	isClosed    bool      // Signifies whether the server is closed or not
	numConnChan chan int  // Used to count the number of active servers
	exitChannel chan bool // Used to end the dispatcher
	listenConn  net.Listener
	wg          sync.WaitGroup      // A global waiter
	serverChans map[int]chan string // List of all server channels
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	return &multiEchoServer{
		isStarted:   false,
		isClosed:    false,
		numConnChan: make(chan int),
		exitChannel: make(chan bool),
		serverChans: make(map[int]chan string),
	}
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!
	var err error
	if mes.isClosed {
		return errors.New("Error: Server has been closed.")
	}
	if mes.listenConn, err = net.Listen("tcp",
		fmt.Sprintf(":%d", port)); err != nil {
		return err
	}
	mes.numConnChan <- 0
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

	mes.wg.Add(1)

	for {
		select {
		case status := <-mes.exitChannel:
			if status {
				// Tell other connection handlers to die

				// terminate itself
				defer mes.wg.Done()
				return
			}
		default:
			conn, err := mes.listenConn.Accept()
			if err != nil {
				fmt.Println("\tError on accept connection: %v",
					err)
			} else {
				go mes.handleConn(conn)
			}
		}
	}
}

/**
 * Handles every connection independently.
 */
func (mes *multiEchoServer) handleConn(conn net.Conn) {
	// After doing all works, remeber to "done" the waiter.
	defer mes.wg.Done()
    defer conn.Close()
	defer mes.deCount()
    defer delete(mes.serverChans, count)
	mes.wg.Add(1)
    serverRead := make(chan string)
    clientExit := make(chan bool)
    go readFromClient(conn, serverRead, clientExit)
	count := <-mes.numConnChan
	count++
	mes.numConnChan <- count
	mes.serverChans[count+1] = make(chan string, NumClientBuffer)
	for !mes.serverExit{
		select {
		case message := <-mes.serverChans[count]:
            if _, err := conn.Write([]byte(message)); err == io.EOF{
                return
            }
        // For now, the server can only cache one message.
        case message := <-serverRead:
            select {
                case <- clientExit:
                    return
                default:
                    for _, channel := range mes.serverChans{
                if len(channel) == cap(channel){
                    <- channel
                }
                channel <- message
            }

           }
       }
	}
    return 
}

func readFromClient(conn net.Conn, serverRead chan string, 
    clientExit chan bool){
    var buf [2048]byte
    for{
        if n, err := conn.Read(buf[:]); err == io.EOF{
            clientExit <- true
            serverRead <- "die"
        }else{
            serverRead <- string(buf[0,n])
        }
    }
}

func (mes *multiEchoServer) deCount() {
	count := <-mes.numConnChan
	count--
	mes.numConnChan <- count
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	if mes.isStarted {
		mes.exitChannel <- true
		mes.wg.Wait()

		// Flags "isStarted" and "isClosed" are only set by
		// server's close and start operations.
		mes.isStarted = false
		mes.isClosed = true
		close(mes.exitChannel)
		close(mes.numConnChan)
		fmt.Println("\tNow the server is down.")
	} else if mes.isClosed {
		fmt.Println("\tThe server has been closed.")
	} else {
		fmt.Println("\tThe server is not started.")
	}
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	count := <-mes.numConnChan
	mes.numConnChan <- count
	return count
}

// TODO: add additional methods/functions below!
