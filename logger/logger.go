package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type logEntry struct {
	Time   time.Time
	Level  string
	Msg    string
	Status int
	Method string
	Path   string
	IP     string
	Cost   time.Duration
}

type AsyncLogger struct {
	ch       chan logEntry
	out      *bufio.Writer
	userJSON bool
	wg       sync.WaitGroup
	once     sync.Once
}

var std *AsyncLogger

func Init(w io.Writer, size int, userJSON bool) {
	std = NewAsyncLogger(w, size, userJSON)
}

func NewAsyncLogger(w io.Writer, size int, userJSON bool) *AsyncLogger {
	if size <= 0 {
		size = 1024
	}
	l := &AsyncLogger{
		ch:       make(chan logEntry, size),
		out:      bufio.NewWriter(w),
		userJSON: userJSON,
	}
	l.wg.Add(1)
	go l.run()
	return l
}

func (l *AsyncLogger) write(e logEntry) {
	select {
	case l.ch <- e:
	default:
	}
}

func (l *AsyncLogger) run() {
	defer l.wg.Done()

	flushTicker := time.NewTicker(3 * time.Second)
	defer flushTicker.Stop()

	for {
		select {
		case en, ok := <-l.ch:
			if !ok {
				return
			}
			_, _ = l.out.WriteString(l.format(en) + "\n")
		case <-flushTicker.C:
			_ = l.out.Flush()
		}
	}
}

func (l *AsyncLogger) format(e logEntry) string {
	if e.Msg != "" {
		if l.userJSON {
			b, _ := json.Marshal(struct {
				Time  string `json:"time"`
				Level string `json:"level"`
				Msg   string `json:"msg"`
			}{
				Time:  e.Time.Format(time.RFC3339),
				Level: e.Level,
				Msg:   e.Msg,
			})
			return string(b)
		}
		return fmt.Sprintf("%s [%s] %s",
			e.Time.Format("2006-01-02 15:04:05"), e.Level, e.Msg)
	}

	if l.userJSON {
		b, _ := json.Marshal(struct {
			Time   string `json:"time"`
			Level  string `json:"level"`
			Method string `json:"method"`
			Path   string `json:"path"`
			Status int    `json:"status"`
			CostMs int64  `json:"cost_ms"`
			IP     string `json:"ip"`
		}{
			Time:   e.Time.Format(time.RFC3339),
			Level:  e.Level,
			Method: e.Method,
			Path:   e.Path,
			Status: e.Status,
			CostMs: e.Cost.Milliseconds(),
			IP:     e.IP,
		})
		return string(b)
	}
	return fmt.Sprintf("%s [%s] %s %s %d %v",
		e.Time.Format("2006-01-02 15:04:05"), e.Level, e.Method, e.Path, e.Status, e.Cost)
}

func levelFromStatus(status int) string {
	switch {
	case status >= 500:
		return "error"
	case status >= 400:
		return "warn"
	default:
		return "info"
	}
}

func (l *AsyncLogger) GinMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		l.write(logEntry{
			Time:   start,
			Level:  levelFromStatus(ctx.Writer.Status()),
			Status: ctx.Writer.Status(),
			Method: ctx.Request.Method,
			Path:   ctx.Request.URL.Path,
			IP:     ctx.ClientIP(),
			Cost:   time.Since(start),
		})
	}
}

func (l *AsyncLogger) Close() {
	l.once.Do(func() {
		close(l.ch)
		l.wg.Wait()
		_ = l.out.Flush()
	})
}

func (l *AsyncLogger) Infof(format string, args ...any) {
	l.write(logEntry{
		Time:  time.Now(),
		Level: "info",
		Msg:   fmt.Sprintf(format, args...),
	})
}

func Infof(format string, args ...any) {
	std.Infof(format, args...)
}

func (l *AsyncLogger) Warnf(format string, args ...any) {
	l.write(logEntry{
		Time:  time.Now(),
		Level: "warn",
		Msg:   fmt.Sprintf(format, args...),
	})
}

func Warng(format string, args ...any) {
	std.Warnf(format, args...)
}

func (l *AsyncLogger) Errorf(format string, args ...any) {
	l.write(logEntry{
		Time:  time.Now(),
		Level: "error",
		Msg:   fmt.Sprintf(format, args...),
	})
}

func Errorf(format string, args ...any) {
	std.Errorf(format, args...)
}

func GinMiddleware() gin.HandlerFunc {
	return std.GinMiddleware()
}

func Close() {
	std.Close()
}
