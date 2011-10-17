package markov

import "runtime"

func init() {
	runtime.GOMAXPROCS(2)
}

type Limiter struct {
	mar        *Markov
	order      int
	buffer     Sentence
	generating chan bool
}

func NewLimiter(m *Markov, order int) *Limiter {
	return &Limiter{
		mar:        m,
		order:      order,
		buffer:     m.GenerateN(order),
		generating: make(chan bool, 1),
	}
}

func (l *Limiter) Get() Sentence {
	select {
	case l.generating <- true:
		go l.generate()
	default:
	}
	return l.buffer
}

func (l *Limiter) generate() {
	l.buffer = l.mar.GenerateN(l.order)
	<-l.generating
}
