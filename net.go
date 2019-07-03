package httpclient

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/qianlidongfeng/httpclient/socks"
	"io"
	"net"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (this *Client) dial(u *url.URL) (conn net.Conn,err error){
	var dest string
	var dialInfo *url.URL
	if this.proxy == nil{
		dialInfo=u
		dest = dialInfo.Host
	}else{
		dialInfo=this.proxy
		dest = dialInfo.Host
	}

	if dialInfo.Scheme=="https"{
		if strings.Index(dialInfo.Host,":") == -1{
			dest=dest+":443"
		}
		dialer:=new(net.Dialer)
		dialer.Timeout=this.cfg.timeout
		conn, err = tls.DialWithDialer(dialer, "tcp", dest, &tls.Config{InsecureSkipVerify:true})
	}else if dialInfo.Scheme=="http"||dialInfo.Scheme==""{
		if strings.Index(dialInfo.Host,":") == -1{
			dest=dest+":80"
		}
		conn,err = net.DialTimeout("tcp",dest,this.cfg.timeout)
	}else if dialInfo.Scheme=="socks5" &&this.proxy!=nil{
		deadline:=time.Now().Add(this.cfg.timeout)
		d:=new(net.Dialer)
		d.Deadline=deadline
		conn,err = d.Dial("tcp",dest)
		if err != nil{
			return
		}
		dialer:=socks.NewDialer("tcp",dest)
		dest = u.Host
		if strings.Index(u.Host,":")==-1{
			if u.Scheme=="http"||u.Scheme==""{
				dest=dest+":80"
			}else if u.Scheme == "https"{
				dest=dest+":443"
			}
		}
		ctx,cancel:=context.WithDeadline(context.Background(),deadline)
		_,err=dialer.DialWithConn(ctx,conn,"tcp",dest)
		cancel()
	}else{
		err = errors.New(fmt.Sprintf("uknown protocol %s",dialInfo.Scheme))
	}
	if err != nil{
		return
	}
	return
}


func (this *Client)readHeader(conn io.Reader) (header textproto.MIMEHeader,statusCode int,err error){
	line, err := this.readLine(conn);
	if err != nil {
		return
	}
	if i := strings.IndexByte(line, ' '); i == -1 {
		err=errors.New("malformed HTTP response:"+line)
		return
	} else {
		status:= strings.TrimLeft(line[i+1:], " ")
		if i := strings.IndexByte(status, ' '); i != -1 {
			statusCode,err= strconv.Atoi(status[:i])
			if err !=nil{
				return
			}
		}
	}
	header=textproto.MIMEHeader{}
	for{
		var line string=""
		line,err=this.readLine(conn)
		if err != nil{
			return
		}
		if line=="\r\n"||line=="\n"{
			break
		}
		index:=strings.Index(line,":")
		if index==-1{
			err=errors.New(fmt.Sprintf("malformed HTTP header:%s",line))
			return
		}
		key:=line[:index]
		value:=strings.TrimRight(strings.TrimRight(strings.TrimLeft(line[index+1:]," "),"\n"),"\r")
		header[key]=append(header[key],value)
	}
	return
}

func (this *Client)readLine(conn io.Reader) (string,error){
	var line []byte
	buf:=make([]byte,1)
	for{
		n,err:=io.ReadFull(conn,buf)
		if err !=nil || n != 1{
			return string(line),err
		}
		line=append(line,buf[0])
		if buf[0]=='\n'{
			break
		}
	}
	return string(line),nil
}

func (this *Client)readContent(conn io.Reader,header textproto.MIMEHeader)(content []byte,err error){
	var buffer []byte
	if encoding,ok:=header["Transfer-Encoding"];ok{
		if strings.ToLower(encoding[0]) !="chunked"{
			err=errors.New(fmt.Sprintf("bad Transfer-Encoding:%s",encoding[0]))
			return
		}
		crlf:=make([]byte,2)
		for{
			var header string = ""
			var l int64 = 0
			var n int = 0
			header,err=this.readLine(conn)
			if err != nil {
				return
			}
			l,err=strconv.ParseInt(strings.TrimRight(header,"\r\n"), 16, 64)
			if err != nil{
				return
			}
			if l==0{
				break
			}
			buf := make([]byte,l)
			n,err=io.ReadFull(conn,buf)
			if err != nil||n!=int(l){
				return
			}
			n,err=io.ReadFull(conn,crlf)
			if err != nil||string(crlf)!="\r\n"||n!=2{
				err=errors.New("bad trunck crlf")
				return
			}
			buffer=append(buffer,buf[:]...)
		}
	}else if contentLength,ok:=header["Content-Length"];ok{
		var cl,n int
		cl,err=strconv.Atoi(contentLength[0])
		if err != nil{
			err=errors.New(fmt.Sprintf("bad Content-Length:%s",contentLength[0]))
			return
		}
		buffer=make([]byte,cl)
		n,err=io.ReadFull(conn,buffer)
		if err != nil{
			return
		}else if n != cl{
			err = errors.New("not enough content with Content-Length")
			return
		}
	}else{
		err = errors.New("no Transfer-Encoding or Content-Length in respone header")
		return
	}
	if encoding,ok:=header["Content-Encoding"];ok{
		content,err=DecodeContent(buffer,strings.ToLower(encoding[0]))
		if err !=nil{
			return
		}

	}else{
		content=buffer
	}
	return
}

func (this *Client) addTls(conn net.Conn,u *url.URL) (net.Conn,error){
	serverName:=u.Host
	index := strings.Index(serverName,":")
	if index != -1{
		serverName=serverName[:index]
	}
	tlsConn := tls.Client(conn, &tls.Config{InsecureSkipVerify:true,ServerName:serverName})
	err:=tlsConn.Handshake()
	if err!=nil{
		return nil,err
	}
	return tlsConn,nil
}

func (this *Client)readHttpsConnectHeader(conn io.Reader) (statusCode int,err error){
	r:=bufio.NewReader(conn)
	tp:= textproto.NewReader(r)
	line, err := tp.ReadLine();
	if err != nil {
		return
	}
	if i := strings.IndexByte(line, ' '); i == -1 {
		err=errors.New("malformed HTTP response:"+line)
		return
	} else {
		status:= strings.TrimLeft(line[i+1:], " ")
		if i := strings.IndexByte(status, ' '); i != -1 {
			statusCode,err= strconv.Atoi(status[:i])
			if err !=nil{
				return
			}
		}
	}
	_,err=tp.ReadMIMEHeader()
	if err != nil{
		return
	}
	return
}

func (this *Client)writeHttpsConnectHeader(u *url.URL,conn io.Writer) error{
	ruri:=u.Host
	if strings.Index(ruri,":")== -1{
		if u.Scheme == "http"{
			ruri=ruri+":80"
		}else if u.Scheme == "https"{
			ruri=ruri+":443"
		}
	}
	w := bufio.NewWriter(conn)
	_, err := fmt.Fprintf(w, "%s %s HTTP/1.1\r\n", "CONNECT", ruri)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Host: %s\r\n", ruri)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "User-Agent: %s\r\n", this.header["User-Agent"])
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "\r\n")
	if err != nil {
		return err
	}
	return w.Flush()
}

func (this *Client) writeHeader(u *url.URL,conn io.Writer,method string,data []byte) error{
	ruri:=u.RequestURI()
	if this.proxy != nil{
		if u.Scheme=="http" && this.proxy.Scheme != "socks5"{
			ruri=u.Scheme+"://"+u.Host+ruri
		}
	}
	w := bufio.NewWriter(conn)
	_, err := fmt.Fprintf(w, "%s %s HTTP/1.1\r\n", method, ruri)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Host: %s\r\n", u.Host)
	if err != nil {
		return err
	}

	for k,v :=range this.header{
		if k=="Cookie"{
			continue
		}
		_, err = fmt.Fprintf(w, "%s: %s\r\n", k,v)
		if err != nil {
			return err
		}
	}
	for k,v :=range this.tempheader{
		if k=="Cookie"{
			continue
		}
		_, err = fmt.Fprintf(w, "%s: %s\r\n", k,v)
		if err != nil {
			return err
		}
	}
	var ck string
	if cookie,ok:=this.header["Cookie"];ok{
		ck=cookie
	}
	if cookie,ok:=this.tempheader["Cookie"];ok{
		if ck != ""{
			ck= ck+"; "+cookie
		}else{
			ck=cookie
		}
	}
	if this.jar != nil{
		for _, cookie := range this.jar.Cookies(u) {
			//req.AddCookie(cookie)
			if ck!=""{
				ck=ck+"; "+cookie.String()
			}else{
				ck=cookie.String()
			}
		}
	}
	if ck != ""{
		_, err = fmt.Fprintf(w, "Cookie: %s\r\n",ck)
		if err != nil {
			return err
		}
	}
	this.tempheader=make(map[string]string)
	_, err = io.WriteString(w, "\r\n")
	if err != nil {
		return err
	}
	if method == "POST"{

	}
	return w.Flush()
}

func (this *Client)do(destUrl string,method string,data []byte) (resp Resp,err error){
	u, err := url.Parse(destUrl)
	if err != nil{
		return
	}
	if u.Host == ""{
		err=errors.New("url missing host")
		return
	}
	begin:=time.Now()
	conn,err:=this.dial(u)
	if err != nil{
		if conn != nil{
			conn.Close()
		}
		return
	}
	defer conn.Close()
	var timeout time.Duration
	if this.cfg.timeout != 0{
		timeout=this.cfg.timeout-time.Since(begin)
		if timeout>0{
			conn.SetDeadline(time.Now().Add(timeout))
		}else{
			err = errors.New("i/o timeout")
			return
		}
	}
	if this.proxy != nil{
		if this.proxy.Scheme=="socks5"{
			//do nothing
		}else if u.Scheme=="https"{
			var statusCode int
			err=this.writeHttpsConnectHeader(u,conn)
			if err != nil{
				return
			}
			statusCode,err= this.readHttpsConnectHeader(conn)
			if err != nil{
				return
			}
			if statusCode!=200{
				err = errors.New(fmt.Sprintf("proxy respone status code:%d",statusCode))
				return
			}
		}else if u.Scheme=="http"{
			//do nothing
		}
		if u.Scheme=="https"{
			conn,err=this.addTls(conn,u)
			if err != nil{
				return
			}
		}
	}
	err= this.writeHeader(u,conn,method,data)
	if err != nil{
		return
	}
	header,statusCode,err:=this.readHeader(conn)
	if err != nil{
		return
	}
	if this.jar != nil{
		if cookies:=getCookieFromRespHeader(header); len(cookies)>0{
			this.jar.SetCookies(u,cookies)
		}
	}
	if statusCode>300&&statusCode<307{
		if this.cfg.redirect{
			if location,ok:=header["Location"];ok{
				resp,err=this.do(location[0],method,data)
				return
			}else{
				err = errors.New(fmt.Sprintf("has no location,status code:%d",resp.StatusCode))
				return
			}
		}
	}else if statusCode != 200{
		resp.StatusCode=statusCode
		err = errors.New(fmt.Sprintf("status code:%d",resp.StatusCode))
		return
	}
	content,err:=this.readContent(conn,header)
	if err !=nil{
		return
	}
	resp.StatusCode=statusCode
	resp.Html=string(content)
	return
}