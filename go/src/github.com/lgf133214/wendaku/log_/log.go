package log_

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

var filePath = ""

func init() {
	s, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath = s
}

func ErrWrite(err interface{}) {
	if data, ok := err.(error); ok {
		data := []byte(time.Now().Format("2006 01/02 03:04:05PM ") + data.Error() + "\n")

		fl, err := os.OpenFile(path.Join(filePath, "error.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		defer fl.Close()

		_, err = fl.Write(data)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	} else {
		data := []byte(time.Now().Format("2006 01/02 03:04:05PM ") + " type is not an error " + fmt.Sprintf("%v", err) + "\n")

		fl, err := os.OpenFile(path.Join(filePath, "error.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		defer fl.Close()

		_, err = fl.Write(data)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}
}
