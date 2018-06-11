// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package txfile

// computePlatformMmapSize computes the maximum amount of bytes to be mmaped,
// depending on the actual file size and the configured maximum file size.
func computePlatformMmapSize(fileSize, maxSize, pageSize uint) (uint, error) {
	return computeMmapSize(fileSize, maxSize, pageSize)
}
