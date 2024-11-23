package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/image/webp"
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

		var (
			img      image.Image
			comments []string
		)
		switch filepath.Ext(job.FilePath) {
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(file)
		case ".png":
			img, err = png.Decode(file)
		case ".webp":
			img, err = webp.Decode(file)
		case ".pbm", ".pgm", ".ppm", ".pnm":
			img, comments, err = parseNetPbm(file)
		default:
			w.updateStatus(job.ID, "failed")
		}

		if err != nil {
			w.updateStatus(job.ID, "failed")
		}
		if err != nil {
			log.Fatal(err)
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			panic(err)
		}

		base64str := base64.StdEncoding.EncodeToString(buf.Bytes())

		job.comments = comments

		w.updateStatus(job.ID, "completed")
		runtime.EventsEmit(w.app.ctx, job.ID, job.comments, w.jobStatus[job.ID], base64str)
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
