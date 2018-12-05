package zk

import (
	"errors"
	"fmt"
	. "logger"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrDeadlock is returned by Lock when trying to lock twice without unlocking first
	ErrDeadlock = errors.New("zk: trying to acquire a lock twice")
	// ErrNotLocked is returned by Unlock when trying to release a lock that has not first be acquired.
	ErrNotLocked = errors.New("zk: not locked")
	// timeout
	ErrLockTimeout = errors.New("zk: lock time out")
	// watch timeout
	ErrLockWatchTimeout = errors.New("zk: lock watch time out")
)

// Lock is a mutual exclusion lock.
type Lock struct {
	c        *Conn
	path     string
	acl      []ACL
	lockPath string
	seq      int
}

// NewLock creates a new lock instance using the provided connection, path, and acl.
// The path must be a node that is only used by this lock. A lock instances starts
// unlocked until Lock() is called.
func NewLock(c *Conn, path string, acl []ACL) *Lock {
	return &Lock{
		c:    c,
		path: path,
		acl:  acl,
	}
}

func parseSeq(path string) (int, error) {
	parts := strings.Split(path, "-")
	return strconv.Atoi(parts[len(parts)-1])
}

func UnlockPath(c *Conn, path string) error {
	if err := c.Delete(path, -1); err != nil {
		return err
	}
	return nil
}

func LockPathWithTimeout(c *Conn, path string, acl []ACL, timeoutSecond int) (pathc string, err error) {
	prefix := fmt.Sprintf("%s/lock-", path)

	for i := 0; i < 3; i++ {
		pathc, err = c.CreateProtectedEphemeralSequential(prefix, []byte{}, acl)
		if err == ErrNoNode {
			// Create parent node.
			parts := strings.Split(path, "/")
			pth := ""
			for _, p := range parts[1:] {
				var exists bool
				pth += "/" + p
				exists, _, err = c.Exists(pth)
				if err != nil {
					Err(err)
					return
				}
				if exists == true {
					continue
				}
				_, err = c.Create(pth, []byte{}, 0, acl)
				if err != nil && err != ErrNodeExists {
					Err(err)
					return
				}
			}
		} else if err == nil {
			break
		} else {
			Err(err)
			return
		}
	}
	if err != nil {
		return
	}

	var seq int
	seq, err = parseSeq(pathc)
	if err != nil {
		Err(err)
		return
	}

	wc := 0
	var children []string
	for {
		if wc >= 3 {
			err = ErrLockTimeout
			return
		}

		children, _, err = c.Children(path)
		if err != nil {
			Err(err)
			return
		}

		lowestSeq := seq
		prevSeq := -1
		prevSeqPath := ""
		var s int
		for _, p := range children {
			s, err = parseSeq(p)
			if err != nil {
				Err(err)
				return
			}
			if s < lowestSeq {
				lowestSeq = s
			}
			if s < seq && s > prevSeq {
				prevSeq = s
				prevSeqPath = p
			}
		}

		if seq == lowestSeq {
			// Acquired the lock
			break
		}

		// Wait on the node next in line for the lock
		var ch <-chan Event
		_, _, ch, err = c.GetW(path + "/" + prevSeqPath)
		if err != nil && err != ErrNoNode {
			return
		} else if err != nil && err == ErrNoNode {
			// try again
			err = nil
			continue
		}

		to := timeoutSecond
		if to <= 0 {
			to = 10
		}

		select {
		case ev := <-ch:
			if ev.Err != nil {
				err = ev.Err
				Err(err)
				return
			}
		case <-time.After(time.Duration(to) * time.Second):
			err = ErrLockWatchTimeout
			return
		}

		wc++
	}
	return
}

// Lock attempts to acquire the lock. It will wait to return until the lock
// is acquired or an error occurs. If this instance already has the lock
// then ErrDeadlock is returned.
func (l *Lock) Lock() error {
	if l.lockPath != "" {
		return ErrDeadlock
	}

	prefix := fmt.Sprintf("%s/lock-", l.path)

	path := ""
	var err error
	for i := 0; i < 3; i++ {
		path, err = l.c.CreateProtectedEphemeralSequential(prefix, []byte{}, l.acl)
		if err == ErrNoNode {
			// Create parent node.
			parts := strings.Split(l.path, "/")
			pth := ""
			for _, p := range parts[1:] {
				var exists bool
				pth += "/" + p
				exists, _, err = l.c.Exists(pth)
				if err != nil {
					return err
				}
				if exists == true {
					continue
				}
				_, err = l.c.Create(pth, []byte{}, 0, l.acl)
				if err != nil && err != ErrNodeExists {
					return err
				}
			}
		} else if err == nil {
			break
		} else {
			return err
		}
	}
	if err != nil {
		return err
	}

	seq, err := parseSeq(path)
	if err != nil {
		return err
	}

	for {
		children, _, err := l.c.Children(l.path)
		if err != nil {
			return err
		}

		lowestSeq := seq
		prevSeq := -1
		prevSeqPath := ""
		for _, p := range children {
			s, err := parseSeq(p)
			if err != nil {
				return err
			}
			if s < lowestSeq {
				lowestSeq = s
			}
			if s < seq && s > prevSeq {
				prevSeq = s
				prevSeqPath = p
			}
		}

		if seq == lowestSeq {
			// Acquired the lock
			break
		}

		// Wait on the node next in line for the lock
		_, _, ch, err := l.c.GetW(l.path + "/" + prevSeqPath)
		if err != nil && err != ErrNoNode {
			return err
		} else if err != nil && err == ErrNoNode {
			// try again
			continue
		}

		ev := <-ch
		if ev.Err != nil {
			return ev.Err
		}
	}

	l.seq = seq
	l.lockPath = path
	return nil
}

// Unlock releases an acquired lock. If the lock is not currently acquired by
// this Lock instance than ErrNotLocked is returned.
func (l *Lock) Unlock() error {
	if l.lockPath == "" {
		return ErrNotLocked
	}
	if err := l.c.Delete(l.lockPath, -1); err != nil {
		return err
	}
	l.lockPath = ""
	l.seq = 0
	return nil
}
