package tools

import (
	"github.com/jinzhu/gorm"
	"github.com/niuniumart/gosdk/martlog"
)

// RollbackUnlessCommitted func
func RollbackUnlessCommitted(tx *gorm.DB) {
	err := tx.RollbackUnlessCommitted().Error
	if err != nil {
		martlog.Errorf("tx.RollbackUnlessCommitted Error %s", err.Error())
	}
}

// Commit sql
func Commit(tx *gorm.DB) {
	err := tx.Commit().Error
	if err != nil {
		martlog.Errorf("tx.Commit Error %s", err.Error())
	}
}
