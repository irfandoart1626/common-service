package log

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var (
	logInfo = EndInfo{
		ApiId:           "Test12345",
		HttpStatusCode:  "200",
		LogMessage:      "Test",
		LogPoint:        "api-test-logpoint",
		ResponsePayload: "",
		ServiceId:       "62123456789",
		TransactionId:   "Trx12345678",
	}

	exceptInfo = ExceptionInfo{
		ApiId:             "Test12345",
		TraceId:           "000000000000000000000000",
		ExceptionCategory: "BUSINESS",
		ExceptionCode:     "20002",
		ExceptionMessage:  "Invalid MSISDN",
		ExceptionSeverity: "1",
		HttpStatusCode:    400,
		ProcessTime:       "500",
		ServiceId:         "62123456789",
		ServiceName:       "Test-service",
		TransactionId:     "Trx123456789",
	}

	faultdetails = FaultDetails{
		Error:      "Error Exception",
		StackTrace: []string{"handler.test"},
	}
)

type LogSuite struct {
	suite.Suite
}

func (s *LogSuite) SetupSuite() {
}

func TestLogStart(t *testing.T) {
	suite.Run(t, new(LogSuite))
}

func (s *LogSuite) TestGetLogger() {
	assert.NotPanics(s.T(), func() { GetLogger() }, "not panic")
}

func (s *LogSuite) TestSet_and_GetLoggerCtx() {
	ctx := context.TODO()
	datas := make(map[string]interface{})
	datas["msisdn"] = "628123456789"
	datas["time"] = time.Now()
	SetLoggerCtxVal(ctx, "transactionId", "Trx123")
	SetLoggerCtxValInterface(ctx, "data", datas)
	SetLoggerCtxValLogInfo(ctx, logInfo)

	assert.NotEmpty(s.T(), GetLoggerCtx(context.TODO()))
}

func (s *LogSuite) TestSetLog() {
	LogTrace(logInfo)
	LogInfo(logInfo)
	LogException(exceptInfo, faultdetails, "")
}

func (s *LogSuite) TestBackendMetrics() {
	BackendMetrics("BackendSample", "BackendService", true, 3*time.Second)
}

func (s *LogSuite) TestProcessMetrics() {
	ProcessMetrics("ApiService", true, 3*time.Second)
}

func (s *LogSuite) TestWriter() {
	assert.NotEmpty(s.T(), NewLevelWriter())
	assert.NotEmpty(s.T(), NewConsoleLevelWriter())
}
