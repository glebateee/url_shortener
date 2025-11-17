package slogdiscard

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{}

func (d *DiscardHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (d *DiscardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (d *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return nil
}

func (d *DiscardHandler) WithGroup(_ string) slog.Handler {
	return nil
}

func NewDiscardhandler() *DiscardHandler {
	return &DiscardHandler{}
}
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardhandler())
}
