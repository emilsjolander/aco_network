package main

import (
	"emergent_network/network"
	"flag"
	"fmt"
	"sync"
	"time"
)

var timeLimit int
var decayRate int
var packetsPerSecond int
var repetitions int

var lock sync.RWMutex

func main() {

	flag.IntVar(&timeLimit, "t", 60, "time in seconds to run the simulation")
	flag.IntVar(&decayRate, "d", 50, "rate of decay between 0 and 100")
	flag.IntVar(&packetsPerSecond, "pps", 1000, "packets sent per second")
	flag.IntVar(&repetitions, "r", 1, "number of the simulation should run")
	flag.Parse()

	network.EdgeWeightDecayRate = float64(decayRate) / 100
	network.SimulateSendingTime = true

	var wg sync.WaitGroup
	for i := 0; i < repetitions; i++ {
		wg.Add(1)
		go func() {
			simulate(false)
			wg.Done()
		}()
	}

	wg.Wait()
}

func simulate(showSimulation bool) {
	startNode := network.NewNetworkNode()

	nodes := make([]*network.NetworkNode, 6)

	nodes[0] = network.NewNetworkNode()
	startNode.ChildPaths = append(startNode.ChildPaths, network.NewEdge(startNode, nodes[0], 4))

	nodes[1] = network.NewNetworkNode()
	startNode.ChildPaths = append(startNode.ChildPaths, network.NewEdge(startNode, nodes[1], 3))

	nodes[2] = network.NewNetworkNode()
	startNode.ChildPaths = append(startNode.ChildPaths, network.NewEdge(startNode, nodes[2], 5))

	nodes[3] = network.NewNetworkNode()
	nodes[0].ChildPaths = append(nodes[0].ChildPaths, network.NewEdge(nodes[0], nodes[3], 2))
	nodes[1].ChildPaths = append(nodes[1].ChildPaths, network.NewEdge(nodes[1], nodes[3], 4))
	nodes[2].ChildPaths = append(nodes[2].ChildPaths, network.NewEdge(nodes[2], nodes[3], 3))

	nodes[4] = network.NewNetworkNode()
	nodes[0].ChildPaths = append(nodes[0].ChildPaths, network.NewEdge(nodes[0], nodes[4], 1))
	nodes[1].ChildPaths = append(nodes[1].ChildPaths, network.NewEdge(nodes[1], nodes[4], 3))
	nodes[2].ChildPaths = append(nodes[2].ChildPaths, network.NewEdge(nodes[2], nodes[4], 2))

	nodes[5] = network.NewNetworkNode()
	nodes[0].ChildPaths = append(nodes[0].ChildPaths, network.NewEdge(nodes[0], nodes[5], 3))
	nodes[1].ChildPaths = append(nodes[1].ChildPaths, network.NewEdge(nodes[1], nodes[5], 2))
	nodes[2].ChildPaths = append(nodes[2].ChildPaths, network.NewEdge(nodes[2], nodes[5], 6))

	endNode := network.NewEndNode()
	nodes[3].ChildPaths = append(nodes[3].ChildPaths, network.NewEdge(nodes[3], endNode, 4))
	nodes[4].ChildPaths = append(nodes[4].ChildPaths, network.NewEdge(nodes[4], endNode, 2))
	nodes[5].ChildPaths = append(nodes[5].ChildPaths, network.NewEdge(nodes[5], endNode, 1))

	quit := make(chan bool)
	go func() {
		time.Sleep(time.Duration(timeLimit) * time.Second)
		quit <- true
	}()

	for {
		select {
		case <-quit:
			lock.Lock()
			fmt.Printf("\t\t[%d]\t%.0f(%d)\t[%d]\n", 1, nodes[0].ChildPaths[0].Weight(), nodes[0].ChildPaths[0].Length, 4)
			fmt.Printf("\t\t\t%.0f(%d)\n", nodes[0].ChildPaths[1].Weight(), nodes[0].ChildPaths[1].Length)
			fmt.Printf("\t%.0f(%d)\t\t%.0f(%d)\t\t%.0f(%d)\n", startNode.ChildPaths[0].Weight(), startNode.ChildPaths[0].Length, nodes[0].ChildPaths[2].Weight(), nodes[0].ChildPaths[2].Length, nodes[3].ChildPaths[0].Weight(), nodes[3].ChildPaths[0].Length)
			fmt.Printf("\n")
			fmt.Printf("\t\t\t%.0f(%d)\n", nodes[1].ChildPaths[0].Weight(), nodes[1].ChildPaths[0].Length)
			fmt.Printf("S\t%.0f(%d)\t[%d]\t%.0f(%d)\t[%d]\t%.0f(%d)\tE\n", startNode.ChildPaths[1].Weight(), startNode.ChildPaths[1].Length, 2, nodes[1].ChildPaths[1].Weight(), nodes[1].ChildPaths[1].Length, 5, nodes[4].ChildPaths[0].Weight(), nodes[4].ChildPaths[0].Length)
			fmt.Printf("\t\t\t%.0f(%d)\n", nodes[1].ChildPaths[2].Weight(), nodes[1].ChildPaths[2].Length)
			fmt.Printf("\n")
			fmt.Printf("\t%.0f(%d)\t\t%.0f(%d)\t\t%.0f(%d)\n", startNode.ChildPaths[2].Weight(), startNode.ChildPaths[2].Length, nodes[2].ChildPaths[0].Weight(), nodes[2].ChildPaths[0].Length, nodes[5].ChildPaths[0].Weight(), nodes[5].ChildPaths[0].Length)
			fmt.Printf("\t\t\t%.0f(%d)\n", nodes[2].ChildPaths[1].Weight(), nodes[2].ChildPaths[1].Length)
			fmt.Printf("\t\t[%d]\t%.0f(%d)\t[%d]\n", 3, nodes[2].ChildPaths[2].Weight(), nodes[2].ChildPaths[2].Length, 6)

			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("\n")
			lock.Unlock()
			return
		default:
			startNode.Queue() <- network.NewPacket()
			if showSimulation {
				fmt.Printf("\t\t[%d]\t%.0f(%d)\t[%d]\n", 1, nodes[0].ChildPaths[0].Weight(), nodes[0].ChildPaths[0].Length, 4)
				fmt.Printf("\t\t\t%.0f(%d)\n", nodes[0].ChildPaths[1].Weight(), nodes[0].ChildPaths[1].Length)
				fmt.Printf("\t%.0f(%d)\t\t%.0f(%d)\t\t%.0f(%d)\n", startNode.ChildPaths[0].Weight(), startNode.ChildPaths[0].Length, nodes[0].ChildPaths[2].Weight(), nodes[0].ChildPaths[2].Length, nodes[3].ChildPaths[0].Weight(), nodes[3].ChildPaths[0].Length)
				fmt.Printf("\n")
				fmt.Printf("\t\t\t%.0f(%d)\n", nodes[1].ChildPaths[0].Weight(), nodes[1].ChildPaths[0].Length)
				fmt.Printf("S\t%.0f(%d)\t[%d]\t%.0f(%d)\t[%d]\t%.0f(%d)\tE\n", startNode.ChildPaths[1].Weight(), startNode.ChildPaths[1].Length, 2, nodes[1].ChildPaths[1].Weight(), nodes[1].ChildPaths[1].Length, 5, nodes[4].ChildPaths[0].Weight(), nodes[4].ChildPaths[0].Length)
				fmt.Printf("\t\t\t%.0f(%d)\n", nodes[1].ChildPaths[2].Weight(), nodes[1].ChildPaths[2].Length)
				fmt.Printf("\n")
				fmt.Printf("\t%.0f(%d)\t\t%.0f(%d)\t\t%.0f(%d)\n", startNode.ChildPaths[2].Weight(), startNode.ChildPaths[2].Length, nodes[2].ChildPaths[0].Weight(), nodes[2].ChildPaths[0].Length, nodes[5].ChildPaths[0].Weight(), nodes[5].ChildPaths[0].Length)
				fmt.Printf("\t\t\t%.0f(%d)\n", nodes[2].ChildPaths[1].Weight(), nodes[2].ChildPaths[1].Length)
				fmt.Printf("\t\t[%d]\t%.0f(%d)\t[%d]\n", 3, nodes[2].ChildPaths[2].Weight(), nodes[2].ChildPaths[2].Length, 6)

				fmt.Printf("\n")
				fmt.Printf("\n")
				fmt.Printf("\n")
				fmt.Printf("\n")
				fmt.Printf("\n")
			}
			time.Sleep(time.Duration(int64(time.Second) / int64(packetsPerSecond)))
		}
	}
}
