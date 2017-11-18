package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/jsternberg/zap-logfmt"
	isatty "github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(w io.Writer) *zap.Logger {
	config := NewConfig()
	l, _ := config.New(w)
	return l
}

func (c *Config) New(defaultOutput io.Writer) (*zap.Logger, error) {
	w := defaultOutput
	format := c.Format
	if c.Format == "" || c.Format == "auto" {
		if isTerminal(w) {
			format = "console"
		} else {
			format = "logfmt"
		}
	}

	encoder, err := newEncoder(format)
	if err != nil {
		return nil, err
	}
	return zap.New(zapcore.NewCore(
		encoder,
		zapcore.AddSync(w),
		c.Level,
	)), nil
}

func newEncoder(format string) (zapcore.Encoder, error) {
	config := newEncoderConfig()
	switch format {
	case "json":
		return zapcore.NewJSONEncoder(config), nil
	case "console":
		return zapcore.NewConsoleEncoder(config), nil
	case "logfmt":
		return zaplogfmt.NewEncoder(config), nil
	default:
		return nil, fmt.Errorf("unknown logging format: %s", format)
	}
}

func newEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format(time.RFC3339))
	}
	config.EncodeDuration = func(d time.Duration, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(d.String())
	}
	return config
}

func isTerminal(w io.Writer) bool {
	if f, ok := w.(interface {
		Fd() uintptr
	}); ok {
		return isatty.IsTerminal(f.Fd())
	}
	return false
}
