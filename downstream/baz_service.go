package downstream

import (
	"context"
	"log"
	"time"
)

type Baz struct {
	Depth       int
	Cancels     bool
	CancelDepth int
	StagePause  time.Duration
}

func (b *Baz) Do(ctx context.Context) (map[string]interface{}, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if b.StagePause > 0 {
		<-time.After(b.StagePause)
	}

	if b.Depth <= 0 {
		out := map[string]interface{}{
			"message": "chain complete",
		}
		return out, nil
	}

	baz := &Baz{
		Depth:      b.Depth - 1,
		StagePause: b.StagePause,
	}

	switch {
	case b.Cancels:
		if b.CancelDepth == 0 {
			var cf context.CancelFunc
			ctx, cf = context.WithCancel(ctx)
			cf()
		} else {
			baz.Cancels = true
			baz.CancelDepth = b.CancelDepth - 1
		}
	}

	log.Printf("dispatching child baz %d\n", baz.Depth)
	return baz.Do(ctx)
}
