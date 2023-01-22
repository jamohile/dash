package monitor

type Monitor struct {
	done chan bool
}

func CreateMonitor() Monitor {
	return Monitor{
		done: make(chan bool, 1),
	}
}

func (m *Monitor) MarkDone() {
	m.done <- true
}

func (m *Monitor) IsDone() bool {
	select {
	case <-m.done:
		return true
	default:
		return false
	}
}
