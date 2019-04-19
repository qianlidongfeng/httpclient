package httpclient

import (
	"bytes"
	"crypto/tls"
	"github.com/headzoo/surf/errors"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)


type HttpClient struct{
	client http.Client
	header map[string]string
	cookiejar http.CookieJar
}

func NewHttpClient () (c HttpClient,err error){
	jar,err:=cookiejar.New(nil)
	if err != nil{
		_=jar
		return
	}
	header:=Headers.One()
	header["User-Agent"] = UserAgents.One()
	c = HttpClient{
		client:http.Client{
			Jar:nil,
			//Transport:&http.Transport{DisableKeepAlives: true,},
	},
		header:header,
		cookiejar:jar,
	}
	return
}

func (this *HttpClient) SetHeader(header map[string]string){
	this.header=make(map[string]string)
	for k,v := range header{
		this.header[k]=v
	}
}

func (this *HttpClient) UnsetHeader(){
	this.header=make(map[string]string)
}

func (this *HttpClient) SetHeaderField(key string,value string){
	this.header[key]=value
}

func (this *HttpClient) UnsetHeaderField(key string){
	delete(this.header,key)
}

func (this *HttpClient) SetHttpProxy(proxy string) error{
	proxyUrl, err := url.Parse(proxy)
	if err != nil{
		return err
	}
	this.client.Transport=&http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
		//DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   this.client.Timeout,
			KeepAlive: this.client.Timeout,
		}).DialContext,
		TLSHandshakeTimeout: this.client.Timeout,
		ExpectContinueTimeout:this.client.Timeout,
		IdleConnTimeout:this.client.Timeout,
	}
	return err
}

func (this *HttpClient) UnsetHttpProxy(){
	this.client.Transport = nil
	return
}

func (this *HttpClient) SetSocksProxy(prox string) error{
	dialSocksProxy, err := proxy.SOCKS5("tcp", prox, nil, proxy.Direct)
	//dialSocksProxy, err := proxy.SOCKS5("tcp", prox, nil, &net.Dialer { Timeout: 30 * time.Second, KeepAlive: 30 * time.Second})
	if err != nil{
		return err
	}
	this.client.Transport=&http.Transport{Dial:dialSocksProxy.Dial}
	return err
}

func (this *HttpClient) UnsetSocksProxy() {
	this.client.Transport=nil
}

func (this *HttpClient) EnableCookie(){
	this.client.Jar=this.cookiejar
}


func (this *HttpClient) DisableCookie(){
	this.client.Jar=nil
}

func (this *HttpClient) ClearCookie() error{
	jar,err:=cookiejar.New(nil)
	if err != nil{
		return err
	}
	this.client.Jar,this.cookiejar=jar,jar
	return nil
}

func (this *HttpClient) SetTimeOut(d time.Duration){
	this.client.Timeout=d
	this.client.Transport=&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   d,
			KeepAlive: d,
		}).DialContext,
		TLSHandshakeTimeout: d,
		ExpectContinueTimeout:d,
		IdleConnTimeout:d,
	}
}

func (this *HttpClient) UnSetTimeOut(){
	this.client.Timeout=0
}

func (this *HttpClient) SetCookies(targetUrl string, cookies []*http.Cookie) error{
	u, err := url.Parse(targetUrl)
	if err != nil{
		return err
	}
	this.cookiejar.SetCookies(u,cookies)
	return nil
}

func (this *HttpClient) Get(url string) (html string,err error){
	req,err:=http.NewRequest("GET", url, nil)
	if err != nil{
		return
	}

	for k,v :=range this.header{
		req.Header.Set(k,v)
	}
	resp,err:=this.client.Do(req)
	if err != nil{
		return
	}
	if resp.StatusCode != 200{
		err = errors.New("status code:%d",resp.StatusCode)
		resp.Body.Close()
		return
	}
	html,err =GetHtml(resp)
	if err != nil{
		resp.Body.Close()
		return
	}
	resp.Body.Close()
	return
}

func (this *HttpClient) Post(url string,data string)(html string,err error){
	req,err:=http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil{
		return
	}
	for k,v :=range this.header{
		req.Header.Set(k,v)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := this.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	html,err =GetHtml(resp)
	if err != nil {
		return
	}
	return
}

func (this *HttpClient) PostJson(url string,data string) (html string,err error){
	req,err:=http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		return
	}
	for k,v :=range this.header{
		req.Header.Set(k,v)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := this.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	html,err =GetHtml(resp)
	if err != nil {
		return
	}
	return
}

func (this *HttpClient) PostBinary(url string,bin []byte)(html string,err error){
	req, err := http.NewRequest("POST", url, bytes.NewReader(bin))
	if err != nil{
		return
	}
	for k,v :=range this.header{
		req.Header.Set(k,v)
	}
	resp, err := this.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	html,err =GetHtml(resp)
	if err != nil {
		return
	}
	return
}
