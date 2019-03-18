package servers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/soupstoregames/go-core/logging"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

var (
	ErrConnectionClosedAbnormally = errors.New("connection closed abnormally")
)

type WebsocketServer struct {
	Connections chan *WebsocketConnection

	server     *http.Server
	disposeLog func()
	addr       string
}

func NewWebsocketServer(addr string) *WebsocketServer {
	logger, _ := logging.WarnLogger()
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      nil,
		ErrorLog:     logger,
	}

	w := &WebsocketServer{
		Connections: make(chan *WebsocketConnection),
		server:      server,
		addr:        addr,
	}
	w.server.Handler = w.UpgradeToWebsocket()

	return w
}

func (ws *WebsocketServer) Start() {
	logging.Info("Websocket server listening on " + ws.addr)
	go func() {
		err := ws.server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				logging.Warn(err.Error())
				return
			}
			logging.Fatal(err.Error())
		}
	}()
}

func (ws *WebsocketServer) Stop(timeout time.Duration) {
	logging.Info("Stopping Websocket server")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := ws.server.Shutdown(ctx); err != nil {
		logging.Error(err.Error())
	}
	close(ws.Connections)
}

func (ws *WebsocketServer) UpgradeToWebsocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logging.Error(err.Error())
			return
		}

		conn := NewWebsocketConnection(c)
		ws.Connections <- conn
	}
}

type WebsocketConnection struct {
	*logging.ConnectionLogger

	Closed bool

	conn *websocket.Conn
	id   string
}

func NewWebsocketConnection(c *websocket.Conn) *WebsocketConnection {
	id := uuid.New().String()
	conn := &WebsocketConnection{
		ConnectionLogger: logging.BuildConnectionLogger(id),
		conn:             c,
		id:               id,
	}

	return conn
}

func (c *WebsocketConnection) ID() string {
	return c.id
}

func (c *WebsocketConnection) Close() error {
	if c.Closed {
		return nil
	}

	c.Info("Closing connection")
	c.Closed = true
	return c.conn.Close()
}

func (c *WebsocketConnection) WriteMessage(p []byte) (err error) {
	return c.conn.WriteMessage(websocket.BinaryMessage, p)
}

func (c *WebsocketConnection) ReadMessage() (msg []byte, err error) {
	_, msg, err = c.conn.ReadMessage()

	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
			err = ErrConnectionClosedAbnormally
		}
	}

	return
}
