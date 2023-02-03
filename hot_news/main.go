package hot_news

import (
	"github.com/gocolly/colly"
)

type HotNewsItem struct {
	Rank    string `json:"rank,omitempty"`
	Title   string `json:"title,omitempty"`
	Url     string `json:"url,omitempty"`
	ViewNum string `json:"viewNum,omitempty"`
}

func GetContent(url string) ([]HotNewsItem, error) {
	c := colly.NewCollector()
	var list []HotNewsItem
	c.OnHTML("#page > div.c-d.c-d-e > div.Zd-p-Sc > div:nth-child(1) > div.cc-dc-c > div > div.jc-c > table > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			Rank := element.ChildText("td:nth-child(1)")
			Title := element.ChildText("td:nth-child(2) > a")
			Url := element.ChildAttr("td:nth-child(2) > a", "href")
			ViewNum := element.ChildText("td:nth-child(3)")
			h := HotNewsItem{
				Rank,
				Title,
				"https://tophub.today" + Url,
				ViewNum,
			}
			list = append(list, h)
		})
	})
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}
	return list, nil
}
