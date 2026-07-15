package utils

import (
	"context"
	"fmt"
	"runtime/debug"
)

type ctxKey string

const RecoveryKey ctxKey = "recovery"

func Recover(ctx context.Context) {
	if rec := recover(); rec != nil {
		WithRecovery(ctx, rec)
	}
}

func WithRecovery(ctx context.Context, rec any) {
	fmt.Printf("panic recovered: %v\n%s", rec, string(debug.Stack()))
}

func Sleep(ctx context.Context) {

}
