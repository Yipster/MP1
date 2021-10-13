# MP1
MP1's goal is to simulate three different types of gossip using go-channels. The three types of gossip are: push, pull, and push/pull switch.

# Specification

On a high level, MP1 runs a simulation of gossip of either push, pull, and push/pull switch. It will ask for the number of nodes you wish to simulate, and then ask for how many times to run the simulation. Finally, it will ask you which type of gossip you would like. Importantly, the number of times to run the simulation will be very helpful in determining the relation between the number of nodes and the average number of rounds it takes to fully infect.

After the user inputs are taken, it will begin simulating that type of gossip. All three types of simulations will start with node 0 being infected, so if the number of nodes is 1, it will always be complete infected right off the bat. Additionally, the program will show you the process of the gossip happening, with each round of infecting being displayed, which includes which nodes are randomly selected to be infected next. At the end of each round, it will display which nodes are infected or not, with true meaning infected and false meaning uninfected. 

Once it has completed the gossip infection for the assigned number of turns, the program will output the average number of rounds required to infect n amount of nodes. 

# Push

The idea behind push gossip is to have currently infected nodes target a random node to send it's message over. The other nodes are also always looking for the message, so that they are able to change their own state whenever they receive it.

I simulate push gossip by using the functions pushSend, pushReceive, push, and runPush. pushSend has a node find a random node to infect, if itself is already infected. pushReceive is the receiving side of that, where if a node is the target of a pushSend then it will turn it's isInfected to true. push function runs both pushSend and pushReceive for each of the nodes. runPush makes and initializes new nodes, channels, and runs the push function for the assigned number of turns, as well as giving the average after it is all done.

# Pull

Pull gossip works by having uninfected nodes targeting random nodes to see their infected status. If the target is infected, then itself will become infected.

Pull gossip is simulated by using the functions pullSend, pullReceive, pull, and runPull. A lot of it works very similarly to the push counterpart, so see above for more details on that. The only big difference here is for pullReceive, which fills the channel with gossip if the node is infected. This was not needed for push gossip.

# Push/Pull Switch

Push/pull switch gossip is a type of gossip that runs both push and pull. It starts off by running push as described above, but when half of the nodes are infected it will switch to using the pull method. 

Push/pull gossip is simulated by making use of existing push and pull functions. pushPull figures out whether or not we are at the halfway mark, and accordingly uses the correct method required. runPushPull runs the simulation much like runPush or runPull.

# isDone()
The isDone function checks at the end of each round whether or not all the nodes have been infected. Additionally, it also returns the number of nodes that are infected, which is used by the pushPull function in determining when we are at the halfway mark.

# How to run

After cloning the repository, simply cd to the file and type "go run main.go". This will begin the program. Make sure you have Go installed!

![Screen Shot 2021-10-12 at 11 09 46 PM](https://user-images.githubusercontent.com/70530925/137061203-abc91fb8-37cd-404a-899b-43513ab3ff21.png)

When the program starts running, it will ask for three different things. Be sure to only input numbers, or there will be an error.

![Screen Shot 2021-10-12 at 11 18 19 PM](https://user-images.githubusercontent.com/70530925/137062335-fe67730c-8ddf-4b54-b38b-0fbf90ffac3a.png)

After this input, it will begin the gossip simulation.

# Sample output

Below are screenshots of an output of using 5 nodes, 10 turns, and push/pull type gossip (number 3 in input).

![Screen Shot 2021-10-12 at 11 21 07 PM](https://user-images.githubusercontent.com/70530925/137061593-630de636-908b-448b-bb1e-26ec8c9a5a16.png)

![Screen Shot 2021-10-12 at 11 21 51 PM](https://user-images.githubusercontent.com/70530925/137061651-5d072b2d-203a-4136-ae2c-60e795fb45d1.png)

![Screen Shot 2021-10-12 at 11 27 02 PM](https://user-images.githubusercontent.com/70530925/137062307-ea6d984f-99a5-4233-89ef-bd2b0b498243.png)


# Resources Used:

https://gobyexample.com/waitgroups

https://pkg.go.dev/time

https://gobyexample.com/channels

https://stackoverflow.com/questions/19208725/example-for-sync-waitgroup-correct

https://www.geeksforgeeks.org/pointers-in-golang/
