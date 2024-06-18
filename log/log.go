package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

var logger = zerolog.New(NewLevelWriter()).With().Logger()

type Property struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type ExceptionInfo struct {
	ApiID                 string `json:"apiID,omitempty"`
	ChannelID             string `json:"chanelID,omitempty"`
	TraceID               string `json:"traceId,omitempty"`
	ExceptionCategory     string `json:"exceptionCategory,omitempty"`
	ExceptionCode         string `json:"exceptionCode,omitempty"`
	ExceptionMessage      string `json:"exceptionMessage,omitempty"`
	ExceptionSeverity     string `json:"exceptionSeverity,omitempty"`
	HttpStatusCode        int    `json:"httpStatusCode,omitempty"`
	InternalTransactionID string `json:"internalTransactionID"`
	NotificationType      string `json:"notificationType,omitempty"`
	ProcessTime           string `json:"processTime,omitempty"`
	ServiceID             string `json:"serviceID"`
	ServiceName           string `json:"serviceName,omitempty"`
	TimeStamp             string `json:"timeStamp"`
	TransactionID         string `json:"transactionID"`
}

type EndInfo struct {
	LogTimestamp          time.Time `json:"-"`
	InternalTransactionID string    `json:"internalTransactionID"`
	TransactionID         string    `json:"transactionID"`
	ServiceID             string    `json:"serviceID"`
	ChannelID             string    `json:"chanelID,omitempty"`
	ApiID                 string    `json:"apiID,omitempty"`
	LogLevel              string    `json:"logLevel"`
	LogPoint              string    `json:"logPoint,omitempty"`
	LogMessage            string    `json:"logMessage,omitempty"`
	NotificationType      string    `json:"notificationType,omitempty"`
	RequestPayload        string    `json:"requestPayload,omitempty"`
	ResponsePayload       string    `json:"responsePayload,omitempty"`
	HttpStatusCode        string    `json:"httpStatusCode,omitempty"`
	ProcessTime           string    `json:"processTime,omitempty"`
	Timestamp             string    `json:"logTimestamp"`
}

// add faultdetails object for exception info
type FaultDetails struct {
	Error      string   `json:"error,omitempty"`
	StackTrace []string `json:"stackTrace,omitempty"`
}

func SetupLogger(devDebugMode bool) {
	if devDebugMode { // log a human-friendly output (not using json), and enabling http trace log
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		logger = log.Output(NewConsoleLevelWriter())
	}
}

func GetLevel(level string) zerolog.Level {
	// using switch case (not map) for performance reason
	switch strings.TrimSpace(strings.ToLower(level)) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disabled":
		return zerolog.Disabled
	}
	return zerolog.DebugLevel
}

// GetLogger will return logger that being used in this env
func GetLogger() zerolog.Logger {
	return logger
}

// GetLoggerCtx will return logger that is associated with provided ctx
func GetLoggerCtx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// SetLoggerCtxVal will set the logger context with provided key & val. If you have set the val, then the next logging of the same context will include this key/value in the output text
func SetLoggerCtxVal(ctx context.Context, key, val string) {
	l := GetLoggerCtx(ctx)
	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(key, val)
	})
}

// SetLoggerCtxValInterface will set the logger context with provided key & val(interface). If you have set the val, then the next logging of the same context will include this key/value in the output text
func SetLoggerCtxValInterface(ctx context.Context, key string, val interface{}) {
	l := GetLoggerCtx(ctx)
	if exceptInfo, ok := val.(ExceptionInfo); ok {
		val = exceptInfo
	}
	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Interface(key, val)
	})
}

// SetLoggerCtxValLogInfo will set the logger context with provided logInfo interface. If you have set the val, then the next logging of the same context will include this LogInfo in the output text
func SetLoggerCtxValLogInfo(ctx context.Context, logInfo EndInfo) {
	l := GetLoggerCtx(ctx)
	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		logInfo.Timestamp = logInfo.LogTimestamp.Format("2006-01-02T15:04:05.000-07:00")
		logPrinted := map[string]interface{}{
			"transactionID":         logInfo.TransactionID,
			"internalTransactionID": logInfo.InternalTransactionID,
			"logTimestamp":          logInfo.Timestamp,
			"serviceID":             logInfo.ServiceID,
			"channelID":             logInfo.ChannelID,
			"apiID":                 logInfo.ApiID,
			"logLevel":              logInfo.LogLevel,
			"logPoint":              logInfo.LogPoint,
			"logMessage":            logInfo.LogMessage,
			"notificationType":      logInfo.NotificationType,
			"requestPayload":        logInfo.RequestPayload,
			"responsePayload":       logInfo.ResponsePayload,
		}
		return c.Fields(logPrinted)
	})
}

// Debugf will print debug in stdout and give new line
func Debugf(format string, args ...interface{}) {
	if e := logger.Debug(); e.Enabled() {
		logger.Debug().Msgf(format, args...)
	}
}

// DebugfCtx will print debug in stdout using the logging keys/values in the context and give new line
func DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	if e := GetLoggerCtx(ctx).Debug(); e.Enabled() {
		e.Msgf(format, args...)
	}
}

// Infof will print info in stdout and give new line
func Infof(format string, args ...interface{}) {
	if e := logger.Info(); e.Enabled() {
		e.Msgf(format, args...)
	}
}

// InfofCtx will print info in stdout using the logging keys/values in the context and give new line
func InfofCtx(ctx context.Context, format string, args ...interface{}) {
	if e := GetLoggerCtx(ctx).Info(); e.Enabled() {
		e.Msgf(format, args...)
	}
}

// LogTrace will print array of key value pair object info in stdout and give new line
func LogTrace(EndInfo EndInfo) {
	if e := logger.Trace(); e.Enabled() {
		EndInfo.Timestamp = EndInfo.LogTimestamp.Format("2006-01-02T15:04:05.000-07:00")
		e.Interface("apiID", EndInfo.ApiID)
		e.Interface("channelID", EndInfo.ChannelID)
		e.Interface("httpStatusCode", EndInfo.HttpStatusCode)
		e.Interface("logLevel", EndInfo.LogLevel)
		e.Interface("logMessage", EndInfo.LogMessage)
		e.Interface("logPoint", EndInfo.LogPoint)
		e.Interface("logTimestamp", EndInfo.Timestamp)
		e.Interface("notificationType", EndInfo.NotificationType)
		e.Interface("requestPayload", EndInfo.RequestPayload)
		e.Interface("responsePayload", EndInfo.ResponsePayload)
		e.Interface("serviceID", EndInfo.ServiceID)
		e.Interface("transactionID", EndInfo.TransactionID)
		e.Interface("internalTransactionID", EndInfo.InternalTransactionID)
		e.Interface("processTime", EndInfo.ProcessTime)
		e.Msg(EndInfo.LogPoint)
	}
}

// LogInfo will print array of key value pair object info in stdout and give new line
func LogInfo(EndInfo EndInfo) {
	if i := logger.Info(); i.Enabled() {
		EndInfo.Timestamp = EndInfo.LogTimestamp.Format("2006-01-02T15:04:05.000-07:00")
		i.Interface("apiID", EndInfo.ApiID)
		i.Interface("channelID", EndInfo.ChannelID)
		i.Interface("httpStatusCode", EndInfo.HttpStatusCode)
		i.Interface("logLevel", EndInfo.LogLevel)
		i.Interface("logMessage", EndInfo.LogMessage)
		i.Interface("logPoint", EndInfo.LogPoint)
		i.Interface("logTimestamp", EndInfo.Timestamp)
		i.Interface("notificationType", EndInfo.NotificationType)
		i.Interface("requestPayload", EndInfo.RequestPayload)
		i.Interface("responsePayload", EndInfo.ResponsePayload)
		i.Interface("serviceID", EndInfo.ServiceID)
		i.Interface("transactionID", EndInfo.TransactionID)
		i.Interface("internalTransactionID", EndInfo.InternalTransactionID)
		i.Interface("processTime", EndInfo.ProcessTime)
		i.Msg(EndInfo.LogPoint)
	}
}

// LogWithoutLvl will print array of key value pair object info in stdout and give new line without key level
func LogWithoutLvl(endInfo *EndInfo, exceptionInfo *ExceptionInfo, details *FaultDetails, payload any) {
	if i := logger.WithLevel(zerolog.NoLevel); i.Enabled() {
		if endInfo != nil {
			endInfo.Timestamp = endInfo.LogTimestamp.Format("2006-01-02T15:04:05.000-07:00")
			i.Interface("apiID", endInfo.ApiID)
			i.Interface("channelID", endInfo.ChannelID)
			i.Interface("httpStatusCode", endInfo.HttpStatusCode)
			i.Interface("logLevel", endInfo.LogLevel)
			i.Interface("logMessage", endInfo.LogMessage)
			i.Interface("logPoint", endInfo.LogPoint)
			i.Interface("logTimestamp", endInfo.Timestamp)
			i.Interface("notificationType", endInfo.NotificationType)
			i.Interface("requestPayload", endInfo.RequestPayload)
			i.Interface("responsePayload", endInfo.ResponsePayload)
			i.Interface("serviceID", endInfo.ServiceID)
			i.Interface("transactionID", endInfo.TransactionID)
			i.Interface("internalTransactionID", endInfo.InternalTransactionID)
			i.Interface("processTime", endInfo.ProcessTime)
			i.Msg("")
		}
		if exceptionInfo != nil {
			i.Interface("ExceptionInfo", exceptionInfo)
			i.Interface("FaultDetails", details)
			i.Interface("requestPayload", payload)
		}

	}
}

// Warnf will print warn in stdout and give new line
func Warnf(format string, args ...interface{}) {
	if e := logger.Warn(); e.Enabled() {
		e.Msgf(format, args...)
	}
}

// WarnfCtx will print info in stdout using the logging keys/values in the context and give new line
func WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	if e := GetLoggerCtx(ctx).Warn(); e.Enabled() {
		e.Msgf(format, args...)
	}
}

// Errorf will print error in stderr and give new line
func Errorf(format string, args ...interface{}) {
	if e := logger.Error(); e.Enabled() {
		e.Msgf(format, args...)
	}
}

// ErrorfCtx will print info in stdout using the logging keys/values in the context and give new line
func ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	if e := GetLoggerCtx(ctx).Error(); e.Enabled() {
		e.Msgf(format, args...)
	}
}

// LogException will print error in stderr and give new line
func LogException(val ExceptionInfo, details FaultDetails, payload string) {
	if e := logger.Error(); e.Enabled() {
		logPoint := val.ApiID + "-" + val.ServiceName + "-End"

		e.Interface("ExceptionInfo", val)
		e.Interface("FaultDetails", details)
		e.Interface("requestPayload", payload)
		e.Msg(logPoint)
	}
}

// Fatalf will print error in stderr and give new line
func Fatalf(format string, args ...interface{}) {
	if e := logger.Fatal(); e.Enabled() {
		e.Msgf(format, args...)
	}
}

func TraceHttpReq(req *http.Request) {
	if e := logger.Trace(); e.Enabled() {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			Errorf("error on dump request %v", err)
			return
		}
		e.Msg(string(dump))
	}
}

func TraceHttpResp(resp *http.Response) {
	if e := logger.Trace(); e.Enabled() {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			Errorf("error on dump response %v", err)
			return
		}
		e.Msg(string(dump))
	}
}

func TraceFastHttpReq(req *fasthttp.Request) {
	if e := logger.Trace(); e.Enabled() {
		var buffer bytes.Buffer
		req.Header.VisitAll(func(key, value []byte) {
			buffer.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		})

		dump := []map[string]string{
			{
				"Method":  string(req.Header.Method()),
				"URI":     string(req.URI().FullURI()),
				"Headers": buffer.String(),
				"Body":    string(req.Body()),
			},
		}

		jsonDump, err := json.Marshal(dump)
		if err != nil {
			Errorf("error on dump request %v", err)
			return
		}

		e.Msg(string(jsonDump))
	}
}

func TraceFastHttpResp(resp *fasthttp.Response) {
	if e := logger.Trace(); e.Enabled() {
		var buffer bytes.Buffer
		resp.Header.VisitAll(func(key, value []byte) {
			buffer.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		})

		dump := []map[string]string{
			{
				"StatusCode": strconv.Itoa(resp.StatusCode()),
				"Headers":    buffer.String(),
				"Body":       string(resp.Body()),
			},
		}

		jsonDump, err := json.Marshal(dump)
		if err != nil {
			Errorf("error on dump response %v", err)
			return
		}

		e.Msg(string(jsonDump))
	}
}
