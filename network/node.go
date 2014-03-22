package network

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Node interface {
	Queue() chan<- *Packet
	SetActive(activate bool)
	Active() bool
}

const MaxQueueSize = 10

type NetworkNode struct {
	ChildPaths []*Edge
	queue      chan *Packet
	destroy    chan bool
	active     bool
}

func NewNetworkNode() *NetworkNode {
	n := &NetworkNode{
		queue:   make(chan *Packet, MaxQueueSize),
		destroy: make(chan bool),
		active:  true,
	}
	go n.processEvents()
	return n
}

func (n *NetworkNode) Queue() chan<- *Packet {
	return n.queue
}

func (n *NetworkNode) SetActive(active bool) {
	if active {
		go n.processEvents()
	} else {
		n.destroy <- true
	}
	n.active = active
}

func (n *NetworkNode) Active() bool {
	return n.active
}

func (node *NetworkNode) processEvents() {
	for {
		select {
		case <-node.destroy:
			// stop processing events
			return
		case packet := <-node.queue:
			a := 1.0
			b := 1.0
			var total float64
			pValues := make([]float64, len(node.ChildPaths))
			for _, e := range node.ChildPaths {
				n := 1.0 / float64(e.Length)
				t := e.Weight()
				total += math.Pow(t, a) * math.Pow(n, b)
			}
			for i, e := range node.ChildPaths {
				n := 1.0 / float64(e.Length)
				t := e.Weight()
				pValues[i] = (math.Pow(t, a) * math.Pow(n, b)) / total
			}
			rouletteWheel := make([]int, 0, 1000)
			for i, p := range pValues {
				for j := 0; j < int(p*1000); j++ {
					rouletteWheel = append(rouletteWheel, i)
				}
			}
			targetIndex := rand.Intn(len(rouletteWheel))
			target := node.ChildPaths[rouletteWheel[targetIndex]]

			target.Send(packet)
		}
	}
}

type EndNode struct {
	queue              chan *Packet
	SuccessfullPackets int
}

func NewEndNode() *EndNode {
	n := &EndNode{
		queue:              make(chan *Packet, MaxQueueSize),
		SuccessfullPackets: 0,
	}
	go n.processEvents()
	return n
}

func (n *EndNode) Queue() chan<- *Packet {
	return n.queue
}

func (n *EndNode) processEvents() {
	for p := range n.queue {
		var length int
		for _, e := range p.Path {
			length += e.Length
		}
		for _, e := range p.Path {
			e.Strengthen(float64(1) / float64(length))
		}
		n.SuccessfullPackets++
	}
}

func (n *EndNode) SetActive(active bool) {
	// i'm always active
}

func (n *EndNode) Active() bool {
	return true
}
