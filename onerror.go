package onerror

import (
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/jwmwalrus/bnp/httpstatus"
	log "github.com/sirupsen/logrus"
)

// LogHTTP logs and HTTP-related error
func LogHTTP(err error, r *http.Response, doNotCloseBody bool) error {
	if err != nil {
		if r != nil {
			withStatus(r.StatusCode, r.Status).Error(err)
		} else {
			log.Error(err)
		}
		return err
	} else if r != nil && httpstatus.IsError(r) {
		if !doNotCloseBody {
			defer r.Body.Close()
		}

		var b []byte
		b, _ = io.ReadAll(r.Body)
		msg := string(b)

		withStatus(r.StatusCode, r.Status, msg).Error(err)
		return fmt.Errorf("ERROR: %v\n\t%v", r.Status, msg)
	}
	return nil
}

// Log logs an error
func Log(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		withCaller(file, line).Error(err)
	}
}

// Panic asserts that no error was given
func Panic(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		wc := withCaller(file, line)
		wc.Error(err)
		wc.Fatal(err)
	}
}

// Warn warns on error
func Warn(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		withCaller(file, line).Warn(err)
	}
}

// WithEntry uses the given logrus.Entry
func WithEntry(e *log.Entry) *Entry {
	return &Entry{e}
}
