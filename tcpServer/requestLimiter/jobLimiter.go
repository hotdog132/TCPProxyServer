package requestlimiter

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// JobLimiter mechanism for the feature of the request rate limiter
type JobLimiter struct {
	pendingJobQueue    chan *Job
	processingJobQueue chan *Job
	deliverPause       chan bool
	currentJobCount    int
	processedJobCount  int
	ticker             *time.Ticker
	limitQPS           int
	externalAPI        string
}

// Init initial the job queue
func (jl *JobLimiter) Init(limitQPS int) {
	jl.ticker = time.NewTicker(2 * time.Second)
	jl.pendingJobQueue = make(chan *Job, 100)
	jl.processingJobQueue = make(chan *Job, 100)
	jl.deliverPause = make(chan bool)
	jl.currentJobCount = 0
	jl.limitQPS = limitQPS

	// Initial workers to do jobs from processingJobQueue
	for i := 1; i < 100; i++ {
		go jl.processJobWorker(jl.processingJobQueue)
	}

	go jl.jobDeliver(jl.pendingJobQueue)

	go func() {
		for t := range jl.ticker.C {
			select {
			case jl.deliverPause <- true:
				log.Println("Resume pause of limited QPS", t)
			default:
			}

			// Clear currentJobCount per second
			if jl.currentJobCount != 0 {
				jl.currentJobCount = 0
			}
		}
	}()

}

// GetRequestRatePerSec Get processing job count
func (jl *JobLimiter) GetRequestRatePerSec() int {
	return jl.currentJobCount
}

// GetRemainingJobCount Get remaining jobs from pendingJobQueue
func (jl *JobLimiter) GetRemainingJobCount() int {
	return len(jl.pendingJobQueue)
}

// GetProcessedJobCount get proccessed job count
func (jl *JobLimiter) GetProcessedJobCount() int {
	return jl.processedJobCount
}

// SetExternalAPI ...
func (jl *JobLimiter) SetExternalAPI(externalAPI string) {
	jl.externalAPI = externalAPI
}

// EnqueueJob ...
func (jl *JobLimiter) EnqueueJob(job *Job) {
	jl.pendingJobQueue <- job
	// log.Println(job.getHost(), "enqueued a job")
}

func (jl *JobLimiter) processJobWorker(jobs chan *Job) {
	for job := range jobs {
		// External api call
		queryURL := jl.externalAPI + "?q=" + job.GetQuery()
		resp, err := http.Get(queryURL)
		if err != nil {
			// Handle the case when external api shut down
			job.WriteResult(job.getHost() + " external api is not available.")
			jl.processedJobCount++
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		job.WriteResult(job.getHost() + " query result: " + string(body))

		jl.processedJobCount++
	}
}

func (jl *JobLimiter) jobDeliver(pendingJobs chan *Job) {
	for job := range pendingJobs {
		if jl.currentJobCount >= jl.limitQPS {
			log.Println("Pause due to limited QPS")
			<-jl.deliverPause
			log.Println("Resume pause and clear currentJobCount")
			jl.currentJobCount = 0
		}

		jl.processingJobQueue <- job
		jl.currentJobCount++
	}
}
