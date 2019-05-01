package servers

import (
	"github.com/soupstoregames/go-core/logging"
	"io"
	"net"
	"strings"
)

type TCPHandler func([]byte, TCPWriter)

type TCPWriter interface {
	Write(b []byte) (n int, err error)
}

type TCPServer struct {
	listener net.Listener
	handler  TCPHandler
	addr     string
}

func NewTCPServer(addr string, handler TCPHandler) *TCPServer {
	return &TCPServer{
		addr:    addr,
		handler: handler,
	}
}

func (t *TCPServer) Start() error {
	var err error

	if t.listener, err = net.Listen("tcp", t.addr); err != nil {
		return err
	}

	go func() {
		logging.Info("TCP Server listening on " + t.addr)

		for {
			// Listen for an incoming connection.
			conn, err := t.listener.Accept()
			if err != nil {
				// net.errClosing is not exported so this
				if strings.Contains(err.Error(), "use of closed network connection") {
					return
				}
				logging.Error(err.Error())
			}

			logging.Debug("Client connected: " + conn.RemoteAddr().String())

			// Handle connections in a new goroutine.
			go func() {
				buffer := make([]byte, 1024)

				for {
					n, err := conn.Read(buffer)
					if err != nil {
						if err == io.EOF {
							logging.Debug("Client disconnected: " + conn.RemoteAddr().String())
							return
						}
						logging.Error(err.Error())
						return
					}

					t.handler(buffer[:n], conn)
				}
			}()
		}
	}()

	return nil
}

func (t *TCPServer) Stop() {
	logging.Info("Stopping TCP Server")
	t.listener.Close()
}
