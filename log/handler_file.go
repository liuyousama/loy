package log

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	textChannelLength = 128
	maxRetryTimes     = 3
	fileCheckTime     = 15
)

type RollingPolicy uint8
const (
	RollingBySize RollingPolicy = 0
	RollingByTime RollingPolicy = 1
)

type FileHandler struct {
	outputFile          *os.File
	outputFilePath		string
	rollingTimeDuration time.Duration
	rollingFileSize     int64
	rollingPolicy       RollingPolicy
	circleCheckTimes    int8
	lastRollingTime     time.Time
	textChan            chan string
}

func NewFileHandler() *FileHandler {
	h := &FileHandler{
		rollingTimeDuration: 24 * time.Hour,
		rollingFileSize: 1024 * 1024,
		rollingPolicy:RollingBySize,
	}

	return h
}

func (h *FileHandler) SetRollingSize(size int64) {
	if size < 1024 * 10 {
		size = 1024 * 10
	}
	h.rollingFileSize = size
}

func (h *FileHandler) SetRollingDuration(duration time.Duration) {
	if duration < 1 * time.Hour {
		duration = 1 * time.Hour
	}
	h.rollingTimeDuration = duration
}

func (h *FileHandler) SetRollingPolicy(policy RollingPolicy) {
	if policy != RollingBySize && policy != RollingByTime {
		policy = RollingBySize
	}
	h.rollingPolicy = policy
}

func (h *FileHandler) SetOutputFilePath(path string) {
	h.outputFilePath = path
}

func (h *FileHandler) Load() error {
	file, err := os.OpenFile(h.outputFilePath, os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		return err
	}
	h.outputFile = file

	h.textChan = make(chan string, textChannelLength)
	go h.handleText()

	return nil
}

func (h *FileHandler) HandleText(text string) {
	if h.textChan == nil {
		h.textChan = make(chan string, textChannelLength)
	}

	select {
	case h.textChan <- text:
		return
	//case <-time.Tick(1 * time.Second):
	//	return
	}

}

func (h *FileHandler) handleText() {
	for {
		select {
		case text := <-h.textChan:
			retryExecutor(func() error {
				_, err := fmt.Fprintln(h.outputFile, text)
				return err
			})

			h.incrCheckTimes()
		default:
			time.Sleep(50 * time.Microsecond)
		}
	}

}

func (h *FileHandler) incrCheckTimes() {
	if h.circleCheckTimes < fileCheckTime {
		h.circleCheckTimes++
		return
	} else {
		h.circleCheckTimes = 0
		h.checkFile()
	}
}

func (h *FileHandler) checkFile() {
	if h.rollingPolicy == RollingByTime {
		h.checkFileTime()
	} else if h.rollingPolicy == RollingBySize {
		h.checkFileSize()
	} else {
		h.rollingPolicy = RollingBySize
		h.checkFileSize()
	}
}

func (h *FileHandler) checkFileSize() {
	fileStat, err := h.outputFile.Stat()
	if err != nil {
		return
	}

	if h.rollingFileSize < fileStat.Size() {
		h.updateLoggerFile()
	}
}

func (h *FileHandler) checkFileTime() {
	if h.lastRollingTime.Add(h.rollingTimeDuration).Before(time.Now()) {
		h.updateLoggerFile()
	}
}

func (h *FileHandler) updateLoggerFile() {
	if time.Now().Before(h.lastRollingTime) {
		time.Sleep(1 * time.Second)
	}

	filePath := h.outputFile.Name()
	err := zipFile(h.outputFile)
	if err != nil {
		return
	}

	newFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0655)
	if err != nil {
		return
	}


	h.outputFile = newFile
	h.lastRollingTime = time.Now()
}

func zipFile(sourceFile *os.File) error {
	defer sourceFile.Close()

	source := sourceFile.Name()
	target := fmt.Sprintf("%s.%s.zip", source, time.Now().Format("20060102150405"))
	zipFile, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0440)
	if err != nil {
		log.Println(err)
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			header.Method = zip.Deflate
		}
		header.Modified = time.Now()
		header.Name = path
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)

		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func retryExecutor(f func() error) {
	var err error
	for i := 0; i < maxRetryTimes; i++ {
		err = f()
		if err != nil {
			continue
		} else {
			break
		}
	}
}