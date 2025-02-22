package js

import (
	"context"
	"strings"

	"github.com/dop251/goja"

	"github.com/sknv/protomock/pkg/log"
)

func ConsoleLog(ctx context.Context) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		var args strings.Builder

		for i, arg := range call.Arguments {
			if i > 0 {
				args.WriteByte(' ')
			}

			args.WriteString(arg.String())
		}

		log.FromContext(ctx).InfoContext(ctx, args.String())

		return nil
	}
}
