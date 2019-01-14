package requestlimiter

import "sync"

// JobCounter Count the number of jobs
type JobCounter struct {
	muxJobCount sync.Mutex
	jobCount    int
}

func (jc *JobCounter) addJobCount(count int) {
	jc.muxJobCount.Lock()
	jc.jobCount += count
	jc.muxJobCount.Unlock()
}

func (jc *JobCounter) setJobCount(count int) {
	jc.muxJobCount.Lock()
	jc.jobCount = count
	jc.muxJobCount.Unlock()
}

func (jc *JobCounter) getJobCount() int {
	jc.muxJobCount.Lock()
	count := jc.jobCount
	jc.muxJobCount.Unlock()
	return count
}

func (jc *JobCounter) getJobCountWithoutLock() int {
	return jc.jobCount
}
