package MySort

func BubbleSort(array Interface, a, b int) {
	for i := a; i < b; i++ {
		for j := a; i < b; j++ {
			if array.Less(i+1, i) {
				array.Swap(i, i+1)
			}
		}
	}
}
