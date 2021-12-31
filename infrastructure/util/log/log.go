package log

import (
	appConfig "app/infrastructure/config"
	"app/infrastructure/util/ctxutil"
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var zapLogger *zap.Logger

type Field = zap.Field

func init() {
	zapLogger = newZapLogger()
}

// log with level

func Debug(msg string, fields ...zapcore.Field) {
	zapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	zapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	zapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	zapLogger.Error(msg, fields...)
}

func FError(msg string, fields ...zapcore.Field) {
	zapLogger.Error(msg, fields...)
	os.Exit(1)
}

func TraceId(traceId string) zapcore.Field {
	return String("traceId", traceId)
}

// construct field

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Float32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func Int(key string, val int) Field {
	return Int64(key, int64(val))
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Int64s(key string, vals []int64) Field {
	return zap.Int64s(key, vals)
}

func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Int16(key string, val int16) Field {
	return zap.Int16(key, val)
}

func Int8(key string, val int8) Field {
	return zap.Int8(key, val)
}

func String(key string, val string) Field {
	return zap.String(key, val)
}

func ErrMsg(msg string) Field {
	return zap.String("error", msg)
}

func TraceIdByCtx(ctx context.Context) Field {
	return zap.String("traceId", ctxutil.GetInTraceId(ctx))
}

func FormatStr(format string, val ...interface{}) string {
	return fmt.Sprintf(format, val...)
}

func Strings(key string, val []string) Field {
	return zap.Strings(key, val)
}

func Interface(key string, val interface{}) Field {
	return zap.Any(key, val)
}

func Uint(key string, val uint) Field {
	return Uint64(key, uint64(val))
}

func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

func Uint32(key string, val uint32) Field {
	return zap.Uint32(key, val)
}

func Uint16(key string, val uint16) Field {
	return zap.Uint16(key, val)
}

func Uint8(key string, val uint8) Field {
	return zap.Uint8(key, val)
}

func Reflect(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}

func ByteString(key string, val []byte) Field {
	return zap.ByteString(key, val)
}

func newZapLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.Level(appConfig.Get().Log.Level))
	config.DisableCaller = appConfig.Get().Log.DisableCaller
	zapLogger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func IsDebugLevel() (isDebug bool) {
	return appConfig.Get().Log.Level == -1
}
