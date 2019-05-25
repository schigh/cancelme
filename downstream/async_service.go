package downstream

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

type Async struct {
	ID int
	MaxWait time.Duration
	Dies bool
}

func (a *Async) Do(ctx context.Context) error {
	waitDur := time.Duration(rand.Int63n(int64(a.MaxWait)))
	<-time.After(waitDur)
	if ctx.Err() != nil {
		log.Printf("async #%d terminated by context after %v: %v\n", a.ID, waitDur, ctx.Err())
		return ctx.Err()
	}
	if a.Dies {
		log.Printf("async #%d dies after %v\n", a.ID, waitDur)
		return errors.New("dies")
	}
	log.Printf("async #%d succeeded after %v\n", a.ID, waitDur)
	return nil
}
