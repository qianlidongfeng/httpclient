package httpclient_test

import (
	"fmt"
	"github.com/qianlidongfeng/string"
	"testing"
	"github.com/qianlidongfeng/httpclient"
	"time"

)

func TestUA_One(t *testing.T) {
	httpclient.NewHttpClient()

	for{
		fmt.Println(httpclient.UserAgents.One())
		time.Sleep(time.Second)
	}
}

func TestHttpClient_Get(t *testing.T) {
	c,err:=httpclient.NewHttpClient()
	if err != nil{
		t.Error(err)
	}
	_=c
}

func TestMakeCookies(t *testing.T) {
	cookies, err := httpclient.MakeCookies(".baidu.com", "/","BAIDUID=B494D4091CAD870FA24C363A40031693:FG=1; BIDUPSID=B494D4091CAD870FA24C363A40031693; PSTM=1551005727; BD_UPN=12314753; BDUSS=Uw0ZVVkTGcwejZyby1jOFZvMHc1WnJzZ35kRFVGN25KVC1wZ0x4N3h4MWtVYUJjQVFBQUFBJCQAAAAAAAAAAAEAAAA4vbUdQW5nbGVfU2VhbgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGTEeFxkxHhcW; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; H_PS_PSSID=1460_21109_28769_28723_28558_28585_26350_28519_22158; delPer=0; BD_CK_SAM=1; PSINO=7; ZD_ENTRY=empty; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; H_PS_645EC=b01fpy%2BKAwQmp6XK7UQBerJCEUKCZTQu754OEftN9DFuxlQRLpWQfQJxmMGRS9A45zLm; BD_HOME=1")
	if err != nil{
		t.Error(err)
	}
	_=cookies
}

func TestHttpClient_SetCookies(t *testing.T) {
	c,err:=httpclient.NewHttpClient()
	if err != nil{
		t.Error(err)
	}
	cookies,err:=httpclient.MakeCookies(".baidu.com","/",
		`BAIDUID=312D691EA3216936E879702AA2C25F38:FG=1; BIDUPSID=312D691EA3216936E879702AA2C25F38; PSTM=1554303051; delPer=0; H_PS_PSSID=1447_21121_18559_28775_28721_28557_28585_26350_28519_28605; BDUSS=Gx1V0lSSlJaV2NQcDJ0dEV0NThTZ2N6eUZydzl3UnNqOXI2cklZWWg3UmhWY3hjRVFBQUFBJCQAAAAAAAAAAAEAAAA4vbUdQW5nbGVfU2VhbgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGHIpFxhyKRcT; TIEBA_USERTYPE=e4dea1ba8138b15b0378a7db; STOKEN=5624e16ecddeadb7315dbd45520953e5d4725a08608d6a26fd74e82868d232c8; TIEBAUID=16f2e80d1c8298db6c89402c; wise_device=0; Hm_lvt_98b9d8c2fd6608d564bf2ac2ae642948=1554303491,1554303893,1554303913,1554303943; Hm_lpvt_98b9d8c2fd6608d564bf2ac2ae642948=1554303943`,
	)
	if err !=nil{
		t.Error(err)
	}
	c.SetCookies("https://www.baidu.com/",cookies)
	html,err:=c.Get("https://www.qq.com/")
	fmt.Println(gstring.GbkToUtf8([]byte(html)))
}

func TestHttpClient_SetSocksProxy(t *testing.T) {
	c,err:=httpclient.NewHttpClient()
	if err != nil{
		fmt.Print(err)
	}
	c.SetSocksProxy("127.0.0.1:1080")
	html,err:=c.Get("https://ip.cn/")
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(html)
}