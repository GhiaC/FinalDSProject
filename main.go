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

type Node struct {
	num        int
	nextEdge   *EdgeLinkedList
	lastEdge   *EdgeLinkedList
	edgesCount int
	visited    int
	//next       *node
}

type EdgeLinkedList struct {
	num          int
	previousEdge *EdgeLinkedList
	nextEdge     *EdgeLinkedList
}

type EdgeStruct struct {
	point float64
	nodes [2]int
}

var NODES map[int]*Node
var EDGES []*EdgeStruct
var contin bool //continue
var end bool
var counter = 0

const file = "test1.txt"

func main() {
	contin = true
	NODES = make(map[int]*Node)
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
		repoint()
		MySort.Sort(byPoint(EDGES))
		//MySort.BubbleSort(byPoint(EDGES), 0, len(EDGES)-1)
		//MySort.InsertionSort(byPoint(EDGES), 0, len(EDGES)-1)
		//MySort.QuickSort(byPoint(EDGES), 0, len(EDGES)-1)
		//MySort.Optimum(byPoint(EDGES), 0, len(EDGES)-1, 25, 0)
		i, j = removeFirstEdge()
		bfs(i, j, 0, 1)
	}
	contin = false
	setFalseToNodes()

	fmt.Println("removed edgs : ", len1-len(EDGES))

	f, err := os.Create("result" + file)
	check(err)
	defer f.Close()
	f.Sync()
	w := bufio.NewWriter(f)
	bfsPrint(i, "A", w)
	bfsPrint(j, "B", w)
}

//remove edge from nodes
func bfs(i, j, depth, direction int) {
	n := NODES[i]
	if !contin && n.visited == 0 {
		n.visited = direction
		if i == j {
			endBfs()
			return
		}
		iterateEdge := n.nextEdge
		for iterateEdge != nil && !contin {
			if NODES[iterateEdge.num].visited == 0 {
				bfs(iterateEdge.num, j, depth, direction)
			}
			iterateEdge = iterateEdge.nextEdge
		}
	}
}

//func findLowest() (int, int) {
//	i, j := 0, 0
//	min := 2.0
//	for _, e := range EDGES {
//		e.point = getPoint(e.nodes[0], e.nodes[1])
//		p := getPoint(e.nodes[0], e.nodes[1])
//		if p < min {
//			min = p
//			i, j = e.nodes[0], e.nodes[1]
//		}
//	}
//	return i, j
//}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//before every dfs
func setFalseToNodes() {
	for _, N := range NODES {
		N.visited = 0
	}
}

//for Test
//func showEdgsOfNode(i int) {
//	iterateEdge := NODES[i].nextEdge
//	fmt.Println(" EDGES  ")
//	for iterateEdge != nil {
//		fmt.Print(iterateEdge.num, " ")
//		iterateEdge = iterateEdge.nextEdge
//	}
//	fmt.Println()
//}

func removeEdgeFromNodes(nodeNum, edgeNum int) {
	iterateEdge := NODES[nodeNum].nextEdge
	var previousEdge *EdgeLinkedList
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
func bfsPrint(i int, str string, w *bufio.Writer) {
	n := NODES[i]
	n.visited = 1
	iterateEdge := n.nextEdge
	w.WriteString("#" + str + " " + strconv.Itoa(i) + "\n")
	for iterateEdge != nil {
		if NODES[iterateEdge.num].visited == 0 {
			bfsPrint(iterateEdge.num, str, w)
		}
		iterateEdge = iterateEdge.nextEdge
	}
}

// check available channel
//func IsClosed(ch <-chan int) bool {
//	select {
//	case <-ch:
//		return true
//	default:
//	}
//
//	return false
//}

func endBfs() {
	contin = true
}

//calc Cij for each [i,j] edge
func getPoint(i, j int) float64 {
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
			point := getPoint(e1, e2)
			nodes := [2]int{e1, e2}
			EDGES = append(EDGES, &EdgeStruct{point, nodes})
		}
	}
}

// create instant of a edge
func NewEdge(edgeNum int, previous *EdgeLinkedList) *EdgeLinkedList {
	return &EdgeLinkedList{edgeNum, previous, nil}
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
func setSortedEdges(i int) (*EdgeLinkedList, *EdgeLinkedList) {
	edgesArray := getEdgesArray(i)
	sort.Sort(byNumber(edgesArray))
	lenArray := len(edgesArray)
	iterateNum := 0
	firstEdge := edgesArray[iterateNum]
	iterateEdges := firstEdge
	var previousEdge *EdgeLinkedList = nil
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
func getEdgesArray(index int) []*EdgeLinkedList {
	var edges []*EdgeLinkedList
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
			NODES[index] = &Node{index, newEdge, newEdge, 1, 0}
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
		}
	}
}

////reCji
func repoint() {
	for _, e := range EDGES {
		e.point = getPoint(e.nodes[0], e.nodes[1])
	}
}

// Thread
//func point(first, end int) {
//	for first < end && contin {
//		edge := EDGES[first]
//		edge.point = getPoint(edge.nodes[0], edge.nodes[1])
//		first++
//	}
//}

// for edge sort
type byNumber []*EdgeLinkedList

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
type byPoint []*EdgeStruct

func (a byPoint) Len() int      { return len(a) }
func (a byPoint) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byPoint) Less(i, j int) bool {
	if a[i].point < a[j].point {
		return true
	} else {
		return false
	}
}
