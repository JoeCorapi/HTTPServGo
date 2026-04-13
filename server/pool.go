package server

import (
	"net"
)

type WorkerPool struct {
	queue   chan net.Conn
	maxSize int
}

func NewWorkerPool(maxWorkers int, queueSize int) *WorkerPool {
	pool := &WorkerPool{
		queue:   make(chan net.Conn, queueSize),
		maxSize: maxWorkers,
	}

	// Pre-spawn workers — they live forever, pulling from the queue
	for i := 0; i < maxWorkers; i++ {
		go pool.worker()
	}

	return pool
}

func (p *WorkerPool) worker() {
	for conn := range p.queue {
		HandleConnection(conn)
	}
}

func (p *WorkerPool) Submit(conn net.Conn) bool {
	select {
	case p.queue <- conn:
		return true
	default:
		// Queue is full — reject immediately, don't block
		return false
	}
}
