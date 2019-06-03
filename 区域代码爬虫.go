package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/* 爬取 省 市 县 镇 村 数据
 *	url：http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/
 */

var file *os.File


func main() {



	fmt.Println("开始爬取数据")

	file, _ = os.Create("text.txt")

	CrawlProvince(`http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/`)

	// CrawlMunicipal("http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/13.html","","")
}

// 根据 url 获取 页面
func getPage(url string)(pageStr string,err error){
	// 更具链接获取相应
	/*resp, err1 := http.Get(url)
	count := 0
	if err1 != nil {
		count++
		err = fmt.Errorf("第 "+strconv.Itoa(count)+" 次获取网页失败，检查 url 或 网络")
		return
	}*/
	var resp *http.Response
	for i:=1;i<=3;i++{
		var err1 error
		resp, err1 = http.Get(url)
		if err1 == nil{
			break
		}else { // 有错误
			if i !=3 {
				print("第 "+strconv.Itoa(i)+" 次获取网页失败"+url)
			}else {
				err = fmt.Errorf("获取网页失败"+url)
				return
			}
		}
	}


	buff := make([]byte, 1024)

	// 读取页面
	for {
		num, err2 := resp.Body.Read(buff)
		if err2 != nil && err2 != io.EOF{
			err = fmt.Errorf("读取相应体失败")
			return
		}
		if num == 0 {
			break
		}

		// .................... 需转码(GBK->utf-8) .................

		pageStr += string(buff[:num])
	}
	return
}

// 获取所有省份数据页面
func CrawlProvince(url string) {
	buffStr,err1 := getPage(url)
	if err1 != nil{
		panic("获取网页失败 url="+url)
	}

	if len(buffStr) < 100{
		fmt.Println(buffStr)
		panic("服务器开启反扒虫")
	}

	// 创建文件准备存储
	/*file, _ := os.Create("province.txt")
	defer file.Close()*/

	// 解析网页
	reg := regexp.MustCompile(`<td><a href='.{1,6}html'>.{1,30}<br/></a></td>`)
	allString := reg.FindAllString(buffStr, -1)
	for _,value := range allString{
		child := strings.Replace(value,"<td><a href='","",1)
		child = strings.Replace(child,"'>",",",1)
		child = strings.Replace(child,"<br/></a></td>","",1)

		split := strings.Split(child, ",")

		curl := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/"+split[0]
		sid := strings.Replace(split[0],".html","",1)+"0000000000"
		name := split[1]

		//file.Write([]byte(sid+" "+name+"\r\n"))
		fmt.Println(curl,sid,name)

		file.Write([]byte(sid+" "+name+"\r\n"))

		// 查询子级
		CrawlMunicipal(curl,sid,name)
	}

	// 获取省份相应网页，代码，全名 。。。。。。。。。

	/*for {
		CrawlMunicipal("","","");
	}*/

	//fmt.Println(buffStr)

}

// 获取 市级数据
func CrawlMunicipal(url, fatherSid, fahterName string)(err error) {
	pageStr, err1 := getPage(url)
	if err1 != nil{
		panic("获取网页失败 url="+url)
	}

	//fmt.Println(pageStr)


	// 解析网页
	reg := regexp.MustCompile(`<td><a href='.{5,20}'>.{5,18}</a></td><td><a href='.{5,20}'>.{1,50}</a></td>`)
	allString := reg.FindAllString(pageStr, -1)

	// 便利结果
	for _,value := range allString{

		child := strings.Replace(value,"<td><a href='","",1)
		child = strings.Replace(child,"'>",",",2)
		child = strings.Replace(child,"</a></td><td><a href='",",",1)
		child = strings.Replace(child,"</a></td>","",1)

		split := strings.Split(child, ",")

		sid := split[1]
		name := split[3]
		curl := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/"+split[0]

		fmt.Println("\t市级",curl,sid,name)

		file.Write([]byte(sid+" "+name+"\r\n"))

		CrawlCounty(curl,sid,name)


	}


	return
}

// 获取 县级数据
func CrawlCounty(url, fatherSid, fahterName string)( err error) {
	pageStr, err1 := getPage(url)
	if err1 != nil{
		panic("获取网页失败 url="+url)
	}

	//fmt.Println(pageStr)


	// 解析网页
	reg := regexp.MustCompile(`<td><a href='.{5,20}'>.{5,18}</a></td><td><a href='.{5,20}'>.{1,50}</a></td>`)
	allString := reg.FindAllString(pageStr, -1)

	// 便利结果
	for _,value := range allString{

		child := strings.Replace(value,"<td><a href='","",1)
		child = strings.Replace(child,"'>",",",2)
		child = strings.Replace(child,"</a></td><td><a href='",",",1)
		child = strings.Replace(child,"</a></td>","",1)

		split := strings.Split(child, ",")



		curl := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/"+string(fatherSid[0:2])+"/"+split[0]
		sid := split[1]
		name := split[3]

		fmt.Println("\t\t县级",curl,sid,name)

		file.Write([]byte(sid+" "+name+"\r\n"))

		CrawlTown(curl,sid,name)


	}


	return
}

// 获取 镇级数据
func CrawlTown(url, fatherSid, fahterName string)( err error) {
	pageStr, err1 := getPage(url)
	if err1 != nil{
		panic("获取网页失败 url="+url)
	}

	// 解析网页
	reg := regexp.MustCompile(`<td><a href='.{5,20}'>.{5,18}</a></td><td><a href='.{5,20}'>.{1,50}</a></td>`)
	allString := reg.FindAllString(pageStr, -1)

	// 便利结果
	for _,value := range allString{

		child := strings.Replace(value,"<td><a href='","",1)
		child = strings.Replace(child,"'>",",",2)
		child = strings.Replace(child,"</a></td><td><a href='",",",1)
		child = strings.Replace(child,"</a></td>","",1)

		split := strings.Split(child, ",")



		curl := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/"+string(fatherSid[0:2])+"/"+string(fatherSid[2:4])+"/"+split[0]
		sid := split[1]
		name := split[3]

		fmt.Println("\t\t\t镇级",curl,sid,name)

		file.Write([]byte(sid+" "+name+"\r\n"))

		CrawlVillage(curl,sid,name)


	}
	return
}

// 获取 村级数据
func CrawlVillage(url, fatherSid, fahterName string)( err error) {
	pageStr, err1 := getPage(url)
	if err1 != nil{
		panic("获取网页失败 url="+url)
	}

	// 解析网页
	// <td>110101001001</td><td>111</td><td>多福巷社区居委会</td></tr>
	reg := regexp.MustCompile(`<tr class='villagetr'><td>.{5,18}</td><td>.{0,12}</td><td>.{1,60}</td></tr>`)
	allString := reg.FindAllString(pageStr, -1)

	// 便利结果
	for _,value := range allString{
		//fmt.Println(value)

		child := strings.Replace(value,"<tr class='villagetr'><td>","",1)
		child = strings.Replace(child,"</td><td>",",",1)
		child = strings.Replace(child,"</td><td>",",",1)
		child = strings.Replace(child,"</td></tr>","",1)

		split := strings.Split(child, ",")

		sid := split[0]
		name := split[2]

		file.Write([]byte(sid+" "+name+"\r\n"))

		fmt.Println("\t\t\t\t村级","____________________________________________",sid,name)

	}
	return
}
