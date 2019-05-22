package httpclient

import (
	"math/rand"
	"time"
)

type HEADER struct{
	AcceptLanguage string
	CacheControl string
	Connection string
	Accept string
	AcceptEncoding string
}

type HEADERS []HEADER

var Headers = HEADERS{
	{
		AcceptLanguage:"zh-CN,zh;q=0.9",
		CacheControl:"no-cache",
		//Connection:"keep-alive",
		//Connection:"close",
		Accept:"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
		AcceptEncoding:"gzip, deflate, br",
	},
	{
		AcceptLanguage:"zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		CacheControl:"no-cache",
		//Connection:"keep-alive",
		//Connection:"close",
		Accept:"text/html,application/xhtml+xm…plication/xml;q=0.9,*/*;q=0.8",
		AcceptEncoding:"gzip, deflate, br",
	},
}

func (this *HEADERS) One() map[string]string{
	l := len(*this)
	rand.Seed(time.Now().UnixNano())
	header:=(*this)[rand.Intn(l-1)]
	m:=make(map[string]string)
	if header.AcceptLanguage!= ""{
		m["Accept-Language"]=header.AcceptLanguage
	}
	if header.CacheControl != ""{
		m["Cache-Control"]=header.CacheControl
	}
	if header.Connection != ""{
		m["Connection"]=header.Connection
	}
	if header.Accept != ""{
		m["Accept"]=header.Accept
	}
	if header.AcceptEncoding != ""{
		m["Accept-Encoding"]=header.AcceptEncoding
	}
	return m
}