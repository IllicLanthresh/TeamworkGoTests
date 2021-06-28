package helperTypes

type StringSlice []string

//IndexOf tries to find element `str` in the StringSlice and returns the index, returns -1 if the element is not found
func (slice StringSlice) IndexOf(str string) int {
	for p, v := range slice {
		if v == str {
			return p
		}
	}
	return -1
}

//Contains tries to find element `str` in the StringSlice and returns whether it could find it or not
func (slice StringSlice) Contains(str string) bool {
	for _, iStr := range slice {
		if iStr == str {
			return true
		}
	}
	return false
}
