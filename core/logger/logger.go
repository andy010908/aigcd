package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger(logpath string, loglevel string) {
	// 日志分割
	hook := lumberjack.Logger{
		Filename:   logpath, // 日志文件路径，默认 os.TempDir()
		MaxSize:    100,     //MAX 100M
		MaxBackups: 30,      // 保留30个备份，默认不限
		MaxAge:     7,       // 保留7天，默认不限
		Compress:   true,    // 是否压缩，默认不压缩
	}
	write := zapcore.AddSync(&hook)
	// 设置日志级别
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		// zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewJSONEncoder(encoderConfig),
		// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&write)), // 打印到控制台和文件
		write,
		level,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段,如：添加一个服务器名称
	filed := zap.Fields(zap.String("serviceName", "aigcd"))
	// 构造日志 如果不需要一些参数可以删除
	Logger = zap.New(core, caller, development, filed)
	//logger = zap.New(core, development)
	Logger.Info("DefaultLogger init success")
}

func Info(msg string, fields ...zapcore.Field) {
	// fields = append(fields, zap.String("time", time.Now().Format(time.RFC3339)))
	Logger.Info(msg, fields...)
}

// Debug logs Debug level msg
func Debug(msg string, fields ...zapcore.Field) {
	// fields = append(fields, zap.String("time", time.Now().Format(time.RFC3339)))
	Logger.Debug(msg, fields...)
}

// Warn logs Warn level msg
func Warn(msg string, fields ...zapcore.Field) {
	// fields = append(fields, zap.String("time", time.Now().Format(time.RFC3339)))

	Logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zapcore.Field) {
	// fields = append(fields, zap.String("time", time.Now().Format(time.RFC3339)))

	Logger.Error(msg, fields...)
}

/*
func main() {
	// 历史记录日志名字为：my.log，服务重新启动，日志会追加，不会删除
	InitLogger("./logs/my.log", "debug")
	// 强结构形式
	logger.Info("test",
		zap.String("string", "xiaotang"),
		zap.Int("int", 3),
		zap.Duration("time", time.Second),
	)

	// // 必须 key-value 结构形式 性能下降一点
	// logger.Sugar().Infow("test-",
	// 	"string", "kk",
	// 	"int", 1,
	// 	"time", time.Second,
	// )

	logger.Error("test02",
		zap.String("string", "x666g"),
		zap.Int("int", 4),
		zap.Duration("time", time.Second),
	)

	for {
		logger.Error("test02",
			zap.String("string", "x666g"),
			zap.Int("int", 4),
			zap.Duration("time", time.Second),
		)
		time.Sleep(5*time.Second)
	}
}
*/
