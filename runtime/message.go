package runtime

type MessageProcessor interface {
	Save(data []byte) error
	GetMsgChannel() chan []byte
}
