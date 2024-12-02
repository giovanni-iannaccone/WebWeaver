package internals

func IsProhibited(prohibitedList []string, directory []byte) bool {
	for _, dir := range prohibitedList {
		if string(directory)[:len(dir)] == dir {
			return true
		}
	}

	return false
}