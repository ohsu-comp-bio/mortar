package cmd

import (
	"time"
)

// Counter logs periodically logs a count.
// e.g. for saying "imported N events" every second.
type Counter struct {
	ch chan int
	t  *time.Ticker
}

func NewCounter(msg string, d time.Duration) *Counter {
	countCh := make(chan int, 100)
	ticker := time.NewTicker(d)
	go func() {
		count := 0
		for {
			select {
			case n := <-countCh:
				count += n
			case <-ticker.C:
				if count != 0 {
					log.Info(msg, "count", count)
					count = 0
				}
			}
		}
	}()
	return &Counter{countCh, ticker}
}

func (c *Counter) Inc() {
	c.ch <- 1
}

func (c *Counter) Add(n int) {
	c.ch <- n
}

func (c *Counter) Stop() {
	c.t.Stop()
}
