package repository

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const pqUniqueConstraintError = "23505"

var ErrRecordNotFound = errors.New("record not found")
var ErrUniqueConstraintFailure = errors.New("record already exist")
var ErrInvalidID = errors.New("invalid id")

func handleGormError(err error, logger *zap.Logger) error {
	fmt.Println(err)
	switch err {
	case gorm.ErrRecordNotFound:
		return ErrRecordNotFound
	}

	switch v := err.(type) {
	case *pq.Error:
		logger.Debug("PQerror",
			zap.String("PQSeverity", string(v.Severity)),
			zap.String("PQCode", string(v.Code)),
			zap.String("PQName", v.Code.Name()),
		)
		switch v.Code {
		case pqUniqueConstraintError:
			return ErrUniqueConstraintFailure
		default:
			return v
		}
	default:
		return err
	}
}
