package exporter

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

func (e *Exporter) gatherData() (SSLInfoArray, error) {

	var data SSLInfoArray
	log.Println(e.Config.Domain)
	wg := sync.WaitGroup{}
	log.Println(len(e.Config.Domain))
	wg.Add(len(e.Config.Domain))

	for _, sslname := range e.Config.Domain {
		go func() {
			result := VerifySSL(sslname)
			log.Println(result)
			MergeSlice(data, result, &wg)
		}()
	}
	wg.Wait()
	return data, nil

}

func VerifySSL(sslname string) SSLInfoArray {

	var sslstruct SSLInfo
	var test SSLInfoArray

	stedout, err := ExecShell("echo | openssl s_client -servername " + sslname + " -connect " + sslname + ":443 2>/dev/null | openssl x509 -noout -dates")
	if err != "" {
		log.Fatalln(err)
		//panic("")
	}
	result := strings.Split(stedout, "\n")
	log.Println(result)

	beforedate := strings.Split(result[0], "=")
	afterdate := strings.Split(result[1], "=")
	//log.Println("------beforedate[1]-----")
	//log.Println(beforedate[1])
	//og.Println("------afterdate[1]-----")
	//log.Println(afterdate[1])

	dataBefore := ParserDateFormat(beforedate[1])
	//log.Println("------dataBefore-----")
	//log.Println(dataBefore)
	dataAfter := ParserDateFormat(afterdate[1])
	//log.Println("------dataAfter-----")
	//log.Println(dataAfter)
	(&sslstruct).SetDomainNmme(sslname)
	(&sslstruct).SetRegistryDate(dataBefore)
	(&sslstruct).SetExpiredDate(dataAfter)

	remain := timeSubDays(dataAfter, dataBefore)
	strremain := strconv.Itoa(remain)
	f64remain, _ := strconv.ParseFloat(strremain, 64)
	(&sslstruct).RemainingDate(f64remain)

	test = append(test, sslstruct)
	return test
}

func ParserDateFormat(date string) time.Time {

	layout := "2006-01-02 15:04:05"
	//date --date="May 30 00:00:00 2020 GMT" --utc +"%Y-%m-%d %T"
	u, _ := uuid.NewV4()
	log.Println("date --date=\"" + date + "\" --utc +\"%Y-%m-%d %T\" > " + u.String())
	_, err := ExecShell("date --date=\"" + date + "\" --utc +\"%Y-%m-%d %T\" > " + u.String())
	if err != "" {
		log.Fatalln(err)
	}

	linuxdate := GetUuidContent(u.String())
	t, _ := time.Parse(layout, linuxdate)
	return t

}

func timeSubDays(t1, t2 time.Time) int {
	log.Println("------timeSubDays--------")
	if t1.Location().String() != t2.Location().String() {
		return -1
	}
	hours := t1.Sub(t2).Hours()

	if hours <= 0 {
		log.Println("------hours <= 0--------")
		return -1
	}
	// sub hours less than 24
	if hours < 24 {
		// may same day
		t1y, t1m, t1d := t1.Date()
		t2y, t2m, t2d := t2.Date()
		isSameDay := (t1y == t2y && t1m == t2m && t1d == t2d)

		if isSameDay {

			return 0
		} else {
			return 1
		}

	} else { // equal or more than 24

		if (hours/24)-float64(int(hours/24)) == 0 { // just 24's times
			return int(hours / 24)
		} else { // more than 24 hours
			return int(hours/24) + 1
		}
	}

}