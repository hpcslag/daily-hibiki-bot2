package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetYanderePictures(num string) ([]string, error) {
	resp, err := http.Get("https://yande.re/post?page=" + num + "&tags=hibiki_%28kancolle%29")
	if err != nil {
		fmt.Println("notthing in this page")
		return nil, fmt.Errorf("Nothing can find in this page.")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)

	//filter 1
	re := regexp.MustCompile("<a class=\"directlink(.*?)\">")
	unfilteredList := re.FindAllString(content, -1)

	//filter 2
	var imageUrls []string
	for _, val := range unfilteredList {
		re := regexp.MustCompile("href=\"(.*?)\"")
		filtered := re.FindStringSubmatch(val)
		imageUrls = append(imageUrls, filtered[1])
	}

	return imageUrls, nil
}
