package httpclient

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	//"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Resp struct{
	Html string
	StatusCode int
}


type HttpClient struct{
	client http.Client
	header map[string]string
	tempheader map[string]string
	cookiejar http.CookieJar
	reqclose bool
}

func NewHttpClient () HttpClient{
	jar,_:=cookiejar.New(nil)
	header:=Headers.One()
	header["User-Agent"] = UserAgents.One()
	return HttpClient{
		client:http.Client{
			Jar:nil,
			Transport:&http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				TLSNextProto:    make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
				MaxIdleConns:1000,
				MaxIdleConnsPerHost:1000,
				//MaxConnsPerHost:1000,
				DisableKeepAlives:true,
			},
	},
		header:header,
		tempheader:make(map[string]string),
		cookiejar:jar,
		reqclose:true,
	}
}

func (this *HttpClient) SetHeader(header map[string]string){
	this.header=make(map[string]string)
	for k,v := range header{
		this.header[k]=v
	}
}

func (this *HttpClient) SetTempHeaderField(key string,value string){
	this.tempheader[key]=value
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
	this.client.Transport.(*http.Transport).Proxy=http.ProxyURL(proxyUrl)
	return err
}

func (this *HttpClient) UnsetHttpProxy(){
	this.client.Transport.(*http.Transport).Proxy = nil
	return
}

func (this *HttpClient) SetProxy(proxy string) error{
	proxyUrl, err := url.Parse(proxy)
	if err != nil{
		return err
	}
	this.client.Transport.(*http.Transport).Proxy=http.ProxyURL(proxyUrl)
	return err
}

func (this *HttpClient) SetSock5Proxy(proxy string) error{
	proxyUrl, err := url.Parse(proxy)
	if err != nil{
		return err
	}
	this.client.Transport.(*http.Transport).Proxy=http.ProxyURL(proxyUrl)
	return err
	//dialSocksProxy, err := proxy.SOCKS5("tcp", prox, nil, proxy.Direct)
	/*dialSocksProxy, err := proxy.SOCKS5("tcp", prox, nil,
		&net.Dialer { Timeout:   this.client.Timeout, KeepAlive: this.client.Timeout},
	)
	if err != nil{
		return err
	}
	this.client.Transport.(*http.Transport).Dial=dialSocksProxy.Dial
	return err*/
}

func (this *HttpClient) UnsetSocksProxy() {
	this.client.Transport.(*http.Transport).Dial=nil
}

func (this *HttpClient) UnsetProxy(){
	this.client.Transport.(*http.Transport).Proxy = nil
	this.client.Transport.(*http.Transport).Dial=nil
}

func (this *HttpClient) SetTimeOut(d time.Duration){
	this.client.Timeout=d
	this.client.Transport.(*http.Transport).TLSHandshakeTimeout=d
	this.client.Transport.(*http.Transport).ExpectContinueTimeout=d
	this.client.Transport.(*http.Transport).IdleConnTimeout=d
	this.client.Transport.(*http.Transport).DialContext=(&net.Dialer{
		Timeout:   d,
		KeepAlive: d,
	}).DialContext
}

func (this *HttpClient) CloseIdleConnections(){
	this.client.CloseIdleConnections()
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
	if this.client.Jar !=nil{
		this.client.Jar=jar
	}
	this.cookiejar=jar
	return nil
}


func (this *HttpClient) SetCookies(targetUrl string, cookies []*http.Cookie) error{
	u, err := url.Parse(targetUrl)
	if err != nil{
		return err
	}
	this.cookiejar.SetCookies(u,cookies)
	return nil
}

func (this *HttpClient) GetCooikes(targetUrl string) ([]*http.Cookie,error){
	u, err := url.Parse(targetUrl)
	if err != nil{
		return nil,err
	}
	return this.cookiejar.Cookies(u),nil
}

func (this *HttpClient) Get(url string) (r Resp,err error){
	req,err:=http.NewRequest("GET", url, nil)
	if err != nil{
		return
	}
	for k,v :=range this.header{
		req.Header.Set(k,v)
	}
	for k,v :=range this.tempheader{
		req.Header.Set(k,v)
	}
	if this.reqclose{
		req.Header.Set("Connection","close")
		req.Close=true
	}else{
		req.Header.Set("Connection","keep-alive")
	}
	this.tempheader=make(map[string]string)
	resp,err:=this.client.Do(req)
	if err != nil{
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return
	}
	r.StatusCode=resp.StatusCode
	if resp.StatusCode != 200{
		err = errors.New(fmt.Sprintf("status code:%d",resp.StatusCode))
		resp.Body.Close()
		return
	}
	html,err :=GetHtml(resp)
	if err != nil{
		resp.Body.Close()
		return
	}
	r.Html=html
	resp.Body.Close()
	return
}

func (this *HttpClient) Post(url string,data string)(r Resp,err error){
	req,err:=http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil{
		return
	}
	for k,v :=range this.header{
		req.Header.Set(k,v)
	}
	for k,v :=range this.tempheader{
		req.Header.Set(k,v)
	}
	if this.reqclose{
		req.Header.Set("Connection","close")
		req.Close=true
	}else{
		req.Header.Set("Connection","keep-alive")
	}
	this.tempheader=make(map[string]string)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := this.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	r.StatusCode=resp.StatusCode
	html,err :=GetHtml(resp)
	if err != nil {
		return
	}
	r.Html=html
	return
}

func (this *HttpClient) PostJson(url string,data string) (r Resp,err error){
	req,err:=http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		return
	}
	for k,v :=range this.header{
		req.Header.Set(k,v)
	}
	for k,v :=range this.tempheader{
		req.Header.Set(k,v)
	}
	if this.reqclose{
		req.Header.Set("Connection","close")
		req.Close=true
	}else{
		req.Header.Set("Connection","keep-alive")
	}
	this.tempheader=make(map[string]string)
	req.Header.Set("Content-Type", "application/json")
	resp, err := this.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	r.StatusCode=resp.StatusCode
	html,err :=GetHtml(resp)
	if err != nil {
		return
	}
	r.Html=html
	return
}

func (this *HttpClient) PostBinary(url string,bin []byte)(r Resp,err error){
	req, err := http.NewRequest("POST", url, bytes.NewReader(bin))
	if err != nil{
		return
	}
	for k,v :=range this.header{
		req.Header.Set(k,v)
	}
	for k,v :=range this.tempheader{
		req.Header.Set(k,v)
	}
	if this.reqclose{
		req.Header.Set("Connection","close")
		req.Close=true
	}else{
		req.Header.Set("Connection","keep-alive")
	}
	this.tempheader=make(map[string]string)
	resp, err := this.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	r.StatusCode=resp.StatusCode
	html,err :=GetHtml(resp)
	if err != nil {
		return
	}
	r.Html=html
	return
}

func (this *HttpClient) SetReqClose(b bool){
	this.reqclose=b
}