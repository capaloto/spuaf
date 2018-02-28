package spuaf

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	regUAs    = regexp.MustCompile("(?s)\\<td[ ]class\\=\"useragent\"\\>(.+?)\\<\\/td")
	headerSet = make([]map[string]string, 0)
)

func Init() error {
	rand.Seed(time.Now().UTC().UnixNano())
	resp, err := http.Get("https://techblog.willshouse.com/2012/01/03/most-common-user-agents/")
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("bad status - %s", resp.Status)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	for _, uaMatch := range regUAs.FindAllSubmatch(bytes, -1) {
		uaString := strings.Trim(string(uaMatch[len(uaMatch)-1]), " \n\t")
		headerSet = append(headerSet, map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.8",
			"User-Agent":      uaString,
		})
	}
	return nil
}

func Spuaf(req *http.Request) error {
	if len(headerSet) < 1 {
		return fmt.Errorf("No headers found did you init Spuaf?")
	}
	headers := headerSet[rand.Intn(len(headerSet))]
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return nil
}
