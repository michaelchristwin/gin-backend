package validation

import (
	"context"
	"errors"
	"gin-backend/internal/domain"
	"gin-backend/internal/util"
	"log/slog"

	"modernc.org/sqlite"
)

var ()

func Validate(ctx context.Context, plaintext string, logger *slog.Logger) error {
	count, err := util.CheckPwned(ctx, plaintext)
	if err != nil {
		// fail open — log and continue rather than blocking the user
		logger.Warn("pwned password check failed", "err", err)
	} else if count > 0 {
		return domain.ErrPasswordPwned
	}

	return nil
}

func IsUniqueViolation(err error) bool {
	if sqliteErr, ok := errors.AsType[*sqlite.Error](err); ok {
		return sqliteErr.Code() == 2067 // SQLITE_CONSTRAINT_UNIQUE
	}
	return false
}
