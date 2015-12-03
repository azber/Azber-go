package xiaoqiao

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	IP_ADDRESS = "117.28.254.159"
	TEST_URL   = "https://www.google.com"
)

type Xiaoqiao struct {
	port int
}

func NewXiaoqiao(port int) (*Xiaoqiao, error) {
	return &Xiaoqiao{
		port: port,
	}, nil
}

func (x *Xiaoqiao) Proxy() {
	//	fmt.Println("begin = " + strconv.Itoa(x.port))

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + IP_ADDRESS + ":" + strconv.Itoa(x.port))
	}
	transport := &http.Transport{Proxy: proxy}
	httpClient := &http.Client{
		Transport: transport,
	}

	response, err := httpClient.Get(TEST_URL)
	if err != nil {
		return
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	respStr := string(data)
	if strings.Contains(respStr, "google") {
		fmt.Println("port = " + strconv.Itoa(x.port))
	}
}
