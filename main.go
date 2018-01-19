package main

import (
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"fmt"
	"strconv"
	"sort"
	"time"
	"DS2/MySort"
)

type node struct {
	num        int
	nextEdge   *edgeLinkedList
	lastEdge   *edgeLinkedList
	edgesCount int
	visited    int
	//next       *node
}

type edgeLinkedList struct {
	num          int
	previousEdge *edgeLinkedList
	nextEdge     *edgeLinkedList
}

type edgeStruct struct {
	point float64
	nodes [2]int
}

var NODES map[int]*node
var EDGES []*edgeStruct
var contin bool //continue
var end bool
var counter = 0

const file = "test1.txt"

func main() {
	contin = true
	NODES = make(map[int]*node)
	initNode()
	sortNodes()
	pointEdges()
	runAlgorithm()
	fmt.Println("end")
}

var counterThread = 0

func runAlgorithm() {
	end = true
	len1 := len(EDGES)
	fmt.Println(len1)
	var i, j = 0, 0
	go showCounter()
	MySort.Sort(byPoint(EDGES))
	for contin {
		counter++
		counterThread = 0
		contin = false
		setFalseToNodes()
		i, j = removeFirstEdge()
		//ch := make(chan int, 1)
		dfs(i, j, 0, 1)
		//dfs(i, j, 0, 2)
		//<- ch
		//close(ch)
	}
	contin = false
	setFalseToNodes()
	fmt.Println("removed edgs : ", len1-len(EDGES))
	dfsPrint(i)
}

//before every dfs
func setFalseToNodes() {
	for _, N := range NODES {
		N.visited = 0
	}
}

//for Test
func showEdgsOfNode(i int) {
	iterateEdge := NODES[i].nextEdge
	fmt.Println(" EDGES  ")
	for iterateEdge != nil {
		fmt.Print(iterateEdge.num, " ")
		iterateEdge = iterateEdge.nextEdge
	}
	fmt.Println()
}

func removeEdgeFromNodes(nodeNum, edgeNum int) {
	iterateEdge := NODES[nodeNum].nextEdge
	var previousEdge *edgeLinkedList
	previousEdge = nil
	for iterateEdge != nil {
		if iterateEdge.num == edgeNum {
			if previousEdge != nil {
				previousEdge.nextEdge = iterateEdge.nextEdge
			} else {
				iterateEdge.nextEdge.previousEdge = nil
			}
			iterateEdge.previousEdge = previousEdge
			iterateEdge = nil
			NODES[nodeNum].edgesCount--
			return
		} else {
			previousEdge = iterateEdge
			iterateEdge = iterateEdge.nextEdge
		}
	}
}

//Remove edge with lower Cij from EDGES(slice)
func removeFirstEdge() (int, int) {
	firstEdge := 0
	i, j := EDGES[firstEdge].nodes[0], EDGES[firstEdge].nodes[1]
	removeEdgeFromNodes(i, j)
	removeEdgeFromNodes(j, i)
	EDGES = EDGES[1:]
	return i, j
}

//monitoring
func showCounter() {
	for end {
		fmt.Println("-------------------------")
		fmt.Println("Remove Edges : ", counter)
		fmt.Println("Thread : ", counterThread)
		time.Sleep(1000 * time.Millisecond)
	}
}

//execute end of runAlgorithm and show Result
func dfsPrint(i int) {
	n := NODES[i]
	n.visited = 1
	iterateEdge := n.nextEdge
	for iterateEdge != nil {
		if NODES[iterateEdge.num].visited == 0 {
			fmt.Println(iterateEdge.num)
			dfsPrint(iterateEdge.num)
		}
		iterateEdge = iterateEdge.nextEdge
	}
}

//remove edge from nodes
func dfs(i, j, depth, direction int) {
	n := NODES[i]
	if !contin && n.visited == 0 {
		n.visited = direction
		if i == j {
			contin = true
			//ch <- 0
			return
		}
		iterateEdge := n.nextEdge
		for iterateEdge != nil && !contin {
			if NODES[iterateEdge.num].visited == 0 {
				//if depth < 3 {
				//	go dfs(iterateEdge.num, j, depth+1, direction, ch)
				//	counterThread++
				//} else {
					dfs(iterateEdge.num, j, depth+1, direction)
				//}
				//} else if NODES[iterateEdge.num].visited != 0 && NODES[iterateEdge.num].visited != direction {
				//	contin = true
				//	return
				//}
				//time.Sleep(1 * time.Millisecond)
			}
			iterateEdge = iterateEdge.nextEdge
		}
	}
	//else {
	//	contin = true
	//	return
	//}
}

//calc Cij for each [i,j] edge
func point(i, j int) float64 {
	z := compareTwoArray(i, j)
	z++
	min := findMin(getEdgeNum(i), getEdgeNum(j))
	min--
	if min != 0 {
		return float64(z) / float64(min)
	}
	return -1
}

func getEdgeNum(index int) int {
	return NODES[index].edgesCount
}

//set point(Cji) for all edges
func pointEdges() {
	f, _ := os.Open(file)
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		e1, _ := strconv.Atoi(record[0])
		e2, _ := strconv.Atoi(record[1])
		if e1 < e2 {
			point := point(e1, e2)
			nodes := [2]int{e1, e2}
			EDGES = append(EDGES, &edgeStruct{point, nodes})
		}
	}
}

// create instant of a edge
func NewEdge(edgeNum int, previous *edgeLinkedList) *edgeLinkedList {
	return &edgeLinkedList{edgeNum, previous, nil}
}

//compare two linked list of two node end return duplicate edges // for calc Cij
func compareTwoArray(i, j int) int {
	edgesOfI := NODES[i].nextEdge
	edgesOfJ := NODES[j].nextEdge
	result := 0
	for edgesOfI != nil && edgesOfJ != nil {
		if edgesOfI.num == edgesOfJ.num {
			edgesOfI = edgesOfI.nextEdge
			edgesOfJ = edgesOfJ.nextEdge
			result++
		} else if edgesOfI.num > edgesOfJ.num {
			edgesOfJ = edgesOfJ.nextEdge
		} else {
			edgesOfI = edgesOfI.nextEdge
		}
	}
	return result
}

func findMin(a, b int) int {
	if a > b {
		return b
	}
	return a
}

//sort nodes after set Cij for all nodes
func sortNodes() {
	for index, i := range NODES {
		i.nextEdge, i.lastEdge = setSortedEdges(index)
	}
}

//sort edges of each nodes (NODES map)
func setSortedEdges(i int) (*edgeLinkedList, *edgeLinkedList) {
	edgesArray := getEdgesArray(i)
	sort.Sort(byNumber(edgesArray))
	lenArray := len(edgesArray)
	iterateNum := 0
	firstEdge := edgesArray[iterateNum]
	iterateEdges := firstEdge
	var previousEdge *edgeLinkedList = nil
	for iterateNum < lenArray {
		iterateEdges.nextEdge = edgesArray[iterateNum]
		iterateEdges.previousEdge = previousEdge
		previousEdge = iterateEdges
		iterateEdges = iterateEdges.nextEdge
		iterateNum++
	}
	return firstEdge, iterateEdges
}

// get array of edges of an node
func getEdgesArray(index int) []*edgeLinkedList {
	var edges []*edgeLinkedList
	iterateEdge := NODES[index].nextEdge
	for iterateEdge != nil {
		edges = append(edges, iterateEdge)
		iterateEdge = iterateEdge.nextEdge
	}
	return edges
}

// read file and save nodes in NODES(map)
func initNode() {
	f, _ := os.Open(file)
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		index, _ := strconv.Atoi(record[0])
		edgeNum, _ := strconv.Atoi(record[1])
		if NODES[index] == nil {
			newEdge := NewEdge(edgeNum, nil)
			NODES[index] = &node{index, newEdge, newEdge, 1, 0}
		} else {
			iterateEdge := NODES[index].nextEdge
			last := iterateEdge
			for iterateEdge != nil {
				last = iterateEdge
				iterateEdge = iterateEdge.nextEdge
			}
			NODES[index].edgesCount++
			last.nextEdge = NewEdge(edgeNum, last)
			NODES[index].lastEdge = last.nextEdge
			//TODO sort and insert
		}
	}
}


////reCji
//func repoint() {
//	for _, e := range EDGES {
//		e.point = point(e.nodes[0], e.nodes[1])
//	}
//}

// for edge sort
type byNumber []*edgeLinkedList

func (a byNumber) Len() int      { return len(a) }
func (a byNumber) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byNumber) Less(i, j int) bool {
	if a[i].num < a[j].num {
		return true
	} else {
		return false
	}
}

// for edge point sort
type byPoint []*edgeStruct

func (a byPoint) Len() int      { return len(a) }
func (a byPoint) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byPoint) Less(i, j int) bool {
	if a[i].point < a[j].point {
		return true
	} else {
		return false
	}
}
