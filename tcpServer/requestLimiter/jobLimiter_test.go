package requestlimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimiterQPSLimit(t *testing.T) {
	jl := &JobLimiter{}
	// set 1 query/sec
	jl.Init(1)
	jl.SetExternalAPI("mockExternalAPI")

	job := &Job{}
	job.SetNetConnection(nil)
	job.SetHost("mockHost")
	job.SetQuery("mockQuery")

	for i := 1; i <= 3; i++ {
		jl.EnqueueJob(job)
	}

	time.Sleep(1 * time.Second)

	assert.NotEqual(t, 0, jl.GetRemainingJobCount(), "Remaining job count should not equal 0")
}

func TestLimiterCompleteJobOnTime(t *testing.T) {
	jl := &JobLimiter{}
	// set 1 query/sec
	jl.Init(10)
	jl.SetExternalAPI("mockExternalAPI")

	job := &Job{}
	job.SetNetConnection(nil)
	job.SetHost("mockHost")
	job.SetQuery("mockQuery")

	for i := 1; i <= 20; i++ {
		jl.EnqueueJob(job)
	}

	time.Sleep(3 * time.Second)

	assert.Equal(t, 0, jl.GetRemainingJobCount(), "Remaining job count should equal 0")
}
