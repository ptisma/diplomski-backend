package backgroundContainer

import (
	"context"
	"fmt"
	"sync"
)

type Scheduler struct {
	ctx               context.Context
	cancel            context.CancelFunc
	waitGroup         *sync.WaitGroup
	backgroundMethods map[int]backgroundWork
	backgroundWorker  BackgroundWorker
}

func (s *Scheduler) ScheduleBackgroundWorks() {

	s.waitGroup.Add(len(s.backgroundMethods))
	//fmt.Println("length", len(s.backgroundMethods))

	for i := 0; i < len(s.backgroundMethods); i++ {
		//fmt.Println("launching i:", i)
		i := i
		go func() {
			s.backgroundMethods[i](s.ctx, s.waitGroup, s.backgroundWorker)

		}()
	}
	fmt.Println("launched all background goroutines")
}

func (s *Scheduler) Exit() {
	s.waitGroup.Wait()
	fmt.Println("exiting scheduler")
}

func NewScheduler(ctx context.Context, worker BackgroundWorker) *Scheduler {
	wg := sync.WaitGroup{}
	return &Scheduler{
		ctx:               ctx,
		backgroundWorker:  worker,
		backgroundMethods: GetBackgroundWorks(),
		waitGroup:         &wg,
	}
}
