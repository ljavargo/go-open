package utils

import (
	"sync"
)

// concurrent limit worker
// task support func() and func() error
type Worker struct {
	rw     sync.RWMutex
	wg     sync.WaitGroup
	err    error
	closed bool

	pools chan struct{}
}

func (w *Worker) RunTask(task interface{}) {
	w.rw.RLock()
	if w.closed {
		w.rw.RUnlock()
		return
	}
	w.wg.Add(1)
	w.pools <- struct{}{}
	w.rw.RUnlock()

	go w.run(task)
}

func (w *Worker) Wait() (err error) {
	w.wg.Wait()
	err = w.err
	w.err = nil
	return
}

func (w *Worker) Close() {
	w.rw.Lock()
	w.close()
	w.rw.Unlock()
}

func (w *Worker) run(task interface{}) {
	switch t := task.(type) {
	case func():
		t()
		w.done()
	case func() error:
		err := t()
		w.done()
		if err != nil {
			w.rw.Lock()
			if !w.closed {
				w.close()
				w.err = err
			}
			w.rw.Unlock()
		}
	}
}

func (w *Worker) done() {
	<-w.pools
	w.wg.Done()
}

func (w *Worker) close() {
	if w.closed {
		return
	}
	w.closed = true
	close(w.pools)
}

func NewWorker(nums int) *Worker {
	wk := new(Worker)
	wk.pools = make(chan struct{}, nums)
	return wk
}
