package httpclient

import (
	"golang.org/x/net/proxy"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type HttpClient struct{
	client http.Client
	header map[string]string
	cookiejar http.CookieJar
}

func NewHttpClient () (c *HttpClient,err error){
	jar,err:=cookiejar.New(nil)
	if err != nil{
		_=jar
		return
	}
	h:=make(map[string]string)
	h["Accept-Encoding"]="zh-CN,zh;q=0.9"
	h["Cache-Control"] = "no-cache"
	h["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
	h["Accept-Encoding"] = "gzip, deflate, br"
	h["User-Agent"] = UserAgents.One()
	c = &HttpClient{
		client:http.Client{Jar:jar},
		header:make(map[string]string),
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

func (this *HttpClient) SetHeaderField(key string,value string){
	this.header[key]=value
}

func (this *HttpClient) SetHttpProxy(proxy string) error{
	proxyUrl, err := url.Parse(proxy)
	if err != nil{
		return err
	}
	this.client.Transport=&http.Transport{Proxy: http.ProxyURL(proxyUrl)}
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

func (this *HttpClient) SetTimeOut(d time.Duration){
	this.client.Timeout=d
}

func (this *HttpClient) UnSetTimeOut(){
	this.client.Timeout=0
}