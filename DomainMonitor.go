package main

import (
	"github.com/likexian/whois-go"
	"github.com/likexian/whois-parser-go"
	"io/ioutil"
	"strings"
	"log"
	"net/http"
	"time"
	"strconv"

	"github.com/kardianos/osext"
	"github.com/araddon/dateparse"
)

func main() {
	path,_ := osext.ExecutableFolder()
	domainFile, err := ioutil.ReadFile(path + "/domain")
	if err != nil{
		log.Fatalln("打开文件失败")
	}

	var min int64 = 9999999
	var minBody string
	var count = 0
	var totalCount = 0
	strs := strings.Split(string(domainFile), "\n")
	for _, line := range strs{
		if len(line) <= 0 {
			continue
		}
		totalCount++

		result, err := whois.Whois(line)
		if err != nil {
			sendFailedPush(line, err.Error())
			continue
		}

		r, e := whois_parser.Parse(result)
		if e != nil {
			sendFailedPush(line, e.Error())
			continue
		}

		expTime := dateFormat(r.Registrar.ExpirationDate)
		updateTime := dateFormat(r.Registrar.UpdatedDate)

		a := (expTime.Unix() - time.Now().Unix()) / 24 / 60 / 60

		body := line + "  " + strconv.Itoa(int(a)) + "天" + "\n过期时间: " + expTime.Format("2006-01-02 15:04:05") + "\n更新时间: " + updateTime.Format("2006-01-02 15:04:05")

		if a <= 7 {
			sendNotification(body)
		}

		if a < min {
			min = a
			minBody = body
		}

		count++
	}

	sendNotification(strconv.Itoa(count) + "个域名(总共" + strconv.Itoa(totalCount) + "个)查询成功. \n" + minBody)

}
func dateFormat(date string)time.Time {
	t,e := dateparse.ParseAny(date)
	if e != nil {
		return time.Now()
	}
	return t
}

func sendFailedPush(domain string, reason string){
	sendNotification("域名: " + domain + " 查询失败 reason: " + reason)
}

func sendNotification(body string){
	url := "https://api.uusing.com/PeNX4RrNFYkwgQYx8jYKck/alert"
	http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("body="+body))
}

