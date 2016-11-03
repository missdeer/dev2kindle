package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Toutiao struct {
}

func (t *Toutiao) resolveFinalURL(u string) string {
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("resolving url %s failed => %v", u, err.Error())
	}

	finalURL := resp.Request.URL.String()
	return finalURL
}

func (t *Toutiao) Fetch(link chan string) {
	now := time.Now()
	u := fmt.Sprintf("https://toutiao.io/prev/%4.4d-%2.2d-%2.2d", now.Year(), now.Month(), now.Day())

	retry := 0
doRequest:
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Println("Could not parse get page list request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not send get page list request:", err)
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doRequest
		}
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("get page list request not 200")
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doRequest
		}
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("cannot read get page list content", err)
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doRequest
		}
	}

	regex := regexp.MustCompile(`/k/([0-9a-zA-Z]+)`)
	list := regex.FindAllSubmatch(content, -1)
	for _, l := range list {
		lnk := fmt.Sprintf("https://toutiao.io/k/%s", string(l[1]))
		link <- t.resolveFinalURL(lnk)
	}
}