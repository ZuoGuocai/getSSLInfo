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
	"net/http"
	"os"
	"strconv"
	//"log"
	//"strings"
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

func main() {

	var myurl = [...]string{"https://youtube.com", "https://google.com", "https://www.hao123.com","https://www.baidu.com"}
 	var tableList [][]string
	for _, v := range myurl {
		oneList:= getSSLInfo(v)
		tableList = append(tableList, oneList)
	}
        genTable(tableList)
}
