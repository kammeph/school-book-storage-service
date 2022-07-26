package domain

import (
	"errors"
	"fmt"
)

var ErrReasonNotSpecified = errors.New("no reason specified")

func ErrUnknownEvent(event Event) error {
	return fmt.Errorf("unhandled event %T", event)
}
