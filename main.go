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
	"runtime"
	"strings"
)

type Node struct {
	num        int
	nextEdge   *EdgeLinkedList
	lastEdge   *EdgeLinkedList
	edgesCount int
	visited    int
}

type EdgeLinkedList struct {
	num          int
	previousEdge *EdgeLinkedList
	nextEdge     *EdgeLinkedList
}
type EdgeStruct struct {
	point        float64
	nodes        [2]int
	foreignIndex int
}
type Edge struct {
	index int
	nodes [2]int
}

//linked list
var NODES map[int]*Node
var EDGES []*EdgeStruct
//sparse
var startNode map[int]float64
var SparseMatrix []*Edge

var contin bool //continue
var end bool
var counter = 0

const file = "test1.txt"

func main() {
	for {
		command := getCommand()
		if len(command) == 3 && command[0] == "run" && command[2] != "optimum" {
			contin = true
			NODES = make(map[int]*Node)
			initNode()
			sortNodes()
			pointEdges()
			runAlgorithm(command[2], 0)
			fmt.Println("end")
		} else if len(command) == 5 && command[0] == "run" && command[2] == "optimum" {
			contin = true
			NODES = make(map[int]*Node)
			initNode()
			sortNodes()
			pointEdges()
			N, _ := strconv.Atoi(command[4])
			if command[3] == "insertion" {
				runAlgorithm("quickInsertion", N)
			} else {
				runAlgorithm("quickBubble", N)
			}
			fmt.Println("end")
		} else {
			fmt.Println("invalid command")
		}
	} //end of while

}

func getCommand() []string {
	fmt.Print("Enter your command: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	commands := strings.Split(text, "\n")
	text = commands[0]
	text = string(text)
	commands = strings.Split(text, " ")
	return commands
}

func runAlgorithm(mode string, N int) {
	end = true
	len1 := len(EDGES)
	fmt.Println("Number of Edges : ", len1)
	var i, j = 0, 0
	go showCounter()
	MySort.Sort(byPoint(EDGES))
	for contin {
		counter++
		contin = false
		i, j = removeFirstEdge()
		repoint(i, j)
		switch mode {
		case "bubble":
			MySort.BubbleSort(byPoint(EDGES), 0, len(EDGES)-1)
		case "insertion":
			MySort.InsertionSort(byPoint(EDGES), 0, len(EDGES)-1)
		case "quickSort", "quick":
			MySort.QuickSort(byPoint(EDGES), 0, len(EDGES)-1)
		case "optimumInsertion":
			MySort.Optimum(byPoint(EDGES), 0, len(EDGES)-1, N, 0)
		case "optimumBubble":
			MySort.Optimum(byPoint(EDGES), 0, len(EDGES)-1, N, 1)
		case "mergeSort", "merge":
			MergeSort(EDGES)
		default:
			MySort.Sort(byPoint(EDGES))
		}
		setFalseToNodes()
		dfs(i, j, 0, 1)
	}
	contin = false
	end = false
	setFalseToNodes()
	fmt.Println("Removed Edges : ", len1-len(EDGES))
	f, err := os.Create("result" + file)
	check(err)
	defer f.Close()
	f.Sync()
	w := bufio.NewWriter(f)
	w.WriteString("runtime " + strconv.Itoa(counterSecond) + "\n")
	dfsPrint(i, "A", w)
	dfsPrint(j, "B", w)
	w.WriteString("alloc " + strconv.Itoa(int(alloc)) + "\n")

}

//remove edge from nodes
func dfs(i, j, depth, direction int) {
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
				dfs(iterateEdge.num, j, depth, direction)
			}
			iterateEdge = iterateEdge.nextEdge
		}
	}
}

//for read file
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
func removeEdgeFromSparseMatrix(i int) {
	SparseMatrix = append(SparseMatrix[:i], SparseMatrix[i+1:]...)
}
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

var counterSecond = 0
//monitoring
func showCounter() {
	for end {
		fmt.Println("-------------------------")
		counterSecond++
		fmt.Println("Second", counterSecond)
		fmt.Println("Removed Edges : ", counter)

		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		alloc = mem.Alloc
		fmt.Println("Memory Used", mem.Alloc)
		//fmt.Println("Memory Used", mem.TotalAlloc)
		//fmt.Println(mem.HeapAlloc)
		//fmt.Println(mem.HeapSys)
		time.Sleep(1000 * time.Millisecond)
	}
}

var alloc uint64

//execute end of runAlgorithm and show Result
func dfsPrint(i int, str string, w *bufio.Writer) {
	n := NODES[i]
	n.visited = 1
	iterateEdge := n.nextEdge
	w.WriteString(strconv.Itoa(i) + ": #" + str + "\n")
	for iterateEdge != nil {
		if NODES[iterateEdge.num].visited == 0 {
			dfsPrint(iterateEdge.num, str, w)
		}
		iterateEdge = iterateEdge.nextEdge
	}
}

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
var counterEdges = 0

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
			counterEdges++
			EDGES = append(EDGES, &EdgeStruct{point, nodes, counterEdges})
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
func repoint(i, j int) {
	for _, e := range EDGES {
		if e.nodes[0] == i || e.nodes[1] == j {
			e.point = getPoint(e.nodes[0], e.nodes[1])
		}
	}
}

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

//Merge sort

func Merge(left, right []*EdgeStruct) []*EdgeStruct {
	result := make([]*EdgeStruct, 0, len(left)+len(right))

	for len(left) > 0 || len(right) > 0 {
		if len(left) == 0 {
			return append(result, right...)
		}
		if len(right) == 0 {
			return append(result, left...)
		}
		if left[0].point <= right[0].point {
			result = append(result, left[0])
			left = left[1:]
		} else {
			result = append(result, right[0])
			right = right[1:]
		}
	}

	return result
}

func MergeSort(arr []*EdgeStruct) []*EdgeStruct {
	if len(arr) <= 1 {
		return arr
	}

	middle := len(arr) / 2

	left := MergeSort(arr[:middle])
	right := MergeSort(arr[middle:])

	return Merge(left, right)
}
