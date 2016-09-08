package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"
	"io/ioutil"
	"net/http"

	"github.com/ronnrein/eps/database"
	"github.com/ronnrein/eps/utils"
	"bytes"
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
	Response   string
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
		split := strings.Split(r.RequestURI, "/")
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
			Response:   Result.Result.(string),
			Error:      errStr,
		}
		formattedTime := logEntry.Date.Format("02/Jan/2006 03:04:05 -0300")
		logStr := fmt.Sprintf(accessFormat, logEntry.Client, formattedTime, logEntry.Method, logEntry.Path, r.Proto, logEntry.StatusCode, unsafe.Sizeof(logEntry.Response))
		fmt.Fprintln(accessFile, logStr)
		fmt.Println(logStr)
		if Result.Error != nil {
			formattedTime = logEntry.Date.Format("Mon Jan _2 15:04:05 2006")
			logStr = fmt.Sprintf(errorFormat, formattedTime, logEntry.Client, logEntry.Error, logEntry.Path)
			fmt.Fprintln(errorFile, logStr)
			fmt.Println(logStr)
		}
		if split[len(split)-1] == "log" {
			return
		}
		query := database.DB.Create(&logEntry)
		if query.Error != nil {
			fmt.Printf("Error logging to DB: %s", query.Error)
		}
	})
}
