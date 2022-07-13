package models

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// Lets redefine the base ID type to use an id from gorm's default uint type
func MarshalID(id uint) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", id)))
	})
}

// And the same for the unmarshaler
func UnmarshalID(v interface{}) (uint, error) {
	str, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("ids must be strings")
	}
	i, err := strconv.Atoi(str[1 : len(str)-1])
	return uint(i), err
}
