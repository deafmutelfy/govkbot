package main

import (
	"context"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/events"
)

type eventManager interface {
	MessageNew(f func(context.Context, events.MessageNewObject))
}

type eventManagerCallback interface {
	eventManager

	HandleFunc(w http.ResponseWriter, r *http.Request)
}

type eventManagerLongpoll interface {
	eventManager

	Run() error
}
