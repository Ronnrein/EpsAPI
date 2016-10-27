package middleware

import (
	"fmt"
	"os"
	"time"
	"unsafe"
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/ronnrein/eps/utils"
)

var accessFormat = "%s - - [%s] \"%s %s %s\" %d %d"
var errorFormat = "[%s] [error] [client %s] %s: %s"

var accessFile *os.File
var errorFile *os.File

type LogEntry struct {
	Date       time.Time
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
	Client     string
	Query      string
	Error      string
}

func init() {
	var err error
	accessFile, err = os.OpenFile(utils.Config.AccessLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	errorFile, err = os.OpenFile(utils.Config.ErrorLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
}

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		body, _ = ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		start := time.Now()
		inner.ServeHTTP(w, r)
		duration := time.Since(start)
		var errStr string
		if Result.Error != nil {
			errObj := *Result.Error
			errStr = errObj.Error()
		}
		logEntry := LogEntry{
			Date:       start,
			Method:     r.Method,
			Path:       r.RequestURI,
			StatusCode: Result.StatusCode,
			Duration:   duration,
			Client:     r.RemoteAddr,
			Query:      string(body),
			Error:      errStr,
		}
		formattedTimeLog := logEntry.Date.Format("02/Jan/2006 03:04:05 -0300")
		logStr := fmt.Sprintf(accessFormat, logEntry.Client, formattedTimeLog, logEntry.Method, logEntry.Path, r.Proto, logEntry.StatusCode, unsafe.Sizeof(Result.Result))
		fmt.Println(logStr)
		if utils.Config.LogAccess {
			fmt.Fprintln(accessFile, logStr)
		}
		if Result.Error != nil {
			formattedTimeError := logEntry.Date.Format("Mon Jan _2 15:04:05 2006")
			errorStr := fmt.Sprintf(errorFormat, formattedTimeError, logEntry.Client, logEntry.Error, logEntry.Path)
			fmt.Println(errorStr)
			if utils.Config.LogError {
				fmt.Fprintln(errorFile, errorStr)
			}
		}
	})
}
