package exporter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Handler  http.Handler
	exporter Exporter
}

type BaseConfig struct {
	Domain []string `json:"domains"`
	//INTERVAL_TIME string
}

type SSLInfoArray []SSLInfo

type SSLInfo struct {
	DomainNmme       string
	SSLRemainingDate float64
	ExpiredDate      time.Time
	RegistryDate     time.Time
}

type JsonStruct struct {
	Name []string `json:"domains"`
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (s *SSLInfo) SetDomainNmme(name string) {
	s.DomainNmme = name
}
func (s *SSLInfo) SetExpiredDate(date time.Time) {
	s.ExpiredDate = date
}
func (s *SSLInfo) SetRegistryDate(date time.Time) {
	s.RegistryDate = date
}
func (s *SSLInfo) RemainingDate(count float64) {
	s.SSLRemainingDate = count
}

func (con *BaseConfig) Initconfig(filepath string) {

	for !Exists(filepath) {
		log.Println("-----wait for gaiaDomains.json , you have to mount 'gaiaDomains.json to container")
		time.Sleep(10 * time.Second)
	}
	//Load(filepath, newconf.Domain)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(data))
	err = json.Unmarshal(data, con)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(con.Domain)
}

/*
func (con *BaseConfig) Initdoman(filepath string) {

	//var domain JsonStruct
	//JsonParse := NewJsonStruct()

	for !Exists(filepath) {
		log.Println("-----wait for gaiaDomains.json , you have to mount 'gaiaDomains.json to container")
		time.Sleep(10 * time.Second)
	}

	data, err := ioutil.ReadFile(filename)
	err = json.Unmarshal(data, v)
	//Load(filepath, con.Domain)
	//con.INTERVAL_TIME = os.Getenv("INTERVAL_TIME")

}
*/
/*
func Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}
*/
func MergeSlice(s1 SSLInfoArray, s2 SSLInfoArray, wg *sync.WaitGroup) []SSLInfo {
	defer wg.Done()
	slice := make([]SSLInfo, len(s1)+len(s2))

	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}
