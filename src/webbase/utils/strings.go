package utils

func DeleteSliceElement(array []string, s string) []string {
	if len(array) == 0 {
		return array
	}

	index := -1
	for i, a := range array {
		if a == s {
			index = i
		}
	}

	if index == -1 {
		return array
	}

	return append(array[:index], array[index+1:]...)
}

func DeleteSliceSubset(array []string, sub []string) []string {
	for _, s := range sub {
		array = DeleteSliceElement(array, s)
	}
	return array
}
