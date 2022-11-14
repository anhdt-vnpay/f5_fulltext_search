package runtime

type SearchProcessor interface {
	IndexData(data []byte) error
	SearchLite(query string) ([]byte, error)
}
