package exporter

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	uuid "github.com/nu7hatch/gouuid"

	"github.com/siangyeh8818/ssl.effective.date.exporter/internal/database"
)

func (e *Exporter) gatherData() (SSLInfoArray, error) {
	log.Println("-------gatherData()------")
	var data SSLInfoArray
	log.Println("-------e.Config.Domain-------")
	log.Println(e.Config.Domain)

	log.Println("-------len(e.Config.Domain)----------")
	log.Println(len(e.Config.Domain))

	var targetDomains []string
	// gaia domains
	for _, domainName := range e.Config.Domain {
		targetDomains = append(targetDomains, domainName)
	}

	// cloudflare domains, fetch redis
	cloudflareKeys := database.Keys(fmt.Sprintf("%s*", database.CloudfalrePrefix))
	for _, key := range cloudflareKeys {
		dName := key[len(database.CloudfalrePrefix):]
		targetDomains = append(targetDomains, dName)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(targetDomains))

	var mux sync.Mutex
	//for i, sslname := range e.Config.Domain {
	for i := 0; i < len(targetDomains); i++ {
		go func(i int) {
			mux.Lock()
			log.Println(i)
			log.Println(targetDomains[i])
			result := VerifySSL(targetDomains[i])
			log.Println(result)
			data = MergeSlice(data, result, &wg)
			mux.Unlock()
		}(i)
	}
	wg.Wait()
	log.Println(data)
	return data, nil

}

func VerifySSL(sslname string) SSLInfoArray {

	var sslstruct SSLInfo
	var test SSLInfoArray

	stedout, err := ExecShell("echo | openssl s_client -servername " + sslname + " -connect " + sslname + ":443 2>/dev/null | openssl x509 -noout -dates")
	if err != "" {
		// Load certificate error
		log.Printf("[Error] Failed to laod certificate of server name: %s, reason: %v", sslname, err)
		//panic("")
	} else {
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

		currentTime := time.Now().UTC()
		currentTime.Location()

		remain := timeSubDays(dataAfter, currentTime)
		if remain < 0 {
			sslstruct.ExpiryStatus = "expired"
		} else if remain >= 0 && remain < 30 {
			sslstruct.ExpiryStatus = "expiring"
		} else {
			sslstruct.ExpiryStatus = "registered"
		}
		strremain := strconv.Itoa(remain)
		f64remain, _ := strconv.ParseFloat(strremain, 64)
		(&sslstruct).RemainingDate(f64remain)
		test = append(test, sslstruct)
	}

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
	defer os.Remove(u.String())
	linuxdate := GetUuidContent(u.String())
	t, _ := time.Parse(layout, linuxdate)
	return t

}

func timeSubDays(t1, t2 time.Time) int {

	log.Println("------timeSubDays--------")
	//local1, _ := time.LoadLocation("Asia/Taipei")
	//log.Println(t1.In(local1).Format("2006-01-02 15:04:05"))
	//log.Println(t2.In(local1).Format("2006-01-02 15:04:05"))
	if t1.Location().String() != t2.Location().String() {
		//log.Println(t1.Location().String())
		//log.Println(t2.Location().String())
		return -1
	}
	hours := t1.Sub(t2).Hours()
	log.Println("------hours--------")
	log.Println(hours)
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
