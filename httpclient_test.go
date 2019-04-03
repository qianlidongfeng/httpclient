package httpclient_test

import (
	"fmt"
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