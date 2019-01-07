package requestlimiter

import (
	"log"
	"time"
)

// JobLimiter mechanism for the feature of the request rate limiter
type JobLimiter struct {
	pendingJobQueue    chan *Job
	processingJobQueue chan *Job
	deliverPause       chan bool
	currentJobCount    int
	ticker             *time.Ticker
	limitQPS           int
}

// Init initial the job queue
func (jl *JobLimiter) Init(limitQPS int) {
	jl.ticker = time.NewTicker(2 * time.Second)
	jl.pendingJobQueue = make(chan *Job, 100)
	jl.processingJobQueue = make(chan *Job, 100)
	jl.deliverPause = make(chan bool)
	jl.currentJobCount = 0
	jl.limitQPS = limitQPS

	for i := 1; i < 100; i++ {
		go jl.processJobWorker(jl.processingJobQueue)
	}

	go jl.jobDeliver(jl.pendingJobQueue)

	go func() {
		for t := range jl.ticker.C {
			select {
			case jl.deliverPause <- true:
				log.Println("resume pause", t)
			default:
			}

			if jl.currentJobCount != 0 {
				jl.currentJobCount = 0
				log.Println("clear count", t)
			}
		}
	}()

}

// EnqueueJob ...
func (jl *JobLimiter) EnqueueJob(job *Job) {
	jl.pendingJobQueue <- job
	log.Println(job.getHost(), "enqueued.")
}

func (jl *JobLimiter) processJobWorker(jobs chan *Job) {
	for job := range jobs {
		time.Sleep(3 * time.Second)
		job.WriteResult(job.getHost() + " job completed.")
	}
}

func (jl *JobLimiter) jobDeliver(pendingJobs chan *Job) {
	for job := range pendingJobs {
		if jl.currentJobCount > jl.limitQPS {
			log.Println("jobDeliver pause")
			<-jl.deliverPause
			log.Println("jobDeliver resume")
		}

		jl.processingJobQueue <- job
		jl.currentJobCount++
	}
}