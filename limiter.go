package zif

import "time"

type Limiter struct {
	Throttle chan time.Time
	Ticker   *time.Ticker
}

func NewLimiter(rate time.Duration, burst int, fill bool) *Limiter {
	tick := time.NewTicker(rate)
	throttle := make(chan time.Time, burst)

	if fill {
		for i := 0; i < burst; i++ {
			throttle <- time.Now()
		}
	}

	go func() {
		for t := range tick.C {
			select {
			case throttle <- t:
			default:
			}
		}
	}()

	return &Limiter{throttle, tick}
}

func (l *Limiter) Wait() {
	_, _ = <-l.Throttle
}

func (l *Limiter) Stop() {
	l.Ticker.Stop()
	close(l.Throttle)
}

// Limits requests from peers
type PeerLimiter struct {
	queryLimiter    *Limiter
	announceLimiter *Limiter
}

func (pl *PeerLimiter) Setup() {
	// Allow an announce every 10 minutes, bursting to allow three.
	// The burst is there as people may make "mistakes" with titles or descriptions
	pl.announceLimiter = NewLimiter(time.Minute*10, 3, true)

	pl.queryLimiter = NewLimiter(time.Second/3, 3, true)
}
