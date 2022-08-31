package onerror

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"

	"github.com/jwmwalrus/bnp/httpstatus"
	log "github.com/sirupsen/logrus"
)

// Entry wraps logrus.Entry
type Entry struct {
	*log.Entry
}

// LogHTTP logs and HTTP-related error
func (e *Entry) LogHTTP(err error, r *http.Response, doNotCloseBody bool) error {
	if err != nil {
		if r != nil {
			e.withStatus(r.StatusCode, r.Status).Error(err)
		} else {
			e.Error(err)
		}
		return err
	} else if r != nil && httpstatus.IsError(r) {
		if !doNotCloseBody {
			defer r.Body.Close()
		}

		var b []byte
		b, _ = io.ReadAll(r.Body)
		msg := string(b)

		e.withStatus(r.StatusCode, r.Status, msg).Error(err)
		return fmt.Errorf("ERROR: %v\n\t%v", r.Status, msg)
	}
	return nil
}

// Log logs an error
func (e *Entry) Log(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		e.withCaller(file, line).Error(err)
	}
}

// Panic asserts that no error was given
func (e *Entry) Panic(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		wc := e.withCaller(file, line)
		wc.Error(err)
		wc.Fatal(err)
	}
}

// Warn warns on error
func (e *Entry) Warn(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		e.withCaller(file, line).Warn(err)
	}
}

// WithFields adds fields to the entry
func (e *Entry) WithFields(f log.Fields) *Entry {
	return &Entry{e.Entry.WithFields(f)}
}

func (e *Entry) withCaller(file string, line int) *Entry {
	f := callerFields(file, line)
	return e.WithFields(f)
}

func (e *Entry) withStatus(statusCode int, status string, msg ...string) *Entry {
	f := statusFields(statusCode, status, msg...)
	return e.WithFields(f)
}

func callerFields(file string, line int) log.Fields {
	return log.Fields{
		"caller":     file,
		"callerLine": line,
	}
}

func statusFields(statusCode int, status string, msg ...string) log.Fields {
	f := log.Fields{
		"statusCode": statusCode,
		"status":     status,
	}
	if len(msg) == 1 {
		f["error"] = msg
	} else if len(msg) > 0 {
		for i, m := range msg {
			f["msg("+strconv.Itoa(i+1)+")"] = m
		}
	}
	return f
}

func withCaller(file string, line int) *Entry {
	f := callerFields(file, line)
	return &Entry{log.WithFields(f)}
}

func withStatus(statusCode int, status string, msg ...string) *Entry {
	f := statusFields(statusCode, status, msg...)
	return &Entry{log.WithFields(f)}
}
