package leveldb

import (
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/storage"
)

// "github.com/Maple-Open-Tech/monorepo/cloud/backend/config"

// NewIncomePropertyEvaluatorModuleEmailer creates a new emailer for the Income Property Evaluator module.
func NewIncomePropertyEvaluatorModuleEmailer(logger *zap.Logger) storage.Storage {
	configProvider := NewLevelDBConfigurationProvider(
		"xxx",
		"yyy",
	)

	return NewDiskStorage(configProvider.GetDBPath(), configProvider.GetDBName(), logger)
}
