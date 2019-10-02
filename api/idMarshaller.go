package api

import (
	// "fmt"

	"io"
	"log"
	"strconv"

	// "log"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Lets redefine the base ID type to use an id from an external library
func MarshalID(id primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(id.Hex()))
	})
}

// And the same for the unmarshaler
func UnmarshalID(v interface{}) (primitive.ObjectID, error) {
	// str, ok := v.(string)

	log.Println("AAAA:", v)

	// if !ok {
	// 	return primitive.ObjectIDFromHex("SALAM"), ok
	// }

	i, err := primitive.ObjectIDFromHex("5d94b4b30e9c3b268c59448f")

	if err != nil {
		return primitive.ObjectIDFromHex("5d94b4b30e9c3b268c59448f")
		// return log.Println("err::", err), err
	}

	return i, err
}
