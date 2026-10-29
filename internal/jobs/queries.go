package jobs

import "slices"

// List returns every job, newest first.
func (m *Manager) List() []Job {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.listLocked()
}

// Get returns one job.
func (m *Manager) Get(id string) (Job, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	j, ok := m.jobs[id]
	if !ok {
		return Job{}, false
	}
	return *j, true
}

func (m *Manager) listLocked() []Job {
	list := make([]Job, 0, len(m.jobs))
	for _, j := range m.jobs {
		list = append(list, *j)
	}
	slices.SortFunc(list, func(a, b Job) int { return b.CreatedAt.Compare(a.CreatedAt) })
	return list
}
