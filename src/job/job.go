package job

import (
	"fmt"
	"sync"
	"time"
)

// Job represents a batch job with an ID and a list of dependent job IDs
type Job struct {
	ID           int
	Name         string
	Dependencies []*Job
	Task         func()
}

// Execute runs the job logic (in this case, it just prints the job ID and simulates a delay)
func (j *Job) Execute(wg *sync.WaitGroup, done chan int) {
	defer wg.Done() // Signal that the job is done
	fmt.Printf("Executing job %s\n", j.Name)
	//time.Sleep(1 * time.Second) // Simulate work
	j.Task()
	done <- j.ID // Notify that the job is finished
}

// RunJobs runs the jobs in batches based on their dependencies
func RunJobs(jobs []*Job) {
	var wg sync.WaitGroup
	done := make(chan int, len(jobs)) // A channel to track finished jobs

	jobMap := make(map[int]*Job) // Maps job IDs to job instances
	for _, job := range jobs {
		jobMap[job.ID] = job
	}

	// Track jobs that are ready to run (i.e., no unmet dependencies)
	readyToRun := make(chan *Job)
	go func() {
		for _, job := range jobs {
			if len(job.Dependencies) == 0 {
				readyToRun <- job
			}
		}
	}()

	time.Sleep(time.Second * 1)

	// Start workers to execute jobs as they become ready
	go func() {
		for job := range readyToRun {
			wg.Add(1)
			go job.Execute(&wg, done)
		}
	}()

	time.Sleep(time.Second * 1)

	// Listen for finished jobs and check for dependent jobs that can now run
	go func() {
		for finishedJobID := range done {
			// Check dependent jobs to see if all their dependencies are met
			for _, job := range jobs {
				for i, dep := range job.Dependencies {
					if dep.ID == finishedJobID {
						// Remove the finished dependency
						job.Dependencies = append(job.Dependencies[:i], job.Dependencies[i+1:]...)
					}
				}

				// If this job now has no dependencies, mark it as ready to run
				if len(job.Dependencies) == 0 {
					readyToRun <- job
				}
			}
		}
	}()

	// Wait for all jobs to finish
	wg.Wait()
	close(done)       // Close the done channel after all jobs are finished
	close(readyToRun) // Close the readyToRun channel
}
