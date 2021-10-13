package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var totalRounds float32 //used to calculate average
var nodeCount int //the number of nodes, as input by the user.
var turns int //number of turns, ie iterations to run the simulation
var roundCount int //number of rounds at the current time
var nodes []Node //all nodes
var gossipType int //type of gossip to simulate


type Node struct {
	isInfected bool
	channel  *chan bool
}

//This function finds a random node to infect, if the node is already infected.
func pushSend(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()

	if node.isInfected == true {
		nextNode := rand.Intn(nodeCount)
		fmt.Printf("\nNode %d will be infected by push next.", nextNode)
		*nodes[nextNode].channel <- node.isInfected
	}
}

//If it receives a node receives an infection message, it will be infected
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

//runs both pushSend and pushReceive for all nodes.
func push(wg *sync.WaitGroup) {
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		pushSend(wg, &nodes[i])
	}
	//nodes go to sleep for 50 miliseconds. This is required or it will cause a deadlock.
	time.Sleep(50 * time.Millisecond) //I've tried using lower sleep timers, but 50 seems to consistently work. It slows down the process a bit.
	wg.Wait()

	for j := 0; j < nodeCount; j++ {
		wg.Add(1)
		pushReceive(wg, &nodes[j])
	}
	wg.Wait()
}

//This function is called by main and runs a simulation of push gossip as many times as the user input for turns.
func runPush(wg *sync.WaitGroup) {
	totalRounds = 0
	for i:= 1; i<=turns; i++{
		fmt.Printf("\nTurn number %d!.\n", i)
		// making a new set of nodes and channels.
		nodes = make([]Node, nodeCount)
		channels := make([]chan bool, nodeCount)

		//initializing the channels and nodes.
		for i := 0; i < nodeCount; i++ {
			channels[i] = make(chan bool, nodeCount)
			nodes[i] = Node{false, &(channels[i])}
		}

		nodes[0].isInfected = true //set node 0 to infected to start things off.
		roundCount = 0 //resetting roundcount every turn
		for true {
			roundCount++
			fmt.Printf("\n\nStarting round %d in turn %d!", roundCount, i)
			push(wg)
			//checks if infection is complete yet.
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

// very similar code to pushReceive. see notes from pushReceive to understand.
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

//This function fills the channel with gossip if the node is infected.
func pullReceive(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()

	if node.isInfected {
		for len(*(*node).channel) < nodeCount {
			*(*node).channel <- node.isInfected
		}
	}
}

//almost exactly the same as push function. Read comments from above to understand.
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

//almost exactly the same as runPush function. Read above to understand.
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


/*This function describes the push/pull switch type of gossip, in which when less than half of the nodes are infected, 
	it will push, and after half of the nodes are infected it will pull.
	This function is why i track infectedCount in isDone.*/
func pushPull(wg *sync.WaitGroup) {
	complete, infectedCount := isDone()
	if !complete{
		//at about halfway point it will swap
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

//This function's code is almost exactly the same as runPush and runPull. See above comments on runPush for explanation.
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

//This function checks whether or not all of the nodes are infected yet.
func isDone() (bool, int) {
	done := true
	tempCount := 0 //tempCount is used to get infectedCount, which is used in pushPull.
	fmt.Print("\nCurrent state of nodes: ")
	for i := 0; i < nodeCount; i++ {
		fmt.Printf("\nNode %d: %t", i, nodes[i].isInfected)
		if !nodes[i].isInfected { // if any node is not infected, it will return false.
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

	// gets user input, exactly as described.
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
