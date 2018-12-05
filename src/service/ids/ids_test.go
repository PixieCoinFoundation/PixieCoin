package ids

import (
	"fmt"
	"sync"
	"testing"
	// "time"
)

var (
	helpID     int
	helpIDLock sync.Mutex
	generated  []int
)

func init() {
	helpIDLock.Lock()
	helpID = 0
	helpIDLock.Unlock()

	generated = make([]int, 0)
}

func Test_GenHelpID(t *testing.T) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		fmt.Println("comparing")

		for i := 0; i < len(generated); i++ {
			for j := i + 1; i < len(generated); j++ {
				if j >= len(generated) {
					break
				}
				if generated[i] == generated[j] {
					fmt.Println("same found:", generated[i])
					return
				}
			}
		}

		fmt.Println("no same found")
	}()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 100000; j++ {
				helpIDLock.Lock()
				helpID++
				tmp := helpID
				generated = append(generated, tmp)
				helpIDLock.Unlock()
				fmt.Println("helpID--->", tmp)
				// time.Sleep(1 * time.Second)
			}

			wg.Done()
		}()
	}

	return
}
