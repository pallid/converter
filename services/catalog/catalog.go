package catalog

import (
	"context"
	"converter/pdf"
	"converter/processor"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"

	"syscall"
	"time"
)

const (
	// Время простоя таймера перед
	idleTimeoutSec = 10
)

var impl CatalogService

type shutdown struct{}

type CatalogService struct {
	startPath   string
	targetPath  string
	converter   pdf.Converter
	qmQueue     chan task
	queueClose  bool
	timeout     time.Duration
	shutdown    chan shutdown
	shutdownErr chan shutdown
	metrics     *metrics
}

type task struct {
	ctx context.Context
	e   *processor.Entity
}

func Start(startPath, targetPath string, c pdf.Converter, timeoutSec time.Duration) {

	if timeoutSec == 0 {
		timeoutSec = idleTimeoutSec
	}

	impl = CatalogService{
		startPath:  startPath,
		targetPath: targetPath,
		converter:  c,
		qmQueue:    make(chan task),
		queueClose: false,
		timeout:    timeoutSec,

		// this channel is for graceful shutdown:
		// if we receive an error, we can send it here to notify the server to be stopped
		shutdown:    make(chan shutdown, 1),
		shutdownErr: make(chan shutdown, 1),

		metrics: &metrics{
			StartDate: time.Now(),
		},
	}

	log.Printf("[INFO] Запуск конвертации")
	go impl.dispatch()
	log.Printf("[INFO] Исходный каталог: %s", impl.startPath)
	impl.copyDirectory(impl.startPath, impl.targetPath)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			log.Print("[INFO] Получил SIGINT...")
		case syscall.SIGTERM:
			log.Print("[INFO] Получил SIGTERM...")
		}
	case <-impl.shutdownErr:
		log.Printf("[ERR] Получил ошибку...")

	case <-impl.shutdown:
		log.Printf("[INFO] Конвертация завершена...")
	}

	impl.Stop()
	log.Print("[INFO] Завершение работы сервиса...")
}

func (cs CatalogService) copyDirectory(src, dst string) {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatalln(err)
	}

	for _, f := range files {
		newSrc := path.Join(src, f.Name())
		newDst := path.Join(dst, f.Name())

		if f.IsDir() {
			cs.copyDirectory(newSrc, newDst)
		} else if !f.IsDir() {
			cs.metrics.IncreaseFind()
			log.Printf(" - task to query: %s\n", newSrc)
			cs.Submit(&processor.Entity{
				SourceFile:        newSrc,
				TargetFolder:      newDst,
				ConvertFileFormat: processor.DefaultFormat,
				Converter:         impl.converter,
			})

		}
	}

}

func (cs *CatalogService) Submit(n *processor.Entity) {
	cs.qmQueue <- task{
		e: n,
	}
}

// закрывает канал при вызове
func (cs *CatalogService) Stop() {
	close(cs.qmQueue)
	cs.queueClose = true
	log.Println(cs.metrics.GetStatistic())
}

func (cs *CatalogService) dispatch() {
	timeout := time.NewTimer(time.Second * cs.timeout)
	var (
		arr        []task
		curTask    task
		ok         bool
		inProgress bool
	)
Loop:
	for {
		timeout.Reset(cs.timeout)
		select {
		case curTask, ok = <-cs.qmQueue:
			if !ok {
				break Loop
			}
			arr = append(arr, curTask)
		case <-timeout.C:
			if inProgress {
				break Loop
			}
			inProgress = true
			if len(arr) != 0 {
				if err := cs.processing(arr); err != nil {
					log.Println("[CONVERT] error: ", err)
					cs.shutdown <- shutdown{}
				}
				arr = make([]task, 0)
			}
			inProgress = false
			if cs.metrics.Find == cs.metrics.Done+cs.metrics.Error && len(cs.qmQueue) == 0 {
				cs.metrics.FinishDate = time.Now()
				cs.shutdown <- shutdown{}
			}
		}
	}
}

//processing
func (cs *CatalogService) processing(tasks []task) error {
	for _, t := range tasks {
		log.Printf(" - convert start: %v\n", t.e.SourceFile)
		if err := t.e.Convert(); err != nil {
			cs.metrics.IncreaseError()
			return err
		}
		log.Printf(" - convert finish: %v\n", t.e.SourceFile)
		cs.metrics.IncreaseDone()
	}
	return nil
}
