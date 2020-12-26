package check

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"wait4it/pkg/model"
)

func RunCheck(ctx context.Context, c *model.CheckContext) error {
	cx, err := findCheckModule(c.Config.CheckType)
	if err != nil {
		return errors.Wrap(err, "can not find the module")
	}

	cx.BuildContext(*c)
	err = cx.Validate()
	if err != nil {
		return errors.Wrap(err, "validation failed")
	}

	fmt.Print("Wait4it...")

	newCtx, cnl := context.WithTimeout(ctx, time.Duration(c.Config.Timeout)*time.Second)
	defer cnl()

	if err := ticker(newCtx, cx); err != nil {
		return errors.Wrap(err, "check failed")
	}

	return nil
}

func findCheckModule(ct string) (model.CheckInterface, error) {
	m, ok := cm[ct]
	if !ok {
		return nil, errors.New("unsupported check type")
	}

	return m, nil
}

func ticker(ctx context.Context, cs model.CheckInterface) error {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			r, err := check(ctx, cs)
			if err != nil {
				return errors.Wrap(err, "check failed")
			}

			if r {
				return nil
			}

			wStdOut(false)
		}
	}
}

func check(ctx context.Context, cs model.CheckInterface) (bool, error) {
	r, eor, err := cs.Check(ctx)
	if err != nil && eor {
		return false, errors.Wrap(err, "failed")
	}

	return r, nil
}
