package zapLog

import (
	"github.com/SilentQianyi/file"
	"github.com/SilentQianyi/timeExtend"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Filename     string `json:"Filename"`
	MaxSize      int    `json:"MaxSize"`
	MaxAge       int    `json:"MaxAge"`
	MaxBackups   int    `json:"MaxBackups"`
	Compress     bool   `json:"Compress"`
	RotationTime int    `json:"RotationTime"`
}

func Init(cfg *Config) *zap.Logger {
	// 创建日志配置
	config := zap.NewProductionConfig()

	// 设置日志文件切割的策略
	lumberjackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,    // 每个日志文件的最大大小（MB）
		MaxBackups: cfg.MaxBackups, // 最多保留的旧日志文件数
		MaxAge:     cfg.MaxAge,     // 保留日志文件的最大天数
		Compress:   cfg.Compress,   // 是否压缩旧日志文件
	}
	writeSyncer := zapcore.AddSync(lumberjackLogger)

	// 创建核心日志器
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		writeSyncer,
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	// 创建带有天数切割的日志器
	logger := zap.New(core, zap.Hooks(func(entry zapcore.Entry) error {
		curTime := timeExtend.GetCurZero()
		createTime := timeExtend.GetZeroByTime(file.GetFileCreatedTime(lumberjackLogger.Filename))
		if !curTime.Equal(createTime) {
			lumberjackLogger.Rotate()
		}
		return nil
	}))

	// 记录日志
	logger.Info("log init success")

	// 关闭日志器
	//logger.Sync()
	return logger
}
