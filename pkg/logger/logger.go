package logger

import (
	"gophermart/pkg/config"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Logger interface {
	Info(msg string, args ...LogMap)
	Warn(msg string, args ...LogMap)
	Error(msg string, args ...LogMap)
}

func NewSlogLogger(c *config.Config) Logger {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove time.
		if a.Key == slog.TimeKey && len(groups) == 0 {
			st := a.Value.Time()
			a.Value = slog.StringValue(st.Format(time.DateTime))
		}
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}
	return &logger{*slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace}))}
}

type LogMap map[string]any

type logger struct {
	slog.Logger
}

func (l *logger) Info(msg string, args ...LogMap) {
	ll := &l.Logger
	for _, arg := range args {
		for k, v := range arg {
			switch t := v.(type) {
			case int:
				ll = ll.With(slog.Int(k, t))
			case string:
				ll = ll.With(slog.String(k, t))
			case float64:
				ll = ll.With(slog.Float64(k, t))
			case bool:
				ll = ll.With(slog.Bool(k, t))
			case time.Duration:
				ll = ll.With(slog.Duration(k, t))
			case error:
				ll = ll.With(slog.String(k, t.Error()))
			default:
				log.Println("Unknown field for logger", k, v)
			}
		}
	}
	ll.Info(msg)
}

func (l *logger) Warn(msg string, args ...LogMap) {
	ll := &l.Logger
	for _, arg := range args {
		for k, v := range arg {
			switch t := v.(type) {
			case int:
				ll = ll.With(slog.Int(k, t))
			case string:
				ll = ll.With(slog.String(k, t))
			case float64:
				ll = ll.With(slog.Float64(k, t))
			case bool:
				ll = ll.With(slog.Bool(k, t))
			case time.Duration:
				ll = ll.With(slog.Duration(k, t))
			case error:
				ll = ll.With(slog.String(k, t.Error()))
			default:
				log.Println("Unknown field for logger", k, v)
			}
		}
	}
	ll.Warn(msg)
}

func (l *logger) Error(msg string, args ...LogMap) {
	ll := &l.Logger
	for _, arg := range args {
		for k, v := range arg {
			switch t := v.(type) {
			case int:
				ll = ll.With(slog.Int(k, t))
			case string:
				ll = ll.With(slog.String(k, t))
			case float64:
				ll = ll.With(slog.Float64(k, t))
			case bool:
				ll = ll.With(slog.Bool(k, t))
			case time.Duration:
				ll = ll.With(slog.Duration(k, t))
			case error:
				ll = ll.With(slog.String(k, t.Error()))
			default:
				log.Println("Unknown field for logger", k, v)
			}
		}
	}
	ll.Error(msg)
}
