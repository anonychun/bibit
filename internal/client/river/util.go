package river

import "github.com/riverqueue/river"

func addJobs[T river.JobArgs](workers *river.Workers, jobs ...river.Worker[T]) error {
	for _, job := range jobs {
		err := river.AddWorkerSafely(workers, job)
		if err != nil {
			return err
		}
	}

	return nil
}
