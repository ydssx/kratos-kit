package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type validator interface {
	ValidateAll() error
}

func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(validator); ok {
				if err := v.ValidateAll(); err != nil {
					tr, ok := transport.FromServerContext(ctx)
					if ok {
						return nil, customizeErrorEncoder(tr, err)
					}
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}
			return handler(ctx, req)
		}
	}
}

func customizeErrorEncoder(tr transport.Transporter, err error) error {
	switch tr.Operation() {
	default:
		return errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
	}
}
