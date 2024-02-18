package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	yello = "\033[33m"
	blue  = "\033[34m"
	white = "\033[37m"
	reset = "\033[0m"
)

func colorRed(s string) string {
	return red + s + reset
}

func colorGreen(s string) string {
	return green + s + reset
}

func colorYello(s string) string {
	return yello + s + reset
}

func colorBlue(s string) string {
	return blue + s + reset
}

func colorWhite(s string) string {
	return white + s + reset
}

func levelColor(l slog.Level) string {
	levelStringColorMap := map[slog.Level]func(s string) string{
		slog.LevelInfo:  colorGreen,
		slog.LevelDebug: colorBlue,
		slog.LevelWarn:  colorYello,
		slog.LevelError: colorRed,
	}
	return levelStringColorMap[l](l.String())
}

type TextHandler struct {
	*slog.TextHandler
	mu *sync.Mutex
	w  io.Writer
}

func NewTextHandler(writer io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &TextHandler{
		TextHandler: slog.NewTextHandler(writer, opts),
		mu:          &sync.Mutex{},
		w:           writer,
	}
}

func (h *TextHandler) Handle(_ context.Context, record slog.Record) error {
	var buf = new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf(
		"[%s] level=%s message=%s ",
		record.Time.Format(time.DateTime),
		record.Level.String(),
		record.Message,
	))
	if record.NumAttrs() > 0 {
		record.Attrs(func(a slog.Attr) bool {
			if a.Key != "" {
				buf.WriteString(a.Key)
				buf.WriteString("=")
			}
			value := a.Value.String()
			if strings.Contains(value, "\n") {
				value = "\n" + value
			}
			buf.WriteString(value)
			buf.WriteByte(' ')
			return true
		})
	}
	buf.WriteString("\n")

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf.Bytes())
	return err
}
