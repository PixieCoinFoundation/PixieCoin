package zk_manager

// import (
// 	"errors"
// 	"github.com/samuel/go-zookeeper/zk"
// 	. "logger"
// 	"math/rand"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )

// const (
// 	INCRE_STEP   = 5
// 	DECRE_STEP   = 1
// 	CHECK_SECOND = 30
// )

// var (
// 	PoolSizeErr = errors.New("pool size not legal")
// )

// type Pool struct {
// 	rawConns  []*zk.Conn
// 	addresses []string

// 	maxSize int

// 	endIndex     int
// 	endIndexLock sync.RWMutex

// 	busyCnt int32
// 	idleCnt int32

// 	// lastGet int64
// }

// func NewZKPool(initSize, maxSize int, addresses []string) (p *Pool, err error) {
// 	if initSize < 1 || maxSize < initSize {
// 		err = PoolSizeErr
// 		return
// 	}

// 	p = &Pool{
// 		rawConns:  make([]*zk.Conn, maxSize),
// 		maxSize:   maxSize,
// 		endIndex:  initSize - 1,
// 		addresses: addresses,
// 		busyCnt:   0,
// 		idleCnt:   0,
// 	}

// 	for i := 0; i < initSize; i++ {
// 		if p.rawConns[i], _, err = zk.Connect(addresses, time.Second); err != nil {
// 			Err(err)
// 			return
// 		}
// 	}

// 	go p.startCheckStatus()

// 	return
// }

// func (p *Pool) startCheckStatus() {
// 	for {
// 		time.Sleep(CHECK_SECOND * time.Second)

// 		p.checkStatus1()
// 	}
// }

// func (p *Pool) checkStatus1() {
// 	busy := atomic.LoadInt32(&p.busyCnt)
// 	idle := atomic.LoadInt32(&p.idleCnt)
// 	Info("zk pool status in", CHECK_SECOND, " seconds busy cnt:", busy, "idle cnt:", idle)
// 	defer atomic.StoreInt32(&p.busyCnt, 0)
// 	defer atomic.StoreInt32(&p.idleCnt, 0)

// 	if busy > idle {
// 		//incre
// 		newEndIndex := p.endIndex + INCRE_STEP
// 		if newEndIndex > p.maxSize-1 {
// 			newEndIndex = p.maxSize - 1
// 		}

// 		if newEndIndex <= p.endIndex {
// 			Err("zk hot!zk pool incre failed.now size:", p.endIndex+1, "max size:", p.maxSize)
// 			return
// 		} else {
// 			var err error
// 			for i := p.endIndex + 1; i <= newEndIndex; i++ {
// 				if p.rawConns[i] == nil {
// 					if p.rawConns[i], _, err = zk.Connect(p.addresses, time.Second); err != nil {
// 						Err(err)
// 						return
// 					}
// 				}
// 			}

// 			//incre success
// 			Err("zk pool size incre to:", newEndIndex+1, "for busy cnt:", busy)
// 			p.endIndexLock.Lock()
// 			p.endIndex = newEndIndex
// 			p.endIndexLock.Unlock()
// 		}
// 	} else if idle >= busy && p.endIndex > 0 && busy == 0 {
// 		//decre
// 		newEndIndex := p.endIndex - DECRE_STEP
// 		if newEndIndex < 0 {
// 			newEndIndex = 0
// 		}

// 		if newEndIndex >= p.endIndex {
// 			Err("decre failed now size:", p.endIndex+1, "new size:", newEndIndex+1)
// 			return
// 		}

// 		oriIndex := p.endIndex
// 		//no new request will come into tail conns
// 		p.endIndexLock.Lock()
// 		p.endIndex = newEndIndex
// 		p.endIndexLock.Unlock()

// 		for i := newEndIndex + 1; i <= oriIndex; i++ {
// 			if p.rawConns[i] != nil {
// 				p.rawConns[i].SafeCloseWithTimeout(time.Second)
// 				p.rawConns[i] = nil
// 				Info("zk conn pool safe close:", i)
// 			}
// 		}
// 		Info("zk pool size decre to:", newEndIndex+1, "for idle:", idle)
// 	}
// }

// func (p *Pool) Close() {
// 	p.endIndexLock.Lock()
// 	p.endIndex = -1
// 	p.endIndexLock.Unlock()
// }

// func (p *Pool) Get() *zk.Conn {
// 	// atomic.StoreInt64(&p.lastGet, time.Now().Unix())
// 	p.endIndexLock.RLock()
// 	if p.endIndex < 0 {
// 		p.endIndexLock.RUnlock()
// 		Err("zk conn pool no conn!!!")
// 		return nil
// 	} else {
// 		c := p.rawConns[rand.Intn(p.endIndex+1)]
// 		p.endIndexLock.RUnlock()

// 		if c.Busy() {
// 			atomic.AddInt32(&p.busyCnt, 1)
// 			// atomic.StoreInt32(&p.idleCnt, 0)
// 		} else {
// 			atomic.AddInt32(&p.idleCnt, 1)
// 			// atomic.StoreInt32(&p.busyCnt, 0)
// 		}

// 		return c
// 	}
// }
