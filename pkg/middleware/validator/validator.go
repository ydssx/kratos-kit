package validator

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if v, ok := req.(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return nil, errors.BadRequest("INVALID_ARGUMENT", err.Error())
				}
			}
			return handler(ctx, req)
		}
	}
}

// SQLInjectionValidator checks for SQL injection attempts
func SQLInjectionValidator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if err := validateSQLInjection(req); err != nil {
				return nil, errors.BadRequest("INVALID_ARGUMENT", err.Error())
			}
			return handler(ctx, req)
		}
	}
}

func validateSQLInjection(req interface{}) error {
	// Convert request to JSON string for checking
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil // If we can't marshal, skip validation
	}
	jsonStr := string(jsonBytes)

	// SQL injection patterns
	patterns := []string{
		`(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|ALTER|CREATE|TRUNCATE).*FROM`,
		`(?i)(\-\-|\/\*|\*\/|;)`,
		`(?i)'.*OR.*'.*=.*'`,
		`(?i)'.*AND.*'.*=.*'`,
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, jsonStr)
		if err != nil {
			continue
		}
		if matched {
			return errors.New(400, "SQL_INJECTION", "potential SQL injection detected")
		}
	}

	return nil
}

// PathTraversalValidator checks for path traversal attempts
func PathTraversalValidator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if err := validatePathTraversal(req); err != nil {
				return nil, errors.BadRequest("INVALID_ARGUMENT", err.Error())
			}
			return handler(ctx, req)
		}
	}
}

func validatePathTraversal(req interface{}) error {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil
	}
	jsonStr := string(jsonBytes)

	// Path traversal patterns
	patterns := []string{
		`\.\./`,
		`\.\.\%2f`,
		`\.\.\%5c`,
		`\%2e\%2e/`,
		`\%2e\%2e\%5c`,
	}

	for _, pattern := range patterns {
		if strings.Contains(strings.ToLower(jsonStr), pattern) {
			return errors.New(400, "PATH_TRAVERSAL", "potential path traversal detected")
		}
	}

	return nil
}

// CommandInjectionValidator checks for command injection attempts
func CommandInjectionValidator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if err := validateCommandInjection(req); err != nil {
				return nil, errors.BadRequest("INVALID_ARGUMENT", err.Error())
			}
			return handler(ctx, req)
		}
	}
}

func validateCommandInjection(req interface{}) error {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil
	}
	jsonStr := string(jsonBytes)

	// Command injection patterns
	patterns := []string{
		`;.*`,
		`\|.*`,
		`\$\(.*\)`,
		`\` + "`" + `.*\` + "`" + ``,
		`&.*`,
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, jsonStr)
		if err != nil {
			continue
		}
		if matched {
			return errors.New(400, "COMMAND_INJECTION", "potential command injection detected")
		}
	}

	return nil
}
