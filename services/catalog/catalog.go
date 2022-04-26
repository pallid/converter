package catalog

import (
	"context"
	"converter/docs"
	"converter/pdf"
	"converter/processor"
	"converter/utils"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"

	"syscall"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/shirou/gopsutil/cpu"
)

var impl CatalogService

type shutdown struct{}

type CatalogService struct {
	startPath   string
	targetPath  string
	converter   pdf.Converter
	shutdown    chan shutdown
	shutdownErr chan shutdown
	metrics     *metrics
	wp          *workerpool.WorkerPool
}

type task struct {
	ctx context.Context
	e   *processor.Entity
}

func Start(startPath, targetPath string, c pdf.Converter) {

	log.Printf("[INFO] Запуск конвертации")

	logocal := false
	MAX_PROC, err := cpu.Counts(logocal)
	if err != nil {
		log.Printf("[WARNING] ошибка определения количества потоков: %v", err)
		MAX_PROC = 1
	}
	runtime.GOMAXPROCS(MAX_PROC)
	log.Printf("[DEBUG] количество потоков: %d", MAX_PROC)

	impl = CatalogService{
		startPath:  startPath,
		targetPath: targetPath,
		converter:  c,

		// this channel is for graceful shutdown:
		// if we receive an error, we can send it here to notify the server to be stopped
		shutdown:    make(chan shutdown, 1),
		shutdownErr: make(chan shutdown, 1),

		metrics: &metrics{
			StartDate: time.Now(),
		},
		wp: workerpool.New(MAX_PROC),
	}

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
		} else if !f.IsDir() && isSupportDocument(f.Name()) {
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
	cs.wp.Submit(func() {
		cs.processing(task{e: n})
	})

}

// закрывает канал при вызове
func (cs *CatalogService) Stop() {
	cs.wp.StopWait()
	log.Println(cs.metrics.GetStatistic())
}

//processing
func (cs *CatalogService) processing(t task) {
	log.Printf(" - convert start: %v\n", t.e.SourceFile)
	if err := t.e.Convert(); err != nil {
		cs.metrics.IncreaseError()
		log.Println("[Processing] error: ", err)
		cs.shutdownErr <- shutdown{}
	}
	log.Printf(" - convert finish: %v\n", t.e.SourceFile)
	cs.metrics.IncreaseDone()
	if cs.metrics.Find == cs.metrics.Done+cs.metrics.Error {
		cs.metrics.FinishDate = time.Now()
		cs.shutdown <- shutdown{}
	}
}

var _supportDocuments = []docs.DocumentFormat{
	docs.PDF,
}

func isSupportDocument(fileName string) bool {
	ext := utils.GetExtensionFile(fileName)
	extDF := docs.DocumentFormat(strings.TrimPrefix(ext, "."))
	for _, doc := range _supportDocuments {
		if extDF == doc {
			return true
		}
	}
	return false
}
