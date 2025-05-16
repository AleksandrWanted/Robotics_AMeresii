package err_stack

import (
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
)

var errorWithStackInterface = (*errbase.StackTraceProvider)(nil)

func WithStack(err error) error {
	if errors.HasInterface(err, errorWithStackInterface) {
		return err
	}
	return errors.WithStackDepth(err, 1)
}
