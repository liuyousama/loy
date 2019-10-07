package text_handler

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	textChannelLength = 128
	maxRetryTimes     = 3
	fileCheckTime     = 10
	rollingByTime     = "time"
	rollingBySize     = "size"
	kb                = 1024
	mb                = 1024 * kb
	gb                = 1024 * mb
)

type FileHandler struct {
	outputFile          *os.File
	textChan            chan string
	circleCheckTimes    int8
	lastRollingTime     time.Time
	rollingTimeDuration time.Duration
	rollingFileSize     int64
	rollingPolicy       string
}

func init()  {
	Handlers["file"] = new(FileHandler)
}

func (h *FileHandler) LoadHandler(option HandlerOption) error {
	file, err := os.OpenFile(option.LogFilePath, os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		return err
	}
	h.outputFile = file

	if strings.ToLower(strings.TrimSpace(option.RollingPolicy)) == rollingBySize {
		h.rollingPolicy = rollingBySize
	} else if strings.ToLower(strings.TrimSpace(option.RollingPolicy)) == rollingByTime {
		h.rollingPolicy = rollingByTime
	} else {
		h.rollingPolicy = rollingBySize
	}

	h.rollingFileSize= option.RollingSize
	h.rollingPolicy = option.RollingPolicy
	h.rollingTimeDuration = option.RollingDuration

	h.textChan = make(chan string, textChannelLength)

	go h.handleText()

	return nil
}

func (h *FileHandler) HandleText(text string, level, minLevel int) {
	if level > minLevel {
		return
	}

	if h.textChan == nil {
		h.textChan = make(chan string, textChannelLength)
	}

	select {
	case h.textChan <- text:
		return
	case <-time.Tick(1 * time.Second):
		return
	}

}

func (h *FileHandler) handleText() {
	for {
		select {
		case text := <-h.textChan:
			retryExecutor(func() error {
				_, err := fmt.Fprintf(h.outputFile, text)
				return err
			})

			h.incrCheckTimes()
		case <-time.Tick(300 * time.Millisecond):
			continue
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
	if h.rollingPolicy == rollingByTime {
		h.checkFileTime()
	} else if h.rollingPolicy == rollingBySize {
		h.checkFileSize()
	} else {
		h.rollingPolicy = rollingBySize
		h.checkFileSize()
	}
}

func (h *FileHandler) checkFileSize() {
	if h.rollingFileSize < 1*mb {
		h.rollingFileSize = 1 * mb
	}

	fileStat, err := h.outputFile.Stat()
	if err != nil {
		return
	}

	if h.rollingFileSize > fileStat.Size() {
		h.updateLoggerFile()
	}
}

func (h *FileHandler) checkFileTime() {
	if h.rollingTimeDuration < 24*time.Hour {
		h.rollingTimeDuration = 24 * time.Hour
	}

	if h.lastRollingTime.Add(h.rollingTimeDuration).Before(time.Now()) {
		h.updateLoggerFile()
	}
}

func (h *FileHandler) updateLoggerFile() {
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
	target := fmt.Sprintf("%s.%s.zip", source, time.Now().Format("20060102150405999"))
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