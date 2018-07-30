package repo

import (
	"errors"
)

var updateQueue chan UpdateRequest

// UpdateRequest ...
type UpdateRequest struct {
	Slug string
	Repo string
}

func init() {
	if updateQueue == nil {
		updateQueue = make(chan UpdateRequest, 200000)
	}
}

// StartUpdateWorkers starts Goroutines to process Plugin and Theme updates
func StartUpdateWorkers(num int, pr *Repo, tr *Repo) {
	for i := 0; i < num; i++ {
		go func(queue <-chan UpdateRequest, pr *Repo, tr *Repo) {
			for {
				ur := <-queue
				var err error
				switch ur.Repo {
				case "plugins":
					err = pr.ProcessUpdate(ur.Slug)
				case "themes":
					err = tr.ProcessUpdate(ur.Slug)
				default:
					err = errors.New("Update failed, Repo not recognized")
				}
				if err != nil {
					// Use the logger embedded into the Plugins Repo
					pr.log.Printf("Update failed for %s (%s)", ur.Slug, ur.Repo)
				}
			}
		}(updateQueue, pr, tr)
	}
}
