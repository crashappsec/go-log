// log defines a logging intefrace and wraps uber-go/zap logger
package log

import (
	"time"

	"go.uber.org/zap"
)

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

func Err(val error) zap.Field {
	return zap.Error(val)
}

func ByteString(key string, val []byte) zap.Field {
	return zap.ByteString(key, val)
}

func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

func Strings(key string, val []string) zap.Field {
	return zap.Strings(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Time(key string, val time.Time) zap.Field {
	return zap.Time(key, val)
}

func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}
