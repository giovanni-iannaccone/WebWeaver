package internals

func IsProhibited(prohibitedList []string, directory []byte) bool {
	for _, dir := range prohibitedList {
		if dir == string(directory) {
			return true
		}
	}

	return false
}