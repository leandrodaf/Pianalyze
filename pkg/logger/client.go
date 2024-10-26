package zaplogger

import (
	"sync"
	"time"

	"github.com/leandrodaf/pianalyze/internal/contracts/logger"
	"go.uber.org/zap"
)

// Pool for reusing ZapFields, saving memory allocation
var zapFieldPool = sync.Pool{
	New: func() interface{} {
		return &ZapField{}
	},
}

// ZapField encapsulates zap.Field and implements the logger.Field interface.
type ZapField struct {
	fields []zap.Field
}

func (z *ZapField) Bool(key string, val bool) logger.Field {
	z.fields = append(z.fields, zap.Bool(key, val))
	return z
}

func (z *ZapField) Int(key string, val int) logger.Field {
	z.fields = append(z.fields, zap.Int(key, val))
	return z
}

func (z *ZapField) Int64(key string, val int64) logger.Field {
	z.fields = append(z.fields, zap.Int64(key, val))
	return z
}

func (z *ZapField) Uint64(key string, val uint64) logger.Field {
	z.fields = append(z.fields, zap.Uint64(key, val))
	return z
}

func (z *ZapField) Uint8(key string, val uint8) logger.Field {
	z.fields = append(z.fields, zap.Uint8(key, val))
	return z
}

func (z *ZapField) Float64(key string, val float64) logger.Field {
	z.fields = append(z.fields, zap.Float64(key, val))
	return z
}

func (z *ZapField) String(key string, val string) logger.Field {
	z.fields = append(z.fields, zap.String(key, val))
	return z
}

func (z *ZapField) Time(key string, val time.Time) logger.Field {
	z.fields = append(z.fields, zap.Time(key, val))
	return z
}

// **Implementação do método Error**
func (z *ZapField) Error(key string, val error) logger.Field {
	z.fields = append(z.fields, zap.NamedError(key, val))
	return z
}

// ZapLogger implements the logger.Logger interface using zap.
type ZapLogger struct {
	logger *zap.Logger
}

var once sync.Once

// New creates and returns a new instance of ZapLogger with the provided options.
func New(options logger.Options) (logger.Logger, error) {
	var err error
	var zapLogger *zap.Logger

	once.Do(func() {
		if options.Environment {
			zapLogger, err = zap.NewProduction()
		} else {
			zapLogger, err = zap.NewDevelopment()
		}
	})

	if err != nil {
		return nil, err
	}

	return &ZapLogger{logger: zapLogger}, nil
}

// Field returns a Field instance (ZapField) from the pool for efficient use.
func (z *ZapLogger) Field() logger.Field {
	return zapFieldPool.Get().(*ZapField)
}

// Logging methods implementing the logger.Logger interface
func (z *ZapLogger) Info(msg string, fields ...logger.Field) {
	z.logger.Info(msg, z.convertFields(fields)...)
	z.releaseFields(fields)
}

func (z *ZapLogger) Error(msg string, fields ...logger.Field) {
	z.logger.Error(msg, z.convertFields(fields)...)
	z.releaseFields(fields)
}

func (z *ZapLogger) Debug(msg string, fields ...logger.Field) {
	z.logger.Debug(msg, z.convertFields(fields)...)
	z.releaseFields(fields)
}

func (z *ZapLogger) Warn(msg string, fields ...logger.Field) {
	z.logger.Warn(msg, z.convertFields(fields)...)
	z.releaseFields(fields)
}

// Fatal logs a fatal message and terminates the application.
func (z *ZapLogger) Fatal(msg string, fields ...logger.Field) {
	z.logger.Fatal(msg, z.convertFields(fields)...)
	z.releaseFields(fields)
}

// convertFields converts logger.Field to zap.Field.
func (z *ZapLogger) convertFields(fields []logger.Field) []zap.Field {
	var zapFields []zap.Field
	for _, field := range fields {
		if f, ok := field.(*ZapField); ok {
			zapFields = append(zapFields, f.fields...)
		}
	}
	return zapFields
}

// releaseFields releases the fields back to the pool after use.
func (z *ZapLogger) releaseFields(fields []logger.Field) {
	for _, field := range fields {
		if f, ok := field.(*ZapField); ok {
			// Clears the content and puts the field back in the pool
			f.fields = f.fields[:0]
			zapFieldPool.Put(f)
		}
	}
}
