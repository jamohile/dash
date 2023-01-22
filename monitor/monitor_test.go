package monitor

import (
	"testing"
	"time"
)

func TestMonitorInitializesNotDone(t *testing.T) {
	m := CreateMonitor()
	if m.IsDone() {
		t.Fatalf("Monitor should not be done.")
	}
}

func TestMonitorFiresSync(t *testing.T) {
	m := CreateMonitor()
	m.MarkDone()
	if !m.IsDone() {
		t.Fatalf("Monitor should be done.")
	}
}

func TestMonitorFiresAsync(t *testing.T) {
	m := CreateMonitor()

	go func() {
		time.Sleep(100 * time.Millisecond)
		m.MarkDone()
	}()

	start := time.Now()
	for time.Since(start) < time.Second {
		if m.IsDone() {
			return
		}
	}

	t.Fatalf("Monitor should have completed.")
}
