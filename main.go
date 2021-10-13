package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var totalRounds float32
var nodeCount int
var turns int
var roundCount int
var nodes []Node
var gossipType int


type Node struct {
	isInfected bool
	channel  *chan bool
}


func pushSend(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()

	if node.isInfected == true {
		nextNode := rand.Intn(nodeCount)
		fmt.Printf("\nNode %d will be infected by push next.", nextNode)
		*nodes[nextNode].channel <- node.isInfected
	}
}

func pushReceive(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()

	select {
	case msg, temp := <- *node.channel:
		if temp {
			node.isInfected = msg
		} else {
			fmt.Println("Error: node could not receive.")
			break
		}
	default:
		break
	}
}

func push(wg *sync.WaitGroup) {
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		pushSend(wg, &nodes[i])
	}
	time.Sleep(50 * time.Millisecond)
	wg.Wait()

	for j := 0; j < nodeCount; j++ {
		wg.Add(1)
		pushReceive(wg, &nodes[j])
	}
	wg.Wait()
}

func runPush(wg *sync.WaitGroup) {
	totalRounds = 0
	for i:= 1; i<=turns; i++{
		fmt.Printf("\nTurn number %d!.\n", i)
		nodes = make([]Node, nodeCount)
		channels := make([]chan bool, nodeCount)

		for i := 0; i < nodeCount; i++ {
			channels[i] = make(chan bool, nodeCount)
			nodes[i] = Node{false, &(channels[i])}
		}

		nodes[0].isInfected = true
		roundCount = 0
		for true {
			roundCount++
			fmt.Printf("\n\nStarting round %d in turn %d!", roundCount, i)
			push(wg)
			complete, _ := isDone()
			if complete {
				totalRounds += (float32)(roundCount)
				break
			}
		}
		fmt.Printf("\nIt took %d rounds to infect %d nodes.\n", roundCount, nodeCount)
	}
	fmt.Printf("On average, over the %d turns, it took %f rounds to infect %d nodes.", turns, totalRounds/(float32)(turns), nodeCount)
}

// very similar code to pushReceive
func pullSend(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()

	if node.isInfected == false {
		nextNode := rand.Intn(nodeCount)
		select {
		case msg, temp := <- *nodes[nextNode].channel:
			if temp {
				node.isInfected = msg
				fmt.Printf("\nNode %d will be infecting by pull next.", nextNode)
			} else {
				fmt.Println("Error: Node could not send.")
				break
			}
		default:
			break
		}
	}
}

func pullReceive(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if node.isInfected {
		for len(*(*node).channel) < nodeCount {
			*(*node).channel <- node.isInfected
		}
	}
}


func pull(wg *sync.WaitGroup) {
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		pullReceive(wg, &nodes[i])
	}
	time.Sleep(50 * time.Millisecond)
	wg.Wait()

	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		pullSend(wg, &nodes[i])
	}
	wg.Wait()
}

func runPull(wg *sync.WaitGroup) {
	for i:= 1; i<=turns; i++{
		nodes = make([]Node, nodeCount)
		channels := make([]chan bool, nodeCount)

		for i := 0; i < nodeCount; i++ {
			channels[i] = make(chan bool, nodeCount)
			nodes[i] = Node{false, &(channels[i])}
		}

		nodes[0].isInfected = true
		roundCount = 0
		for true {
			roundCount++
			fmt.Printf("\n\nStarting round %d in turn %d!", roundCount, i)
			pull(wg)
			complete, _ := isDone()
			if complete {
				break
			}
		}
		fmt.Printf("\nIt took %d rounds to infect %d nodes.\n", roundCount, nodeCount)
	}
	fmt.Printf("On average, over the %d turns, it took %f rounds to infect %d nodes.", turns, totalRounds/(float32)(turns), nodeCount)
}

func pushPull(wg *sync.WaitGroup) {
	complete, infectedCount := isDone()
	if !complete{
		if infectedCount <= nodeCount/2 {
			fmt.Print("\nPushing...")
			push(wg)
		} else {
			fmt.Print("\nPulling...")
			pull(wg)
		}
	} else {
		fmt.Println( "Already complete!")
	}
}

func runPushPull(wg *sync.WaitGroup) {
	for i:= 1; i<=turns; i++{
		nodes = make([]Node, nodeCount)
		channels := make([]chan bool, nodeCount)

		for i := 0; i < nodeCount; i++ {
			channels[i] = make(chan bool, nodeCount)
			nodes[i] = Node{false, &(channels[i])}
		}

		nodes[0].isInfected = true
		roundCount = 0
		for true {
			roundCount++
			fmt.Printf("\n\nStarting round %d in turn %d!", roundCount, i)
			pushPull(wg)
			complete, _ := isDone()
			if complete {
				break
			}
		}
		fmt.Printf("\nIt took %d rounds to infect %d nodes.\n", roundCount, nodeCount)
	}
}

func isDone() (bool, int) {
	done := true
	tempCount := 0
	fmt.Print("\nCurrent state of nodes: ")
	for i := 0; i < nodeCount; i++ {
		fmt.Printf("\nNode %d: %t", i, nodes[i].isInfected)
		if !nodes[i].isInfected {
			done = false
		} else {
			tempCount++
		}
	}
	return done, tempCount
}

func main() {
	wg := &sync.WaitGroup{}
	rand.Seed(time.Now().UnixNano()) //Make sure that the seed is randomized.

	fmt.Println("How many nodes would you like to test? Please only enter numbers!")
	fmt.Scanf("%d", &nodeCount)
	fmt.Println("How many times would you like to run the test? Please only enter numbers!")
	fmt.Scanf("%d", &turns)
	fmt.Println("Which type of gossip? Please only enter numbers! \n1: Push   2: Pull   3: Push/Pull")
	fmt.Scanf("%d", &gossipType)

	if gossipType == 1 {
		runPush(wg)
	} else if gossipType == 2 {
		runPull(wg)
	} else if gossipType == 3 {
		runPushPull(wg)
	} else {
		fmt.Println("Error with input! Please try again.")
	}
}
