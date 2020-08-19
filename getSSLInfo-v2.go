/* 
  Author: zuoguocai@126.com
Function: getSSLInfo
 version: 1.0
*/
package main

import (
	"crypto/tls"
	"fmt"
	"github.com/olekukonko/tablewriter"
    "github.com/go-gomail/gomail"
	"net/http"
	"os"
	"strconv"
	//"log"
	"strings"
	"math"
	"time"
)

//var messages chan []string = make(chan []string, 1000)

func getSSLInfo(url string) (data []string) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	seedUrl := url
	resp, err := client.Get(seedUrl)
	defer resp.Body.Close()

	if err != nil {
		fmt.Errorf(seedUrl, " 请求失败")
		panic(err)
	}

	//fmt.Println(resp.TLS.PeerCertificates[0])
	certInfo := resp.TLS.PeerCertificates[0]
	NotBefore := certInfo.NotBefore.Format("2006-01-02 15:04:05")
	NotAfter := certInfo.NotAfter.Format("2006-01-02 15:04:05")
	bb := certInfo.NotAfter
	now := time.Now()
	nn := math.Ceil(bb.Sub(now).Hours() / 24)
	mm := strconv.FormatFloat(nn, 'f', 2, 64)
	//DNSNames := strings.Join(certInfo.DNSNames,"")
	//Subject := certInfo.Subject.String()
	//Issuer := certInfo.Issuer.String()

	//data := [][]string{
	//	[]string{"https", seedUrl, NotBefore, NotAfter, mm},
	//}
	data = []string{"https", seedUrl, NotBefore, NotAfter, mm}
        return
}


func genTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"监控类型", "监控URL", "证书开始时间", "证书结束时间", "证书还有n天到期"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output

}

func genHtml(tableList [][]string) (html string) {

head := `<h2 style="color:red ;display:inline;text-align: center; ">https证书过期时间一览</h2>
<table align="center" border="1" cellpadding="1" cellspacing="1" width="100%" style="border-collapse: collapse;border:1px #9cf solid;">
<tr>
<td style="border:1px #fff solid;background:#9cf;padding: 1px 0 1px 0; color: #fff; font-size: 14px; font-weight: bold; font-family: Arial, sans-serif;height:28px;text-align:center">监控类型</td>
<td style="border:1px #fff solid;background:#9cf;padding: 1px 0 1px 0; color: #fff; font-size: 14px; font-weight: bold; font-family: Arial, sans-serif;height:28px;text-align:center">域名</td>
<td style="border:1px #fff solid;background:#9cf;padding: 1px 0 1px 0; color: #fff; font-size: 14px; font-weight: bold; font-family: Arial, sans-serif;height:28px;text-align:center">开始时间</td>
<td style="border:1px #fff solid;background:#9cf;padding: 1px 0 1px 0; color: #fff; font-size: 14px; font-weight: bold; font-family: Arial, sans-serif;height:28px;text-align:center">到期时间</td>
<td style="border:1px #fff solid;background:#9cf;padding: 1px 0 1px 0; color: #fff; font-size: 14px; font-weight: bold; font-family: Arial, sans-serif;height:28px;text-align:center">还有n天到期</td>
</tr>`


      line1 :=`<td style="border:1px #9cf solid;background:#fff;padding: 1px 0 1px 0; color: #153643; font-size: 12px; font-family: Arial, sans-serif;height:28px;text-align:center">`
      line2 :=`</td>`
      var html_1 string
      for i:=0;i<len(tableList);i++ {
            for ii:=0;ii<len(tableList[i]);ii++{
             //      fmt.Println(line1,tableList[i][ii],line2)
                   line3 := strings.Join([]string{line1,tableList[i][ii]},"")
                  //fmt.Println(line3)
                  html_1 = strings.Join([]string{html_1,line3,line2},"")
            }
            html_1 = strings.Join([]string{html_1,"</tr>"},"")

      }
html = strings.Join([]string{head,html_1,"</table>"},"")
return

}

func  sendMail(html_content string){

	m := gomail.NewMessage()
	m.SetHeader("From", "Admin@zuoguocai.com.cn")
	m.SetHeader("To", "zuoguocai@zuoguocai.com.cn","zuoguocai@126.com")
	// m.SetAddressHeader("Cc", "zuoguocai@zuoguocai.com.cn", "Dan") //抄送
	m.SetHeader("Subject", "https证书过期时间监测") // 邮件标题
	m.SetBody("text/html", html_content) // 邮件内容
	// m.Attach("/home/aa.jpg") //附件

	d := gomail.NewDialer("relay.zuoguocai.com.cn", 25, "Admin@zuoguocai.com.cn", "0HUxreN^LX!y")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
        panic(err)
    }

}

func main() {

	var myUrl = [...]string{"https://youtube.com", "https://google.com", "https://www.hao123.com","https://www.baidu.com"}
 	var tableList [][]string
	for _, v := range myUrl {
	        oneList:= getSSLInfo(v)
		tableList = append(tableList, oneList)
	}
        genTable(tableList)
        aaa := genHtml(tableList)
        sendMail(aaa)
}
