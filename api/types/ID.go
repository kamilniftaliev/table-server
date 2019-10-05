package types

import (
	// "fmt"

	"io"
	"strconv"

	// "log"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ID *primitive.ObjectID

// Lets redefine the base ID type to use an id from an external library
func MarshalID(id primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(id.Hex()))
	})
}

// And the same for the unmarshaler
func UnmarshalID(v interface{}) (primitive.ObjectID, error) {
	str, _ := v.(string)

	i, err := primitive.ObjectIDFromHex(str)

	return i, err
}
