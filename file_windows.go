package txfile

// computePlatformMmapSize computes the maximum amount of bytes to be mmaped,
// depending on the actual file size and the configured maximum file size.
// On Windows, the size returned MUST NOT exceed the actual file size.
func computePlatformMmapSize(fileSize, maxSize, pageSize uint) (uint, error) {
	if maxSize == 0 {
		return fileSize, nil
	}

	sz, err := computeMmapSize(fileSize, maxSize, pageSize)
	if fileSize < sz {
		sz = fileSize
	}
	return sz, err
}
