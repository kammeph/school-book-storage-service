package application

import (
	"context"
)

type EventHandler interface {
	Handle(ctx context.Context, eventData []byte)
}

type EventSubscriber interface {
	Subscribe(exchange string, handler EventHandler) error
}
