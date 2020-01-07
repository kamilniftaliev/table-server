package types

import (
	// "fmt"

	"io"
	"strconv"
	"time"

	// "log"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const HUMAN_TIME_FORMAT = "02.01.2006 15:04"

var DateTime *primitive.DateTime

// Lets redefine the base ID type to use an id from an external library
func MarshalDateTime(dateTime primitive.DateTime) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(dateTime.Time().Format(HUMAN_TIME_FORMAT)))
	})
}

// And the same for the unmarshaler
func UnmarshalDateTime(v interface{}) (primitive.DateTime, error) {
	str, _ := v.(string)

	time, err := time.Parse(HUMAN_TIME_FORMAT, str)
	dateTime := primitive.NewDateTimeFromTime(time)

	return dateTime, err
}
