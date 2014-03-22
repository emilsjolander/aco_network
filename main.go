package main

import (
	"emergent_network/network"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var timeLimit int
var decayRate int
var packetsPerSecond int
var repetitions int
var destinationFile string

func main() {

	flag.IntVar(&timeLimit, "t", 60, "time in seconds to run the simulation")
	flag.IntVar(&decayRate, "d", 50, "rate of decay between 0 and 100")
	flag.IntVar(&packetsPerSecond, "pps", 1000, "packets sent per second")
	flag.IntVar(&repetitions, "r", 1, "number of the simulation should run")
	flag.StringVar(&destinationFile, "csv", "", "csv file to save data to")
	flag.Parse()

	network.EdgeWeightDecayRate = float64(decayRate) / 100

	var err error
	var output *os.File
	if destinationFile != "" {
		output, err = os.Create(destinationFile)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := output.Close(); err != nil {
				panic(err)
			}
		}()
		output.Write([]byte("\"time_limit\",\"decay_rate\",\"packets_per_second\",\"repititions\",\"mean_throughput\"\n"))
	}

	result := make(chan string)
	var wg sync.WaitGroup

	go func() {
		for {
			if output != nil {
				output.Write([]byte(<-result))
			} else {
				fmt.Println(<-result)
			}
			wg.Done()
		}
	}()

	for i := 0; i < repetitions; i++ {
		wg.Add(1)
		go simulate(output == nil, result)
	}

	wg.Wait()
}

func simulate(visualize bool, result chan<- string) {
	startNode := network.NewNetworkNode()

	nodes := make([]*network.NetworkNode, 9)

	nodes[0] = network.NewNetworkNode()
	startNode.ChildPaths = append(startNode.ChildPaths, network.NewEdge(startNode, nodes[0], 1))

	nodes[1] = network.NewNetworkNode()
	startNode.ChildPaths = append(startNode.ChildPaths, network.NewEdge(startNode, nodes[1], 3))

	nodes[2] = network.NewNetworkNode()
	startNode.ChildPaths = append(startNode.ChildPaths, network.NewEdge(startNode, nodes[2], 5))

	nodes[3] = network.NewNetworkNode()
	nodes[0].ChildPaths = append(nodes[0].ChildPaths, network.NewEdge(nodes[0], nodes[3], 2))
	nodes[1].ChildPaths = append(nodes[1].ChildPaths, network.NewEdge(nodes[1], nodes[3], 1))
	nodes[2].ChildPaths = append(nodes[2].ChildPaths, network.NewEdge(nodes[2], nodes[3], 3))

	nodes[4] = network.NewNetworkNode()
	nodes[0].ChildPaths = append(nodes[0].ChildPaths, network.NewEdge(nodes[0], nodes[4], 6))
	nodes[1].ChildPaths = append(nodes[1].ChildPaths, network.NewEdge(nodes[1], nodes[4], 2))
	nodes[2].ChildPaths = append(nodes[2].ChildPaths, network.NewEdge(nodes[2], nodes[4], 4))

	nodes[5] = network.NewNetworkNode()
	nodes[0].ChildPaths = append(nodes[0].ChildPaths, network.NewEdge(nodes[0], nodes[5], 1))
	nodes[1].ChildPaths = append(nodes[1].ChildPaths, network.NewEdge(nodes[1], nodes[5], 3))
	nodes[2].ChildPaths = append(nodes[2].ChildPaths, network.NewEdge(nodes[2], nodes[5], 9))

	nodes[6] = network.NewNetworkNode()
	nodes[3].ChildPaths = append(nodes[3].ChildPaths, network.NewEdge(nodes[3], nodes[6], 2))
	nodes[4].ChildPaths = append(nodes[4].ChildPaths, network.NewEdge(nodes[4], nodes[6], 2))
	nodes[5].ChildPaths = append(nodes[5].ChildPaths, network.NewEdge(nodes[5], nodes[6], 2))

	nodes[7] = network.NewNetworkNode()
	nodes[3].ChildPaths = append(nodes[3].ChildPaths, network.NewEdge(nodes[3], nodes[7], 4))
	nodes[4].ChildPaths = append(nodes[4].ChildPaths, network.NewEdge(nodes[4], nodes[7], 3))
	nodes[5].ChildPaths = append(nodes[5].ChildPaths, network.NewEdge(nodes[5], nodes[7], 8))

	nodes[8] = network.NewNetworkNode()
	nodes[3].ChildPaths = append(nodes[3].ChildPaths, network.NewEdge(nodes[3], nodes[8], 2))
	nodes[4].ChildPaths = append(nodes[4].ChildPaths, network.NewEdge(nodes[4], nodes[8], 3))
	nodes[5].ChildPaths = append(nodes[5].ChildPaths, network.NewEdge(nodes[5], nodes[8], 1))

	endNode := network.NewEndNode()
	nodes[6].ChildPaths = append(nodes[6].ChildPaths, network.NewEdge(nodes[6], endNode, 4))
	nodes[7].ChildPaths = append(nodes[7].ChildPaths, network.NewEdge(nodes[7], endNode, 5))
	nodes[8].ChildPaths = append(nodes[8].ChildPaths, network.NewEdge(nodes[8], endNode, 2))

	go func() {
		var dead *network.NetworkNode
		for {
			time.Sleep(10 * time.Second)
			if dead != nil {
				dead.SetActive(true)
			}
			n := rand.Intn(9)
			dead = nodes[n]
			dead.SetActive(false)
		}
	}()

	quit := make(chan bool)
	go func() {
		time.Sleep(time.Duration(timeLimit) * time.Second)
		quit <- true
	}()

	start := time.Now()
	for {
		select {
		case <-quit:
			throughput := endNode.SuccessfullPackets / int(time.Since(start).Seconds())
			result <- fmt.Sprintf("\"%d\",\"%d\",\"%d\",\"%d\",\"%d\"\n", timeLimit, decayRate, packetsPerSecond, repetitions, throughput)
			return
		default:
			if visualize {
				fmt.Printf("\t\t[%d]\t%.0f(%d)\t[%d]\t%.0f(%d)\t[%d]\n", 1, nodes[0].ChildPaths[0].Weight(), nodes[0].ChildPaths[0].Length, 4, nodes[3].ChildPaths[0].Weight(), nodes[3].ChildPaths[0].Length, 7)
				fmt.Printf("\t\t\t%.0f(%d)\t\t%.0f(%d)\n", nodes[0].ChildPaths[1].Weight(), nodes[0].ChildPaths[1].Length, nodes[3].ChildPaths[1].Weight(), nodes[3].ChildPaths[1].Length)
				fmt.Printf("\t%.0f(%d)\t\t%.0f(%d)\t\t%.0f(%d)\t\t%.0f(%d)\n", startNode.ChildPaths[0].Weight(), startNode.ChildPaths[0].Length, nodes[0].ChildPaths[2].Weight(), nodes[0].ChildPaths[2].Length, nodes[3].ChildPaths[2].Weight(), nodes[3].ChildPaths[2].Length, nodes[6].ChildPaths[0].Weight(), nodes[6].ChildPaths[0].Length)
				fmt.Printf("\n")
				fmt.Printf("\t\t\t%.0f(%d)\t\t%.0f(%d)\n", nodes[1].ChildPaths[0].Weight(), nodes[1].ChildPaths[0].Length, nodes[4].ChildPaths[0].Weight(), nodes[4].ChildPaths[0].Length)
				fmt.Printf("S\t%.0f(%d)\t[%d]\t%.0f(%d)\t[%d]\t%.0f(%d)\t[%d]\t%.0f(%d)\tE\n", startNode.ChildPaths[1].Weight(), startNode.ChildPaths[1].Length, 2, nodes[1].ChildPaths[1].Weight(), nodes[1].ChildPaths[1].Length, 5, nodes[4].ChildPaths[1].Weight(), nodes[4].ChildPaths[1].Length, 8, nodes[7].ChildPaths[0].Weight(), nodes[7].ChildPaths[0].Length)
				fmt.Printf("\t\t\t%.0f(%d)\t\t%.0f(%d)\n", nodes[1].ChildPaths[2].Weight(), nodes[1].ChildPaths[2].Length, nodes[4].ChildPaths[2].Weight(), nodes[4].ChildPaths[2].Length)
				fmt.Printf("\n")
				fmt.Printf("\t%.0f(%d)\t\t%.0f(%d)\t\t%.0f(%d)\t\t%.0f(%d)\n", startNode.ChildPaths[2].Weight(), startNode.ChildPaths[2].Length, nodes[2].ChildPaths[0].Weight(), nodes[2].ChildPaths[0].Length, nodes[5].ChildPaths[0].Weight(), nodes[5].ChildPaths[0].Length, nodes[8].ChildPaths[0].Weight(), nodes[8].ChildPaths[0].Length)
				fmt.Printf("\t\t\t%.0f(%d)\t\t%.0f(%d)\n", nodes[2].ChildPaths[1].Weight(), nodes[2].ChildPaths[1].Length, nodes[5].ChildPaths[1].Weight(), nodes[5].ChildPaths[1].Length)
				fmt.Printf("\t\t[%d]\t%.0f(%d)\t[%d]\t%.0f(%d)\t[%d]\n", 3, nodes[2].ChildPaths[2].Weight(), nodes[2].ChildPaths[2].Length, 6, nodes[5].ChildPaths[2].Weight(), nodes[5].ChildPaths[2].Length, 9)

				fmt.Printf("\n")
				fmt.Printf("\n")
				fmt.Printf("\n")
				fmt.Printf("\n")
				fmt.Printf("\n")
			}
			startNode.Queue() <- network.NewPacket()
			time.Sleep(time.Duration(int64(time.Second) / int64(packetsPerSecond)))
		}
	}
}
