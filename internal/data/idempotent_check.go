package data

import (
	"context"
	"fmt"

	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/pkg/util"
)

type IdempotentCheck struct {
	data *Data
}

// IsIdempotent implements biz.IdempotencyChecker.
func (i *IdempotentCheck) IsIdempotent(ctx context.Context, uid int, req interface{}) (bool, error) {
	key := fmt.Sprintf("idempotent:%d:%s", uid, util.CalculateChecksum(req))
	if i.data.rdb.Exists(ctx, key).Val() == 1 {
		return true, nil
	}

	return false, nil
}

// MarkIdempotent implements biz.IdempotencyChecker.
func (i *IdempotentCheck) MarkIdempotent(ctx context.Context, uid int, req interface{}) error {
	key := fmt.Sprintf("idempotent:%d:%s", uid, util.CalculateChecksum(req))
	return i.data.rdb.Set(ctx, key, 1, 0).Err()
}

func NewIdempotentCheck(data *Data) *IdempotentCheck {
	return &IdempotentCheck{
		data: data,
	}
}

var _ biz.IdempotencyChecker = (*IdempotentCheck)(nil)
