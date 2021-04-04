package model

import (
	"context"
)

// CheckInterface is an interface to handle single check
type CheckInterface interface {
	Check(ctx context.Context) (bool, bool, error)
}
