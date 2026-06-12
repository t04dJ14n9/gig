package runner

import (
	"context"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// ProgramRunner is the execution surface consumed by the root gig package.
type ProgramRunner interface {
	Run(funcName string, params ...any) (any, error)
	RunWithContext(ctx context.Context, funcName string, params ...any) (any, error)
	RunWithValues(ctx context.Context, funcName string, args []value.Value) (value.Value, error)
	Wait()
	WaitContext(ctx context.Context) error
	InternalProgram() *bytecode.CompiledProgram
	Close()
}
