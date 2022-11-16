package runtime

import (
	"fmt"
	"strings"
)

func createMessage(tipe string, data interface{}) ([]byte, error) {
	tipe = strings.Trim(tipe, " ")
	tipe = strings.ToLower(tipe)
	m := fmt.Sprintf("%s_%v", tipe, data)
	fmt.Println("message: ", m)
	return []byte(m), nil
}
