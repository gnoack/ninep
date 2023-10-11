package ninep

import "sync"

// fidPool is a thread-safe pool of FIDs.
//
// The main operations on fidPool are acquisition and release of FIDs.
type fidPool struct {
	mu sync.Mutex

	// TODO: This is an improper implementation.
	// We should track used FIDs instead of just cycling.
	nextFID uint32
}

func (p *fidPool) Acquire() uint32 {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.nextFID++
	return p.nextFID
}

func (p *fidPool) Release(fid uint32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// TODO: This is just a hack to allow for limited
	// pseudo-reuse.  We can return FIDs if they were the last
	// ones that were acquired.
	if p.nextFID == fid {
		p.nextFID--
	}
}
