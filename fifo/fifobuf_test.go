package fifo

import (
	"context"
	"testing"
	"time"
)

// TestFifoBuffer puhes
func TestFifoBuffer(t *testing.T) {
	set, head := NewBuffer()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	done := make(chan interface{})

	go func() {
		for i := 0; i < 10; i++ {
			h := <-head
			if h != i {
				t.Errorf("fifo violation: expected: %d, got: %d", i, h)
				return
			}
			set.Check()
		}

		close(done)
	}()

	for i := 0; i < 10; i++ {
		set.Add(i)
	}

	select {
	case <-done:
		return
	case <-ctx.Done():
		t.Fatal("test timeout")
	}
}
