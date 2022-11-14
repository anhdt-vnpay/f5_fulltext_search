package runtime

type SearchProcessor interface {
	SearchLite(query string) ([]byte, error)
}