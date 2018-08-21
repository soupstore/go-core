package connections

type Connection interface {
	ID() string
	Close() error
	WriteMessage([]byte) error
	ReadMessage() (msg []byte, err error)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}
