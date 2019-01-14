package requestlimiter

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// JobLimiter mechanism for the feature of the request rate limiter
type JobLimiter struct {
	cPendingJobQueue    chan *Job
	cProcessingJobQueue chan *Job
	currentJobCounter   *JobCounter
	cDeliverPause       chan bool
	cProcessedJobCount  chan int
	ticker              *time.Ticker
	externalAPI         string
	processedJobCount   int
	limitQPS            int
}

// Init initial the job queue
func (jl *JobLimiter) Init(limitQPS int) {
	jl.ticker = time.NewTicker(2 * time.Second)
	jl.cPendingJobQueue = make(chan *Job, 100)
	jl.cProcessingJobQueue = make(chan *Job, 100)
	jl.cDeliverPause = make(chan bool)
	jl.cProcessedJobCount = make(chan int, 100)
	jl.limitQPS = limitQPS
	jl.currentJobCounter = &JobCounter{}

	// Initial workers to do jobs from cProcessingJobQueue
	for i := 1; i < 100; i++ {
		go jl.processJobWorker(jl.cProcessingJobQueue)
	}

	jl.processedJobCountWorker()

	go jl.jobDeliver(jl.cPendingJobQueue)

	go func() {
		defer jl.ticker.Stop()

		for t := range jl.ticker.C {
			select {
			case jl.cDeliverPause <- true:
				log.Println("Resume pause of limited QPS", t)
			default:
			}

			// Clear currentJobCount per second
			// if jl.currentJobCounter.getJobCount() != 0 {
			jl.currentJobCounter.setJobCount(0)
			// }
		}
	}()

}

// GetRequestRatePerSec Get processing job count
func (jl *JobLimiter) GetRequestRatePerSec() int {
	return jl.currentJobCounter.getJobCountWithoutLock()
}

// GetRemainingJobCount Get remaining jobs from cPendingJobQueue
func (jl *JobLimiter) GetRemainingJobCount() int {
	return len(jl.cPendingJobQueue)
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
	jl.cPendingJobQueue <- job
	// log.Println(job.getHost(), "enqueued a job")
}

// processedJobCountWorker Set a worker for handling processed job count to avoid race condition
func (jl *JobLimiter) processedJobCountWorker() {
	go func() {
		for count := range jl.cProcessedJobCount {
			jl.processedJobCount += count
		}
	}()
}

func (jl *JobLimiter) processJobWorker(jobs chan *Job) {
	for job := range jobs {
		// External api call
		queryURL := jl.externalAPI + "?q=" + job.GetQuery()
		resp, err := http.Get(queryURL)
		if err != nil {
			// Handle the case when external api shut down
			job.WriteResult(job.getHost() + " external api is not available.")
			jl.cProcessedJobCount <- 1
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		job.WriteResult(job.getHost() + " query result: " + string(body))

		jl.cProcessedJobCount <- 1
	}
}

func (jl *JobLimiter) jobDeliver(pendingJobs chan *Job) {
	for job := range pendingJobs {
		if jl.currentJobCounter.getJobCount() >= jl.limitQPS {
			log.Println("Pause due to limited QPS")
			<-jl.cDeliverPause
			log.Println("Resume pause and clear currentJobCount")
			jl.currentJobCounter.setJobCount(0)
		}

		jl.cProcessingJobQueue <- job

		jl.currentJobCounter.addJobCount(1)
	}
}
