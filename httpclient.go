package httpclient

import (
	"errors"
	"golang.org/x/net/proxy"
	"io/ioutil"
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

func NewHttpClient () (c *HttpClient,err error){
	jar,err:=cookiejar.New(nil)
	if err != nil{
		_=jar
		return
	}
	h:=make(map[string]string)
	h["Accept-Encoding"]="zh-CN,zh;q=0.9"
	h["Cache-Control"] = "no-cache"
	h["Connection"] = "keep-alive"
	h["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
	h["Accept-Encoding"] = "gzip, deflate, br"
	h["User-Agent"] = UserAgents.One()
	c = &HttpClient{
		client:http.Client{Jar:jar},
		header:h,
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

func (this *HttpClient) UsetHeaderField(key string){
	delete(this.header,key)
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

func (this *HttpClient) Get(url string) (html string,err error){
	req,err:=http.NewRequest("GET", url, nil)
	if err != nil{
		return
	}
	for k,v :=range this.header{
		req.Header.Add(k,v)
	}
	resp,err:=this.client.Do(req)
	if err != nil{
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		resp.Body.Close()
		return
	}
	html = string(b)
	return
}

func (this *HttpClient) SetCookies(targetUrl string, cookies []*http.Cookie) error{
	u, err := url.Parse(targetUrl)
	if err != nil{
		return err
	}
	this.cookiejar.SetCookies(u,cookies)
	return nil
}

func MakeCookies(domain string,path string,cookie string)(cookies []*http.Cookie,err error){
	cks:=strings.Split(cookie,";")
	for _,v :=range cks{
		n:=strings.Index(v,"=")
		if n == -1{
			err = errors.New("bad cookie")
			return
		}
		s:=0
		if v[0]==' '{
			s=1
		}
		if len(v)<=s+1{
			err = errors.New("bad cookie")
			return
		}
		name:=v[s:n]
		value:=v[n+1:]
		cookies=append(cookies,&http.Cookie{
			Name:name,
			Value:value,
			Path:path,
			Domain:domain,
		})
	}
	return
}