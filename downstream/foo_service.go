package downstream

import (
	"context"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

type Foo struct {
	Async       bool
	AsyncDies   int
	Timeout     time.Duration
	Depth       int
	Cancels     bool
	CancelDepth int
	StagePause  time.Duration
}

func (f *Foo) Do(ctx context.Context) (map[string]interface{}, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if f.Async {
		return f.DoAsync(ctx)
	}

	baz := &Baz{
		Depth:      f.Depth - 1,
		StagePause: f.StagePause,
	}

	switch {
	case f.Cancels:
		if f.CancelDepth == 0 {
			var cf context.CancelFunc
			ctx, cf = context.WithCancel(ctx)
			cf()
		} else {
			baz.Cancels = true
			baz.CancelDepth = f.CancelDepth - 1
		}
	case f.Timeout > 0:
		ctx, _ = context.WithTimeout(ctx, f.Timeout)
	}

	log.Printf("dispatching baz %d\n", baz.Depth)
	return baz.Do(ctx)
}

func (f *Foo) DoAsync(ctx context.Context) (map[string]interface{}, error) {
	errGroup, ctx := errgroup.WithContext(ctx)
	for i := 0; i < f.Depth; i++ {
		ii := i
		errGroup.Go(func() error {
			a := &Async{
				ID:      ii,
				MaxWait: f.StagePause,
				Dies:    f.AsyncDies == ii,
			}
			return a.Do(ctx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "async chain complete",
	}, nil
}
