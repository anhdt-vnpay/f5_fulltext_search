package runtime

type SearchProcessor interface {
	Subscribe() error
	IndexData(data []byte) error
	SearchLite(query string) ([]byte, error)
}
