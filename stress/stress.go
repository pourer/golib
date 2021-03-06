package stress

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type waitGroupWrapper struct {
	wg sync.WaitGroup
}

func (this *waitGroupWrapper) Wrap(seq int, cb func(seq int)) {
	this.wg.Add(1)
	atomic.AddInt32(&concurrency, 1)
	go func() {
		cb(seq)
		this.wg.Done()
		atomic.AddInt32(&concurrency, -1)
	}()
}

func (this *waitGroupWrapper) Wait() {
	this.wg.Wait()
	atomic.StoreInt32(&concurrency, 0)
}

func RunStress(cb func(seq int)) {
	go runReporter()

	var waitGroup waitGroupWrapper
	var t0 = time.Now()
	for c := Flags.C1; c <= Flags.C2; c += Flags.Step {
		for r := 0; r < Flags.Round; r++ {
			if !Flags.Neat {
				log.Printf("concurrency: %6d started, loops: %d", c, r)
			}
			t1 := time.Now()
			for i := 0; i < c; i++ {
				waitGroup.Wrap(i, cb)
			}
			waitGroup.Wait()
			if !Flags.Neat {
				log.Printf("concurrency: %6d elapsed: %s", c, time.Since(t1))
			}
		}

	}

	log.Printf("All done in %s", time.Since(t0))
}
