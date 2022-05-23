package deadlock

// type to save info about caller
type callerInfo struct {
	file string
	line int
}

// create a new caller info
func newInfo(file string, line int) callerInfo {
	return callerInfo{
		file: file,
		line: line,
	}
}
