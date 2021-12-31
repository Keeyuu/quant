package arrayutil

func ArrayIndexOfStr(arr []string, target string) int {
	index := -1
	for k, v := range arr {
		if v == target {
			index = k
			break
		}
	}
	return index
}
