package MySort

func partition(data Interface, lo, hi int) int {
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

func QuickSort(data Interface, lo, hi int) {
	if lo > hi {
		return
	}
	p := partition(data, lo, hi)
	QuickSort(data, lo, p-1)
	QuickSort(data, p+1, hi)
}
