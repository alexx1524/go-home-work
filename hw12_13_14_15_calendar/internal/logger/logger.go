package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
)

type Logger struct {
	filePath string
	logger   *zap.Logger
}

func New(level string, filePath string) (*Logger, error) {
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(zapConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig)

	var logLevel zapcore.Level
	if err := logLevel.Set(level); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, logLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		filePath: filePath,
		logger:   logger,
	}, nil
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Logger) Warning(msg string) {
	l.logger.Warn(msg)
}

func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *Logger) LogHTTPRequest(r *http.Request, statusCode int, duration time.Duration) {
	l.logger.Debug(fmt.Sprintf("%s, %v, %v", r.RequestURI, statusCode, duration))
}

func (l *Logger) LogGRPCRequest(code codes.Code, method, address string, requestDuration time.Duration) {
	l.logger.Debug(fmt.Sprintf("%v, %s, %s, %v", code, method, address, requestDuration))
}
