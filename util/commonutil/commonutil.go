package commonutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

var configInstance ConfigIns
var logInstance LogIns

var lockLogIns *sync.Mutex = &sync.Mutex{}
var lockConfigIns *sync.Mutex = &sync.Mutex{}
var once sync.Once

type ConfigIns interface {
	String(key string) string
}

type LogIns interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
}

type logStruct struct {
	log *logs.BeeLogger
}

type configuration struct {
	conf config.Configer
}

func (this *configuration) String(key string) string {
	return this.conf.String(key)
}

func (this *logStruct) Debug(format string, v ...interface{}) {
	this.log.Debug(format)
}

func (this *logStruct) Info(format string, v ...interface{}) {
	this.log.Info(format)
}

func LogInstance() LogIns {

	lockLogIns.Lock()
	defer lockLogIns.Unlock()

	if logInstance != nil {
		return logInstance
	} else {
		log := logs.NewLogger(10000)
		log.SetLogger("file", `{"filename":"test.log"}`)
		return &logStruct{log}
	}

}

func ConfigInstance(adapterName, filename string) ConfigIns {
	lockConfigIns.Lock()
	defer lockConfigIns.Unlock()
	if configInstance != nil {
		return configInstance
	} else {
		conf, err := config.NewConfig(adapterName, filename)
		if err != nil {
			return nil
		}
		configInstance = &configuration{conf}
		return configInstance
	}

}

func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}
