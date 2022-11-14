package runtime

type MessageProcessor interface {
	Save(data []byte) error
}