package utils

import "sync"

type Promise struct {
	wg sync.WaitGroup
}

func (p *Promise) Add(f func()) {
	p.wg.Add(1)
	go func() {
		f()
		defer p.wg.Done()
	}()
}

func (p *Promise) End() {
	p.wg.Wait()
}
