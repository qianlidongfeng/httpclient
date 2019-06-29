package httpclient

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type clientConfig struct{
	redirect bool
	timeout time.Duration
}

type Client struct{
	proxy *url.URL
	header map[string]string
	tempheader map[string]string
	cfg clientConfig
	jar http.CookieJar
	cookieJar http.CookieJar
}

func NewClient() Client{
	cj,_:=cookiejar.New(nil)
	return Client{
		header:Headers.One(),
		tempheader:make(map[string]string),
		cfg:clientConfig{
			redirect:true,
		},
		jar:nil,
		cookieJar:cj,
	}
}

func (this *Client)Get(destUrl string) (resp Resp,err error){
	return this.do(destUrl,"GET",nil)
}


func (this *Client) SetProxy(proxy string) error{
	var err error
	this.proxy,err=url.Parse(proxy)
	if err != nil{
		return err
	}
	if this.proxy.Host == ""{
		return errors.New("proxy url missing host")
	}
	return nil
}

func (this *Client) UnSetProxy(){
	this.proxy=nil
}

func (this *Client) SetTimeOut(timeOut time.Duration){
	this.cfg.timeout=timeOut
}

func (this *Client) EnableRedirect(){
	this.cfg.redirect=true
}

func (this *Client) DisableRedirect(){
	this.cfg.redirect=false
}

func (this *Client) EnableCookie(){
	this.jar=this.cookieJar
}

func (this *Client) DisableCookie(){
	this.jar=nil
}

func (this *Client) SetCookies(targetUrl string, cookies []*http.Cookie) error{
	u, err := url.Parse(targetUrl)
	if err != nil{
		return err
	}
	this.cookieJar.SetCookies(u,cookies)
	return nil
}

func (this *Client) GetCooikes(targetUrl string) ([]*http.Cookie,error){
	u, err := url.Parse(targetUrl)
	if err != nil{
		return nil,err
	}
	return this.cookieJar.Cookies(u),nil
}

func (this *Client) ClearCookie() error{
	jar,err:=cookiejar.New(nil)
	if err != nil{
		return err
	}
	if this.jar!= nil{
		this.jar=jar
	}
	this.cookieJar=jar
	return nil
}

func (this *Client) SetHeaderField(key string,value string){
	this.header[key]=value
}

func (this *Client) UnsetHeaderField(key string){
	delete(this.header,key)
}

func (this *Client) SetHeader(header map[string]string){
	this.header=make(map[string]string)
	for k,v := range header{
		this.header[k]=v
	}
}

func (this *Client) UnsetHeader(){
	this.header=make(map[string]string)
}

func (this *Client) SetTempHeaderField(key string,value string){
	this.tempheader[key]=value
}