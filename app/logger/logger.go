package logger

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nishanthrk/aspire-lms/app/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

var Sugar *zap.SugaredLogger
var App *newrelic.Application

type newRelicCore struct {
	zapcore.LevelEnabler
	enc zapcore.Encoder
	app *newrelic.Application
}

func (c *newRelicCore) With(fields []zapcore.Field) zapcore.Core {
	return &newRelicCore{
		LevelEnabler: c.LevelEnabler,
		enc:          c.enc,
		app:          c.app,
	}
}

func (c *newRelicCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, c)
	}
	return checkedEntry
}

func (c *newRelicCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	var sb strings.Builder
	sb.WriteString(entry.Message)
	for _, field := range fields {
		sb.WriteString(" ")
		sb.WriteString(field.Key)
		sb.WriteString("=")
		sb.WriteString(field.String)
	}

	c.app.RecordCustomEvent("Log", map[string]interface{}{
		"message": sb.String(),
		"level":   entry.Level.String(),
		"time":    entry.Time,
	})

	return nil
}

func (c *newRelicCore) Sync() error {
	return nil
}

func InitLogger() {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Initialize New Relic Application
	if len(configs.GetConfig().NewRelicLicense) > 0 {
		App, err = newrelic.NewApplication(
			newrelic.ConfigAppName(configs.GetConfig().GetTenantName()),
			newrelic.ConfigLicense(configs.GetConfig().NewRelicLicense),
			newrelic.ConfigDistributedTracerEnabled(true),
		)
		if err != nil {
			logger.Sugar().Info("failed to initialize New Relic application: %v", err)
		}

		// Create custom New Relic core
		nrCore := &newRelicCore{
			LevelEnabler: zap.InfoLevel,
			enc:          zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			app:          App,
		}

		// Combine the custom New Relic core with the default zap core
		core := zapcore.NewTee(
			logger.Core(), // default core for stdout and stderr
			nrCore,        // custom New Relic core
		)

		logger = zap.New(core)
	}

	Sugar = logger.Sugar()
}
