package lib

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func GetNhentaiBooks() *nHentai {
	var covers, links, titles, length = GetNhentaiList()
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)
	randomSelect := r.Intn(length - 1)

	book := &nHentai{
		Cover: covers[randomSelect],
		Link:  links[randomSelect],
		Title: titles[randomSelect],
	}

	book.Content = book.nHentaiContents()
	return book
}

type nHentai struct {
	Cover   string
	Link    string
	Title   string
	Content []string
}

//send cover and link
func GetNhentaiList() (covers []string, links []string, titles []string, length int) {
	resp, _ := http.Get("https://nhentai.net/character/hibiki/")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)

	//filter 1
	re := regexp.MustCompile("<div class=\"gallery\"(.*?)</div>")
	unfiltered := re.FindAllString(content, -1)

	//set covers
	re = regexp.MustCompile("data-src=\"(.*?)\"")
	for _, val := range unfiltered {
		d := re.FindStringSubmatch(val)[1]
		covers = append(covers, d)
	}

	//set links
	re = regexp.MustCompile("<a href=\"(.*?)\"")
	for _, val := range unfiltered {
		d := re.FindStringSubmatch(val)[1]
		links = append(links, "https://nhentai.net"+d)
	}

	//set titles
	re = regexp.MustCompile("<div class=\"caption\">(.*?)</div>")
	for _, val := range unfiltered {
		d := re.FindStringSubmatch(val)[1]
		titles = append(titles, d)
	}

	return covers, links, titles, len(covers)
}

func (n *nHentai) nHentaiContents() []string {
	resp, _ := http.Get(n.Link)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)

	//filter (remove head and end)
	re := regexp.MustCompile("data-src=\"(.*?)\" src")

	thumblinkraw := []string{}
	for _, val := range re.FindAllString(content, -1) {
		rawlink := re.FindStringSubmatch(val)[1]

		link := strings.Replace(rawlink, "t.nhentai.net", "i.nhentai.net", -1)
		splink := strings.Split(link, "/")

		if strings.Index(splink[5], "cover") == -1 && strings.Index(splink[5], "thumb") == -1 {
			changeToRawImage := strings.Split(splink[5], "t")
			rawImageLink := splink[0] + "/" + splink[1] + "/" + splink[2] + "/" + splink[3] + "/" + splink[4] + "/" + changeToRawImage[0] + changeToRawImage[1]
			thumblinkraw = append(thumblinkraw, rawImageLink)
		}
	}
	return thumblinkraw
}
