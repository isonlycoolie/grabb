package jobs

import "context"

// Start runs the workers until ctx is done.
func (m *Manager) Start(ctx context.Context) {
	for range m.workers {
		go m.work(ctx)
	}
	m.signal()
}

func (m *Manager) work(ctx context.Context) {
	for {
		if id, ok := m.next(); ok {
			m.execute(ctx, id)
			continue
		}
		select {
		case <-ctx.Done():
			return
		case <-m.wake:
		}
	}
}

func (m *Manager) next() (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for len(m.pending) > 0 {
		id := m.pending[0]
		m.pending = m.pending[1:]
		if j := m.jobs[id]; j != nil && j.Status == StatusQueued {
			if len(m.pending) > 0 {
				m.signal() // let another idle worker take the next one
			}
			return id, true
		}
	}
	return "", false
}

func (m *Manager) signal() {
	select {
	case m.wake <- struct{}{}:
	default:
	}
}
