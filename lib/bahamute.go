package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
	/today
	/day

	guessDayFunction (if page 1 dont have this date, then guess!)
*/

func RequestImageByDate(date string) ([]string, error) {
	fmt.Println(date)
	setDate := date
	c := &Class{}
	c.GetList(CalcPagePosition(setDate))

	link, ifmatch := c.GetLinkBySearchList(setDate)
	if !ifmatch {
		fmt.Println("sorry i can't find any thing.")
		return nil, fmt.Errorf("can't find this date")
	}

	images, err := GetPostPictures(link)
	if err != nil {
		fmt.Println("can't make a request")
		return nil, fmt.Errorf("find post failed")
	}

	return images, nil
}

type Class struct {
	queryString string
}

func (c *Class) GetList(page int) {

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "ckFORUM_sval",
		Value:  "%E6%AF%8F%E6%97%A5%E9%9F%BF",
		Path:   "/",
		Domain: "forum.gamer.com.tw",
	}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{
		Name:   "ckFORUM_stype",
		Value:  "title",
		Path:   "/",
		Domain: "forum.gamer.com.tw",
	}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{
		Name:   "ckFORUM_searchType",
		Value:  "baha",
		Path:   "/",
		Domain: "forum.gamer.com.tw",
	}
	cookies = append(cookies, cookie)

	u, _ := url.Parse("https://forum.gamer.com.tw/B.php?page=" + strconv.Itoa(page) + "&bsn=60076&forumSearchQuery=%E6%AF%8F%E6%97%A5%E9%9F%BF&subbsn=0")
	jar.SetCookies(u, cookies)
	client := &http.Client{
		Jar: jar,
	}

	postData := url.Values{}
	req, _ := http.NewRequest("GET", "https://forum.gamer.com.tw/B.php?page="+strconv.Itoa(page)+"&bsn=60076&forumSearchQuery=%E6%AF%8F%E6%97%A5%E9%9F%BF", strings.NewReader(postData.Encode()))
	req.Header.Add("Content-Type", "text/html")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request Failed!")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	c.queryString = string(body)
}

func dateConventor(date string) string {
	s := strings.Split(date, "/")

	converted := ""

	//remove zero from number head of month field.
	if len(s[0]) > 1 && string(s[0][0]) == "0" {
		converted += string(s[0][1]) + "/"
	} else {
		converted += string(s[0] + "/")
	}

	//remove zero from number head of date field.
	if len(s[1]) > 1 && string(s[1][0]) == "0" {
		converted += string(s[1][1])
	} else {
		converted += string(s[1])
	}

	return converted
}

func (c *Class) GetLinkBySearchList(date string) (string, bool) {
	//filter 1
	re := regexp.MustCompile("<a data-gtm=\"B頁文章列表\"(.*?)</a>")
	unfilteredList := re.FindAllString(c.queryString, -1)

	//filter 2
	findStr := ""
	for _, val := range unfilteredList {
		re = regexp.MustCompile("【艦隊】" + dateConventor(date))
		if re.MatchString(val) {
			findStr = val
		}
	}

	//filter 3
	re = regexp.MustCompile("href=\"(.*?)\"")
	matchdata := re.FindStringSubmatch(findStr)
	if len(matchdata) < 1 {
		//if not match anything
		return "", false
	}
	link := "https://forum.gamer.com.tw/" + matchdata[1]
	return link, true
}

func GetPostPictures(link string) ([]string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return []string{}, fmt.Errorf("can't make a request.")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)

	//fmt.Println(content)
	//make filter to this page
	//filter 1 (find main post)
	re := regexp.MustCompile("<div class=\"c-post__body\">((.|\n|\r)*?)</div>")
	mainPost := re.FindAllStringSubmatch(content, -1)[0][0]

	//filter 2 (find image url)
	re = regexp.MustCompile("<a name=\"attachImgName\" href=\"(.*?)\"")
	imageUrlsRaw := re.FindAllStringSubmatch(mainPost, -1)

	var imageUrls []string
	for _, val := range imageUrlsRaw {
		imageUrls = append(imageUrls, val[1])
	}
	return imageUrls, nil
}

func CalcPagePosition(date string) int {

	dateall := strings.Split(date, "/")
	setyear, _ := strconv.Atoi(string(dateall[0]))
	setdate, _ := strconv.Atoi(string(dateall[1]))

	datediff := diffDays(time.Month(setyear), setdate)

	//one page (30 raw)
	return (datediff / daysIn(time.Now().Month(), time.Now().Year())) + 1
}

func diffDays(month1 time.Month, day1 int) int {
	now := time.Now()

	d2 := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, time.UTC,
	)

	d1 := time.Date(
		now.Year(), month1, day1,
		0, 0, 0, 0, time.UTC,
	)
	return int(d2.Sub(d1) / (24 * time.Hour))
}

func daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
