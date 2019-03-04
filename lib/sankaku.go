package lib

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func GetSankakuPictures(num string, preview bool) ([]string, error) {
	resp, err := http.Get("https://chan.sankakucomplex.com/ja/post/index.content?tags=響（艦これ）&page=" + num)
	if err != nil {
		fmt.Println("notthing in this page")
		return nil, fmt.Errorf("Nothing can find in this page.")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)

	if preview {
		//filter 1
		re := regexp.MustCompile("(<img class=preview src=\").*?(\")")
		unfilteredList := re.FindAllString(content, -1)

		//filter 2
		var imageUrls []string
		for _, val := range unfilteredList {
			re := regexp.MustCompile("src=\"(.*?)\"")
			filtered := re.FindStringSubmatch(val)
			imageUrls = append(imageUrls, "https:"+filtered[1])
		}

		return imageUrls, nil
	} else {
		//filter 1
		re := regexp.MustCompile("(<a href=\").*?(\" onclick)")
		unfilteredList := re.FindAllString(content, -11)

		//filter 2
		var imageUrls []string
		for _, val := range unfilteredList {
			re := regexp.MustCompile("href=\"(.*?)\"")
			filtered := re.FindStringSubmatch(val)
			imageUrls = append(imageUrls, filtered[1])
		}

		//randonNum := randInt(0, len(imageUrls))

		resp, err := http.Get("https://chan.sankakucomplex.com" + imageUrls[4])
		if err != nil {
			fmt.Println(err)
			fmt.Println("notthing in this page")
			return nil, fmt.Errorf("Nothing can find in this page.")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		content := string(body)

		imgUrl := "" //init

		re = regexp.MustCompile("(<a id=image-link class=sample href=\").*?(\">)")
		if len(re.FindAllString(content, -1)) == 0 {
			re = regexp.MustCompile("(?s)(<a id=image-link class=full>).*?(</a>)")
			imagesRawPattern := re.FindAllString(content, -1)[0]
			re = regexp.MustCompile("src=\"(.*?)\"")
			filtered := re.FindStringSubmatch(imagesRawPattern)[1]
			imgUrl = "https:" + filtered
		} else {
			imagesRawPattern := re.FindAllString(content, -1)[0]
			re = regexp.MustCompile("href=\"(.*?)\"")
			filtered := re.FindStringSubmatch(imagesRawPattern)[1]
			imgUrl = "https:" + filtered
		}
		imgUrl, _ = url.QueryUnescape(imgUrl)
		imgUrl = strings.Replace(imgUrl, "&amp;", "&", -1)

		return []string{imgUrl}, nil
	}
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
