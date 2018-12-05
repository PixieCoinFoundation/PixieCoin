package tasksequence

import (
	"errors"
	"fmt"
	"sync"
)

type task_seq struct {
	name             string
	concurrencyLevel int
	channel          chan func()
	wg               sync.WaitGroup

	workingCount     int
	workingCountLock sync.Mutex
}

var seq_map map[string]*task_seq

func init() {
	seq_map = make(map[string]*task_seq)
}

func Init(taskName string, concurrencyLevel int) error {
	if _, ok := seq_map[taskName]; ok {
		return errors.New("save task sequence name exists")
	}

	ts := &task_seq{
		name:             taskName,
		concurrencyLevel: concurrencyLevel,
		channel:          make(chan func(), concurrencyLevel),
	}

	seq_map[taskName] = ts

	go schedule_ts(seq_map[taskName])

	return nil
}

func schedule_ts(ts *task_seq) {
	for {
		select {
		case task := <-ts.channel:
			ts.workingCountLock.Lock()
			ts.workingCount++
			ts.workingCountLock.Unlock()
			ts.wg.Add(1)
			go func() {
				task()
				ts.workingCountLock.Lock()
				ts.workingCount--
				ts.workingCountLock.Unlock()
				ts.wg.Done()
			}()
			if ts.workingCount >= ts.concurrencyLevel {
				ts.wg.Wait()
			} else {
				break
			}
		}

	}
}

// need to init task sequence named "taskName" first
func Add(taskName string, task func()) {
	if ts, ok := seq_map[taskName]; !ok {
		fmt.Println("no task sequence name is", taskName)
		return
	} else {
		ts.channel <- func() {
			task()
		}
		return
	}
}
