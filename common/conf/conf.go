package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

// MustLoad loads the configuration from the given file path and scans it into the
// given Bootstrap struct. It panics if there are any errors loading or scanning
// the configuration. The returned function closes the underlying config source.
func MustLoad(out *Bootstrap, path string) func() error {
	c := config.New(config.WithSource(file.NewSource(path), env.NewSource()))

	if err := c.Load(); err != nil {
		panic(err)
	}

	if err := c.Scan(out); err != nil {
		panic(err)
	}

	// Validate all of the required fields are populated
	if err := out.ValidateAll(); err != nil {
		panic(err)
	}
	
	return c.Close
}
