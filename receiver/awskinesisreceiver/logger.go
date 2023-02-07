package awskinesisreceiver

import (
	"github.com/vmware/vmware-go-kcl-v2/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type log struct {
	logger *zap.Logger
}

func (l log) Debugf(format string, args ...interface{}) {
	l.logger.Sugar().Debugf(format, args...)
}

func (l log) Infof(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}

func (l log) Warnf(format string, args ...interface{}) {
	l.logger.Sugar().Warnf(format, args...)
}

func (l log) Errorf(format string, args ...interface{}) {
	l.logger.Sugar().Errorf(format, args...)
}

func (l log) Fatalf(format string, args ...interface{}) {
	l.logger.Sugar().Fatalf(format, args...)
}

func (l log) Panicf(format string, args ...interface{}) {
	l.logger.Sugar().Fatalf(format, args...)
}

func (l log) WithFields(fields logger.Fields) logger.Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	l.logger.Sugar().With(f...)
	return l
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case logger.Info:
		return zapcore.InfoLevel
	case logger.Warn:
		return zapcore.WarnLevel
	case logger.Debug:
		return zapcore.DebugLevel
	case logger.Error:
		return zapcore.ErrorLevel
	case logger.Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
