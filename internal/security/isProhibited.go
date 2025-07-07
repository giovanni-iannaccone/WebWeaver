package security

func IsProhibited(prohibitedList []string, directory []byte) bool {
	for _, dir := range prohibitedList {
		if len(dir) < len(directory) && 
			string(directory)[:len(dir)] == dir {
			return true
		}
	}

	return false
}