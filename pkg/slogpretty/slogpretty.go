package slogpretty

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l      *stdLog.Logger
	attrs  []slog.Attr
	groups []string
}

func (opts PrettyHandlerOptions) NewPrettyHandler(
	out io.Writer,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]any, r.NumAttrs())
	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	r.Attrs(func(a slog.Attr) bool {
		setNested(fields, h.groups, a.Key, a.Value.Any())
		return true
	})

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:04:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		"Comment:",
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
		groups:  h.groups,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	groups := make([]string, 0, len(h.groups)+1)
	groups = append(groups, h.groups...)
	groups = append(groups, name)

	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
		attrs:   h.attrs,
		groups:  groups,
	}
}

func setNested(fields map[string]any, groups []string, key string, val any) {
	curr := fields

	for _, g := range groups {
		if _, ok := curr[g]; !ok {
			curr[g] = make(map[string]any)
		}
		if nextMap, ok := curr[g].(map[string]any); ok {
			curr = nextMap
		}
	}

	curr[key] = val
}
