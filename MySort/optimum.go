package MySort

func part(data Interface, lo, hi int) int {
	p := hi
	for j := lo; j < hi; j++ {
		if data.Less(j, p) {
			data.Swap(j, lo)
			lo++
		}
	}
	data.Swap(lo, hi)
	return lo
}

func Optimum(data Interface, lo, hi, N, mode int) {
	if lo > hi {
		return
	}
	p := part(data, lo, hi)
	if p-1-lo > N {
		Optimum(data, lo, p-1, N, mode)
	} else if mode == 0 {
		InsertionSort(data, lo, p-1)
	} else {
		BubbleSort(data, lo, p-1)
	}

	if p+1-hi > N {
		Optimum(data, p+1, hi, N, mode)
	} else if mode == 0 {
		InsertionSort(data, p+1, hi)
	} else {
		BubbleSort(data, p+1, hi)
	}
}
