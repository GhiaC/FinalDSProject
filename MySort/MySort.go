package MySort

type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

func insertionSort(data Interface, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

func siftDown(data Interface, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.Less(first+child, first+child+1) {
			child++
		}
		if !data.Less(first+root, first+child) {
			return
		}
		data.Swap(first+root, first+child)
		root = child
	}
}

func heapSort(data Interface, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown(data, lo, i, first)
	}
}

func medianOfThree(data Interface, m1, m0, m2 int) {
	if data.Less(m1, m0) {
		data.Swap(m1, m0)
	}
	if data.Less(m2, m1) {
		data.Swap(m2, m1)
		if data.Less(m1, m0) {
			data.Swap(m1, m0)
		}
	}
}

func swapRange(data Interface, a, b, n int) {
	for i := 0; i < n; i++ {
		data.Swap(a+i, b+i)
	}
}

func doPivot(data Interface, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1)
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThree(data, lo, lo+s, lo+2*s)
		medianOfThree(data, m, m-s, m+s)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree(data, lo, m, hi-1)
	pivot := lo
	a, c := lo+1, hi-1
	for ; a < c && data.Less(a, pivot); a++ {
	}
	b := a
	for {
		for ; b < c && !data.Less(pivot, b); b++ { // data[b] <= pivot
		}
		for ; b < c && data.Less(pivot, c-1); c-- { // data[c-1] > pivot
		}
		if b >= c {
			break
		}
		data.Swap(b, c-1)
		b++
		c--
	}
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		dups := 0
		if !data.Less(pivot, hi-1) { // data[hi-1] = pivot
			data.Swap(c, hi-1)
			c++
			dups++
		}
		if !data.Less(b-1, pivot) { // data[b-1] = pivot
			b--
			dups++
		}
		if !data.Less(m, pivot) { // data[m] = pivot
			data.Swap(m, b-1)
			b--
			dups++
		}

		protect = dups > 1

	}

	if protect {

		for {

			for ; a < b && !data.Less(b-1, pivot); b-- { // data[b] == pivot

			}

			for ; a < b && data.Less(a, pivot); a++ { // data[a] < pivot

			}

			if a >= b {

				break

			}


			data.Swap(a, b-1)

			a++

			b--

		}

	}
	data.Swap(pivot, b-1)

	return b - 1, c

}

func quickSort(data Interface, a, b, maxDepth int) {

	for b-a > 12 { // Use ShellSort for slices <= 12 elements

		if maxDepth == 0 {

			heapSort(data, a, b)

			return

		}

		maxDepth--

		mlo, mhi := doPivot(data, a, b)

		// Avoiding recursion on the larger subproblem guarantees

		// a stack depth of at most lg(b-a).

		if mlo-a < b-mhi {

			quickSort(data, a, mlo, maxDepth)

			a = mhi // i.e., quickSort(data, mhi, b)

		} else {

			quickSort(data, mhi, b, maxDepth)

			b = mlo // i.e., quickSort(data, a, mlo)

		}

	}

	if b-a > 1 {

		// Do ShellSort pass with gap 6

		// It could be written in this simplified form cause b-a <= 12

		for i := a + 6; i < b; i++ {

			if data.Less(i, i-6) {

				data.Swap(i, i-6)

			}

		}

		insertionSort(data, a, b)

	}

}

func Sort(data Interface) {
	n := data.Len()
	quickSort(data, 0, n, maxDepth(n))
}

func maxDepth(n int) int {

	var depth int

	for i := n; i > 0; i >>= 1 {

		depth++

	}

	return depth * 2

}

func Stable(data Interface) {

	stable(data, data.Len())

}

func stable(data Interface, n int) {

	blockSize := 20 // must be > 0

	a, b := 0, blockSize

	for b <= n {

		insertionSort(data, a, b)

		a = b

		b += blockSize

	}

	insertionSort(data, a, n)


	for blockSize < n {

		a, b = 0, 2*blockSize

		for b <= n {

			symMerge(data, a, a+blockSize, b)

			a = b

			b += 2 * blockSize

		}

		if m := a + blockSize; m < n {

			symMerge(data, a, m, n)

		}

		blockSize *= 2

	}

}

func symMerge(data Interface, a, m, b int) {

	if m-a == 1 {

		i := m

		j := b

		for i < j {

			h := int(uint(i+j) >> 1)

			if data.Less(h, a) {

				i = h + 1

			} else {

				j = h

			}

		}

		for k := a; k < i-1; k++ {

			data.Swap(k, k+1)

		}

		return

	}

	if b-m == 1 {

		i := a

		j := m

		for i < j {

			h := int(uint(i+j) >> 1)

			if !data.Less(m, h) {

				i = h + 1

			} else {

				j = h

			}

		}

		// Swap values until data[m] reaches the position i.

		for k := m; k > i; k-- {

			data.Swap(k, k-1)

		}

		return

	}


	mid := int(uint(a+b) >> 1)

	n := mid + m

	var start, r int

	if m > mid {

		start = n - b

		r = mid

	} else {

		start = a

		r = m

	}

	p := n - 1


	for start < r {

		c := int(uint(start+r) >> 1)

		if !data.Less(p-c, c) {

			start = c + 1

		} else {

			r = c

		}

	}


	end := n - start

	if start < m && m < end {

		rotate(data, start, m, end)

	}

	if a < start && start < mid {

		symMerge(data, a, start, mid)

	}

	if mid < end && end < b {

		symMerge(data, mid, end, b)

	}

}

func rotate(data Interface, a, m, b int) {

	i := m - a

	j := b - m


	for i != j {

		if i > j {

			swapRange(data, m-i, m, j)

			i -= j

		} else {

			swapRange(data, m-i, m+j-i, i)

			j -= i

		}

	}

	// i == j

	swapRange(data, m-i, m, i)

}