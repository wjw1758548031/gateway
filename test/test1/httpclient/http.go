package httpclient

import (
	"crypto/tls"
	"errors"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var HttpClient Http

type Http struct {
	//代理列表
	Proxys  []string
	Timeout int
}

func (this *Http) randProxy() (string, error) {
	if len(this.Proxys) == 0 {
		return "", errors.New("代理不存在")
	}
	n := rand.Intn(len(this.Proxys))
	return this.Proxys[n], nil
}
func (this *Http) Get(urlStr string, proxy bool) (string, error) {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy == true {
		proxyUrl, err := this.randProxy()
		if err != nil {
			logger.Error(err)
			return "", err
		}
		urli := url.URL{}
		urlproxy, _ := urli.Parse(proxyUrl)
		tr.Proxy = http.ProxyURL(urlproxy)
	}
	c := http.Client{
		Transport: &tr,
		Timeout:   time.Second * time.Duration(this.Timeout),
	}
	if resp, err := c.Get(urlStr); err != nil {
		//log.Error(err)
		return "", err
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), err
	}
}
func (this *Http) GetHeader(url string, headers map[string]string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(this.Timeout),
	}
	//提交请求
	reqest, err := http.NewRequest("GET", url, nil)
	//增加header选项
	for k, v := range headers {
		reqest.Header.Add(k, v)
	}
	if err != nil {
		return "", err
	}
	//处理返回结果
	if response, err := client.Do(reqest); err != nil {
		return "", err
	} else {
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		return string(data), nil
	}
}
func (this *Http) Post(urlStr string, post string, proxy bool,token map[string]string) (string, error) {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy == true {
		proxyUrl, err := this.randProxy()
		if err != nil {
			logger.Error(err)
			return "", err
		}
		urli := url.URL{}
		urlproxy, _ := urli.Parse(proxyUrl)
		tr.Proxy = http.ProxyURL(urlproxy)
	}
	c := http.Client{
		Transport: &tr,
		Timeout:   time.Second * time.Duration(this.Timeout),
	}
	resp, err := c.Post(urlStr, "application/x-www-form-urlencoded", strings.NewReader(post))
	resp.Header = http.Header{}
	for k , v :=range token{
		resp.Header.Add(k,v)
	}
	if  err != nil {
		//log.Error(err)
		return "", err
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), err
	}
}

func (this *Http) PostClient(urlStr string, post string, proxy bool,token map[string]string) (string, error) {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy == true {
		proxyUrl, err := this.randProxy()
		if err != nil {
			logger.Error(err)
			return "", err
		}
		urli := url.URL{}
		urlproxy, _ := urli.Parse(proxyUrl)
		tr.Proxy = http.ProxyURL(urlproxy)
	}
	c := http.Client{
		Transport: &tr,
		Timeout:   time.Second * time.Duration(this.Timeout),
	}

	req, err := http.NewRequest("POST",urlStr,  strings.NewReader(post))
	if err != nil {
		return "", err
		// handle error
	}
	//resp, err := c.Post(urlStr, "application/x-www-form-urlencoded", strings.NewReader(post))
	//resp.Header = http.Header{}
	//请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k , v :=range token{
		req.Header.Set(k,v)
	}
	resp, err := c.Do(req)
	if  err != nil {
		//log.Error(err)
		return "", err
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), err
	}
}

func (this *Http) PostHeader(url string, data string, headers map[string]string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(this.Timeout),
	}
	//提交请求
	reqest, err := http.NewRequest("POST", url, strings.NewReader(data))
	//增加header选项
	for k, v := range headers {
		reqest.Header.Add(k, v)
	}
	if err != nil {
		return "", err
	}
	//处理返回结果
	if response, err := client.Do(reqest); err != nil {
		return "", err
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return string(body), nil
	}
}

func (this *Http) MapToQueryString(queryMap map[string]string) (string, error) {
	queryValues := url.Values{}
	for k, v := range queryMap {
		queryValues.Add(k, string(v))
	}
	logger.Debug(queryValues.Encode())
	return queryValues.Encode(), nil
}
