// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

const (
	NumClientBuffer = 100
	Debug           = true
)

var logger *log.Logger

// Package initialization function
func init() {
	if Debug {
		logger = log.New(os.Stdout, "", log.Lshortfile)
	}
}

type multiEchoServer struct {
	// TODO: implement this!
	isStarted   bool     //Signifies whether the server is started or not
	isClosed    bool     // Signifies whether the server is closed or not
	serverExit  bool     //Broadcast to all service channels to exit
	numConnChan chan int // Used to count the number of active servers
	count       int
	countLock   *sync.Mutex
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
		serverExit:  false,
		count:       0,
		countLock:   &sync.Mutex{},
		numConnChan: make(chan int, 1),
		exitChannel: make(chan bool, 1),
		serverChans: make(map[int]chan string),
	}
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!
	fuck("red", "We are starting the Server")
	var err error
	if mes.isClosed {
		return errors.New("Error: Server has been closed.")
	}
	if mes.listenConn, err = net.Listen("tcp",
		fmt.Sprintf(":%d", port)); err != nil {
		return err
	}
	mes.numConnChan <- 0
	fuck("red", "Gonna start dispatcher")
	go mes.dispatch()
	mes.isStarted = true
	fuck("red", "Start server successfully")
	return nil
}

// After the server started listenning on the ServerPort,
// this dispatcher will accept a connection and fork another
// go routine to handle this connection. When recieved the
// "exit" from server, it stops all of its connections
// immediately.
func (mes *multiEchoServer) dispatch() {

	mes.wg.Add(1)
	fuck("red", "We are starting the Dispatcher")

	for {
		select {
		case status := <-mes.exitChannel:
			fuck("cyan", fmt.Sprintf("Got exit message: %v", status))
			if status {
				// Tell other connection handlers to die
				mes.serverExit = true
				// terminate itself
				defer mes.wg.Done()
				return
			}
		default:
			fuck("cyan", "Waiting for connection...")
			conn, err := mes.listenConn.Accept()
			fuck("cyan", "Got an connection")
			mes.countLock.Lock()
			mes.count++
			mes.countLock.Unlock()
			if err != nil {
				fmt.Println("\tError on accept connection: %v",
					err)
			} else {
				fuck("red", "Gonna fork a handleConn")
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
	serverRead := make(chan string)
	clientExit := make(chan bool)
	go readFromClient(conn, serverRead, clientExit)
	count := <-mes.numConnChan
	count++
	mes.numConnChan <- count
	fuck("cyan", fmt.Sprintf("Handling %vth connection", count))
	defer mes.deCount()
	mes.wg.Add(1)
	mes.serverChans[count] = make(chan string, NumClientBuffer)
	defer delete(mes.serverChans, count)
	for !mes.serverExit {
		select {
		case message := <-mes.serverChans[count]:
			fuck("red", fmt.Sprintf("%v receieved:%s",
				count, message))
			if _, err := conn.Write([]byte(message)); err == io.EOF {
				fuck("cyan", fmt.Sprintf("%v gonna Exit", count))
				return
			}
		// For now, the server can only cache one message.
		case message := <-serverRead:
			select {
			case exit := <-clientExit:
				fuck("cyan",
					fmt.Sprintf("client %v send an exit(%v)",
						count, exit))
				if exit {
					fuck("red",
						fmt.Sprintf("Gonna lose %v connection.",
							count))
					return
				}
			default:
				for _, channel := range mes.serverChans {
					fuck("red", fmt.Sprintf("Client %v sends: %s",
						count, message))
					if len(channel) == cap(channel) {
						<-channel
					}
					channel <- message
				}

			}
		}
	}
	return
}

func readFromClient(conn net.Conn, serverRead chan string,
	clientExit chan bool) {
	var buf [2048]byte
	for {
		if n, err := conn.Read(buf[:]); err == io.EOF {
			fuck("cyan", "Detect one client is DOWN.")
			serverRead <- "die"
			clientExit <- true
		} else {
			serverRead <- string(buf[0:n])
		}
	}
}

func (mes *multiEchoServer) deCount() {
	mes.countLock.Lock()
	mes.count++
	mes.countLock.Unlock()
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	if mes.isStarted {
		fuck("cyan", "Server is closing")
		mes.exitChannel <- true
		mes.wg.Wait()

		// Flags "isStarted" and "isClosed" are only set by
		// server's close and start operations.
		mes.isStarted = false
		mes.isClosed = true
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

func fuck(color string, content string) {
	if logger != nil {
		logger.Output(2, colorText(color, content))
	}
}
func colorText(color string, content string) string {
	var code int
	if color == "cyan" {
		code = 36
	} else if color == "red" {
		code = 31
	} else {
		code = 37
	}
	return fmt.Sprintf("\x1b[%vm%s\x1b[0m", code, content)
}
