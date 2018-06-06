package osfs

type osFileState struct {
	mmap mmapState
	lock lockState
}
