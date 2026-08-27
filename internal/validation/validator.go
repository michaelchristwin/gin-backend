package validation

import (
	"context"
	"errors"
	"fmt"
)

const (
	MinLength = 12
	MaxLength = 64
)

var (
	ErrTooShort = fmt.Errorf("password must be at least %d characters", MinLength)
	ErrTooLong  = fmt.Errorf("password must be at most %d characters", MaxLength)
	ErrPwned    = errors.New("this password has appeared in a known data breach — please choose a different one")
)

func Validate(ctx context.Context, plaintext string) error {
	if len([]rune(plaintext)) < MinLength {
		return ErrTooShort
	}
	if len([]rune(plaintext)) > MaxLength {
		return ErrTooLong
	}

	count, err := CheckPwned(ctx, plaintext)
	if err != nil {
		// fail open — log and continue rather than blocking the user
		// logger.Warn("pwned password check failed", "err", err)
	} else if count > 0 {
		return ErrPwned
	}

	return nil
}
