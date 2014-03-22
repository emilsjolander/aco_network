package network

import (
	"sync"
	"time"
)

var EdgeWeightDecayRate = 0.5
var SimulateSendingTime = true

type Edge struct {
	sync.RWMutex
	From   Node
	To     Node
	Length int
	weight float64
}

func NewEdge(from, to Node, length int) *Edge {
	e := &Edge{
		From:   from,
		To:     to,
		Length: length,
		weight: 1.0,
	}
	go e.decay()
	return e
}

func (e *Edge) decay() {
	for {
		time.Sleep(time.Second)
		e.Lock()
		e.weight *= (1 - EdgeWeightDecayRate)
		if e.weight < 1 {
			e.weight = 1
		}
		e.Unlock()
	}
}

func (e *Edge) Send(packet *Packet) {
	go func() {
		if SimulateSendingTime {
			time.Sleep(time.Duration(e.Length) * time.Millisecond)
		}
		if e.To.Active() {
			select {
			case e.To.Queue() <- packet:
				packet.Path = append(packet.Path, e)
			default:
				// throw packet
				return
			}
		}
	}()
}

func (e *Edge) Strengthen(amount float64) {
	e.Lock()
	defer e.Unlock()
	e.weight += amount
}

func (e *Edge) Weight() float64 {
	e.RLock()
	defer e.RUnlock()
	return e.weight
}
