package log

import (
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type dailyFileWriter struct {
	// 日志文件名称
	filename string
	// 上一次写入日期
	lastYearDay int
	outPutFile  *os.File
	// 考虑多线程，加锁
	lock *sync.Mutex
}

func (dfw *dailyFileWriter) Write(byteArray []byte) (n int, err error) {
	if byteArray == nil || len(byteArray) <= 0 {
		return 0, nil
	}
	outputFile, err := dfw.getOutPutFile()
	if err != nil {
		return 0, err
	}
	_, _ = os.Stderr.Write(byteArray)  // 屏幕输出
	_, _ = outputFile.Write(byteArray) // 文件输出
	return len(byteArray), nil
}

// 获取输出文件
// 每天创建一个新日志文件
func (dfw *dailyFileWriter) getOutPutFile() (io.Writer, error) {
	// 文件日期和当前时间一样，则还是为之前的文件句柄
	yearDay := time.Now().YearDay()
	if dfw.outPutFile != nil && yearDay == dfw.lastYearDay {
		return dfw.outPutFile, nil
	}
	if nil == dfw.lock {
		return nil, errors.New("lock is nil")
	}
	dfw.lock.Lock()
	defer dfw.lock.Unlock()
	// check double lock
	if dfw.outPutFile != nil && yearDay == dfw.lastYearDay {
		return dfw.outPutFile, nil
	}
	dfw.lastYearDay = yearDay
	// 创建目录，访问新的文件句柄
	err := os.MkdirAll(path.Dir(dfw.filename), os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "创建目录失败")
	}
	// 定义新的日志文件名称
	newDayFile := dfw.filename + "." + time.Now().Format("20060102")
	outputFile, err := os.OpenFile(newDayFile,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644)
	if err != nil || outputFile == nil {
		return nil, errors.Errorf("打开文件 %s 失败！ err = %v", outputFile, err)
	}
	// 处理句柄文件
	if dfw.outPutFile != nil {
		dfw.outPutFile.Close()
	}
	dfw.outPutFile = outputFile
	return dfw.outPutFile, nil
}
