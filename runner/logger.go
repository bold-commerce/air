package runner

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	rawColor = "raw"
	// TODO: support more colors
	colorMap = map[string]color.Attribute{
		"red":     color.FgRed,
		"green":   color.FgGreen,
		"yellow":  color.FgYellow,
		"blue":    color.FgBlue,
		"magenta": color.FgMagenta,
		"cyan":    color.FgCyan,
		"white":   color.FgWhite,
	}
)

type logFunc func(string, ...interface{})

type logger struct {
	config  *config
	colors  map[string]string
	loggers map[string]logFunc
}

func newLogger(cfg *config) *logger {
	if cfg == nil {
		return nil
	}

	colors := cfg.colorInfo()
	loggers := make(map[string]logFunc, len(colors))
	for name, nameColor := range colors {
		loggers[name] = newLogFunc(nameColor, cfg.Log)
	}
	loggers["default"] = defaultLogger()
	return &logger{
		config:  cfg,
		colors:  colors,
		loggers: loggers,
	}
}

func newLogFunc(colorname string, cfg cfgLog) logFunc {
	return func(msg string, v ...interface{}) {
		msg = strings.TrimSuffix(msg, "\n")
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			return
		}
		// TODO: filter msg by regex
		msg = msg + "\n"
		if cfg.AddTime {
			t := time.Now().Format("15:04:05")
			msg = fmt.Sprintf("[%s] %s", t, msg)
		}
		if colorname == rawColor {
			fmt.Fprintf(os.Stdout, msg, v...)
		} else {
			color.New(getColor(colorname)).Fprintf(color.Output, msg, v...)
		}
	}
}

func getColor(name string) color.Attribute {
	if v, ok := colorMap[name]; ok {
		return v
	}
	return color.FgWhite
}

func (l *logger) main() logFunc {
	return l.getLogger("main")
}

func (l *logger) build() logFunc {
	return l.getLogger("build")
}

func (l *logger) runner() logFunc {
	return l.getLogger("runner")
}

func (l *logger) watcher() logFunc {
	return l.getLogger("watcher")
}

func (l *logger) app() logFunc {
	return l.getLogger("app")
}

func rawLogger() logFunc {
	return newLogFunc("raw", defaultConfig().Log)
}

func defaultLogger() logFunc {
	return newLogFunc("white", defaultConfig().Log)
}

func (l *logger) getLogger(name string) logFunc {
	v, ok := l.loggers[name]
	if !ok {
		return rawLogger()
	}
	return v
}

type appLogWriter struct {
	l logFunc
}

func newAppLogWriter(l *logger) appLogWriter {
	return appLogWriter{
		l: l.app(),
	}
}

func (w appLogWriter) Write(data []byte) (n int, err error) {
	w.l(string(data))
	return len(data), nil
}
