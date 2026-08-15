package replay

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConsumeAtMostOnceConcurrent(t *testing.T) {
	s:=NewMemoryStore()
	key:=Key{OrganizationID:"org",RequestID:"req",Nonce:"nonce"}
	exp:=time.Now().Add(time.Minute)
	var success atomic.Int64
	var replayed atomic.Int64
	var wg sync.WaitGroup
	for i:=0;i<100;i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			err:=s.Consume(context.Background(),key,exp)
			switch {
			case err==nil: success.Add(1)
			case errors.Is(err,ErrReplay): replayed.Add(1)
			default: t.Errorf("unexpected error: %v",err)
			}
		}()
	}
	wg.Wait()
	if success.Load()!=1 { t.Fatalf("success=%d want=1",success.Load()) }
	if replayed.Load()!=99 { t.Fatalf("replay=%d want=99",replayed.Load()) }
}
