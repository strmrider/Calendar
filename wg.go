package calendar

import "sync"

const (
	wgAdd  = 0
	wgDone = 1
)

type SynWG struct {
	wg *sync.WaitGroup
}

func (syncwg *SynWG) set(wg *sync.WaitGroup) {
	syncwg.wg = wg
}

func (syncwg *SynWG) operate(action int) {
	if syncwg.wg != nil {
		if action == wgAdd {
			syncwg.wg.Add(1)
		} else if action == wgDone {
			syncwg.wg.Done()
		}
	}
}

var syncwg SynWG

func SetWG(wg *sync.WaitGroup) {
	syncwg.wg = wg
}