package js

import (
	"context"
	"strings"

	"github.com/dop251/goja"

	"github.com/sknv/protomock/pkg/log"
)

type Console struct {
	ctx context.Context //nolint:containedctx // should use a new instance for every evaluation
}

func NewConsole(ctx context.Context) Console {
	return Console{
		ctx: ctx,
	}
}

func (c Console) Log(call goja.FunctionCall) goja.Value { //nolint:ireturn // contract
	var args strings.Builder

	for i, arg := range call.Arguments {
		if i > 0 {
			args.WriteByte(' ')
		}

		args.WriteString(arg.String())
	}

	ctx := c.ctx
	log.FromContext(ctx).InfoContext(ctx, args.String())

	return nil
}
