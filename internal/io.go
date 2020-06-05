package exporter

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func GetUuidContent(filepath string) string {
	for !Exists(filepath) {
		log.Println("-----wait for xzcom-exporter to write output.csv'")
		time.Sleep(1 * time.Second)
	}

	fin, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	bytes, err := ioutil.ReadAll(fin)
	if err != nil {
		panic(err)
	}
	csccontent := string(bytes)
	//fmt.Println(csccontent)
	csccontent = strings.Replace(csccontent, "\n", "", -1)

	return csccontent
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
