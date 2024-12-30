package log

import (
    "fmt"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "os"
    "path/filepath"
    "runtime"
    "time"
)

var (
    Logger *DynamicLogger
)

type DynamicLogger struct {
    zapLogger *zap.Logger
    level     zap.AtomicLevel
}

func newLogger(logFile string, maxSize, maxBackups, maxAge int) (*DynamicLogger, error) {
    timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
        enc.AppendString(t.In(time.Local).Format("2006-01-02 15:04:05"))
    }

    level := zap.NewAtomicLevelAt(zap.InfoLevel)
    var core zapcore.Core

    lumberjackLogger := &lumberjack.Logger{
        Filename:   logFile,    // log path
        MaxSize:    maxSize,    // max size
        MaxBackups: maxBackups, // max backup count
        MaxAge:     maxAge,     // maximum number of days
        Compress:   true,       // compressed
    }
    core = zapcore.NewCore(
        zapcore.NewJSONEncoder(zapcore.EncoderConfig{
            TimeKey:      "timestamp",
            LevelKey:     "level",
            MessageKey:   "message",
            CallerKey:    "caller",
            EncodeTime:   timeEncoder,
            EncodeLevel:  zapcore.CapitalLevelEncoder,
            EncodeCaller: zapcore.ShortCallerEncoder,
        }),
        zapcore.AddSync(lumberjackLogger),
        level,
    )

    // create logger
    logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

    return &DynamicLogger{
        zapLogger: logger,
        level:     level,
    }, nil
}

// CheckWritePermission check write permission
func CheckWritePermission(path string) error {
    dir := filepath.Dir(path)
    info, err := os.Stat(dir)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("directory does not exist: %s", dir)
        }
        return fmt.Errorf("failed to stat directory: %s, error: %v", dir, err)
    }

    if !info.IsDir() {
        return fmt.Errorf("path is not a directory: %s", dir)
    }

    if info.Mode().Perm()&(1<<(uint(7))) == 0 {
        return fmt.Errorf("directory is not writable: %s", dir)
    }

    return nil
}

func init() {
    var logFilePath = os.Getenv("LOG_FILE_PATH")
    if logFilePath == "" {
        logFilePath = defaultPath()
    }
    err := CheckWritePermission(logFilePath)
    if err != nil {
        panic(err)
    }
    logger, err := newLogger(logFilePath, 100, 5, 30) // 文件路径，最大100MB，保留5个文件，保留30天
    if err != nil {
        panic(err)
    }
    logger.SetLogLevel(zap.DebugLevel)
    Logger = logger
}

func defaultPath() string {
    if runtime.GOOS == "windows" {
        return "C:\\Windows\\Temp\\transme.log"
    }
    return "/tmp/transme.log"
}

func (l *DynamicLogger) SetLogLevel(level zapcore.Level) {
    l.level.SetLevel(level)
}

func (l *DynamicLogger) Debug(msg string, fields ...zapcore.Field) {
    l.zapLogger.Debug(msg, fields...)
}

func (l *DynamicLogger) Info(msg string, fields ...zapcore.Field) {
    l.zapLogger.Info(msg, fields...)
}

func (l *DynamicLogger) Warn(msg string, fields ...zapcore.Field) {
    l.zapLogger.Warn(msg, fields...)
}

func (l *DynamicLogger) Error(msg string, fields ...zapcore.Field) {
    l.zapLogger.Error(msg, fields...)
    l.zapLogger.Sync()
}

func (l *DynamicLogger) DPanic(msg string, fields ...zapcore.Field) {
    l.zapLogger.DPanic(msg, fields...)
}

func (l *DynamicLogger) Panic(msg string, fields ...zapcore.Field) {
    l.zapLogger.Panic(msg, fields...)
}

func (l *DynamicLogger) Fatal(msg string, fields ...zapcore.Field) {
    l.zapLogger.Fatal(msg, fields...)
}
