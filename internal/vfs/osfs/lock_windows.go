package osfs

import (
	"fmt"

	flock "github.com/theckman/go-flock"
)

const (
	lockExt = ".lock"
)

type lockState struct {
	*flock.Flock
}

func (f *File) Lock(exclusive, blocking bool) error {
	if f.state.lock.Flock != nil {
		return fmt.Errorf("file %v is already locked", f.Name())
	}

	var ok bool
	var err error
	lock := flock.NewFlock(f.Name() + lockExt)
	if blocking {
		err = lock.Lock()
		ok = err == nil
	} else {
		ok, err = lock.TryLock()
	}

	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("file %v can not be locked right now", f.Name())
	}

	f.state.lock.Flock = lock
	return nil
}

func (f *File) Unlock() error {
	if f.state.lock.Flock == nil {
		return fmt.Errorf("file %v is not locked", f.Name())
	}

	err := f.state.lock.Unlock()
	if err == nil {
		f.state.lock.Flock = nil
	}
	return err
}
