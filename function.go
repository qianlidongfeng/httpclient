package httpclient

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"gopkg.in/kothar/brotli-go.v0/dec"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetHtml(resp *http.Response) (html string,err error){
	var b []byte
	switch strings.ToLower(resp.Header.Get("Content-Encoding")) {
	case "gzip":
		var reader io.Reader
		reader,err=gzip.NewReader(resp.Body)
		if err != nil{
			return
		}
		b,err=ioutil.ReadAll(reader)
		if err != nil{
			return
		}
	case "br":
		var br []byte
		br,err=ioutil.ReadAll(resp.Body)
		if err != nil{
			return
		}
		b, err = dec.DecompressBuffer(br, nil)
	default:
		b,err=ioutil.ReadAll(resp.Body)
		if err != nil{
			return
		}
	}
	html=string(b)
	return
}

func DecodeContent(buf []byte,encode string)(content []byte,err error){
	switch encode{
	case "gzip":
		var reader io.Reader
		reader, err = gzip.NewReader(bytes.NewReader(buf))
		if err != nil{
			return
		}
		content,err=ioutil.ReadAll(reader)
		if err != nil{
			return
		}
	case "br":
		content, err = dec.DecompressBuffer(buf, nil)
		if err !=nil{
			return
		}
	default:
		err = errors.New(fmt.Sprintf("unknown Content-Encoding:%s",encode))
		return
	}
	return
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

func GetCookieString(cookies []*http.Cookie) string{
	var s string
	for _,cookie := range cookies{
		s+=cookie.Name+"="+cookie.Value+"; "
	}
	return strings.TrimRight(s,"; ")
}

func GetBaseDomain(url string) string{
	start:= strings.Index(url,".")
	if start != -1{
		url=url[start+1:]
	}
	end := strings.Index(url,"/")
	if end != -1{
		url=url[:end]
	}
	return "."+url
}
