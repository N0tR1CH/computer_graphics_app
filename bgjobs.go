package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Job struct {
	ID       string
	FilePath string
	comments []string
}

type Worker struct {
	jobQueue  chan Job
	jobStatus map[string]string
	lock      sync.Mutex
	app       *App
}

func NewWorker(app *App) *Worker {
	worker := &Worker{
		jobQueue:  make(chan Job, 1), // Buffer size set to 1
		jobStatus: make(map[string]string),
		app:       app,
	}
	go worker.processJobs()
	return worker
}

func (w *Worker) UploadNetPbmImg() string {
	filepath, err := runtime.OpenFileDialog(w.app.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return ""
	}
	if filepath == "" {
		return filepath
	}
	fmt.Println("Selected file:", filepath)

	jobID := uuid.New().String()
	job := Job{ID: jobID, FilePath: filepath}

	w.lock.Lock()
	w.jobStatus[jobID] = "queued"
	w.lock.Unlock()

	w.jobQueue <- job
	return jobID
}

func (w *Worker) processJobs() {
	for job := range w.jobQueue {
		w.updateStatus(job.ID, "processing")

		fmt.Println("Processing file:", job.FilePath)

		file, err := os.Open(job.FilePath)
		if err != nil {
			log.Fatal(err)
		}

		img, comments, err := parseNetPbm(file)
		if err != nil {
			w.updateStatus(job.ID, "failed")
		}
		if err != nil {
			log.Fatal(err)
		}

		outFile, err := os.Create(fmt.Sprintf("./images/file%s.jpeg", job.ID))
		if err != nil {
			log.Fatal(err)
		}
		if err := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 50}); err != nil {
			panic(err)
		}
		if err := outFile.Close(); err != nil {
			panic(err)
		}

		job.comments = comments

		w.updateStatus(job.ID, "completed")
		runtime.EventsEmit(w.app.ctx, job.ID, job.comments, w.jobStatus[job.ID])
	}
}

func (w *Worker) updateStatus(jobID, status string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.jobStatus[jobID] = status
}

func (w *Worker) GetJobStatus(jobID string) string {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.jobStatus[jobID]
}
