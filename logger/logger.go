package logger

import "go.uber.org/zap"

const ProductionMode = true
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}
type implLoggerService struct {
	zapLogger		*zap.Logger
	environmentMode	bool
}

func (service *implLoggerService) Debug(msg string, fields ...zap.Field) {
	if service.environmentMode == ProductionMode {
		service.zapLogger.Info(msg, fields...)
		return
	}
	service.zapLogger.Debug(msg, fields...)
	return
}

func (service *implLoggerService) Info(msg string, fields ...zap.Field) {
	service.zapLogger.Info(msg, fields...)
	return
}

func (service *implLoggerService) Warn(msg string, fields ...zap.Field) {
	service.zapLogger.Warn(msg, fields...)
	return
}

func (service *implLoggerService) Error(msg string, fields ...zap.Field) {
	service.zapLogger.Error(msg, fields...)
	return
}

func initZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"stdout",
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return logger, err
}

func NewLoggerService(envMode bool) (Logger, error) {
	logger, err := initZapLogger()
	if err != nil {
		return nil, err
	}
	return &implLoggerService{
		environmentMode: envMode,
		zapLogger: logger,
	}, nil
}
