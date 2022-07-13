package common

import (
	"context"
)

type EventHandler interface {
	Handle(ctx context.Context, eventData []byte)
}
