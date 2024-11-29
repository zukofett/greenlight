package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	jsonValue = strconv.Quote(jsonValue)

	return []byte(jsonValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSON, err := strconv.Unquote(string(jsonValue))
	if nil != err {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSON, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	n, err := strconv.ParseInt(parts[0], 10, 32)
	if nil != err {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(n)

	return nil
}
