//go:generate mockgen -source=interfaces.go -destination=interfaces_mock_test.go -package=method_test
package method

import "go.uber.org/zap"

type Logger interface {
	Warn(msg string, fields ...zap.Field)
}
