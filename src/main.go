package main

import (
	"fmt"

	"github.com/nedvisol/go-connectdots/job"
)

func main() {
	// Create jobs with dependencies

	job1 := &job.Job{
		ID:   1,
		Name: "first job",
		Task: func() { fmt.Println("job1") },
	}
	// job2 := &job.Job{ID: 2, Dependencies: []*job.Job{job1}}
	// job3 := &job.Job{ID: 3, Dependencies: []*job.Job{job1}}
	// job4 := &job.Job{ID: 4, Dependencies: []*job.Job{job2, job3}}

	// List of all jobs
	jobs := []*job.Job{job1}

	// Run jobs with dependencies
	job.RunJobs(jobs)

	// fmt.Println("All jobs executed.")
}
