//go:generate mockgen -source=interfaces.go -destination=logger_mock_test.go -package=interfaces_test
package interfaces

import "go.uber.org/zap"

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}
