package nba

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

type TodaySchedule struct {
	MatchType     string `json:"matchType"`
	Mid           string `json:"mid"`
	LeftId        string `json:"leftId"`
	LeftName      string `json:"leftName"`
	LeftBadge     string `json:"leftBadge"`
	LeftGoal      string `json:"leftGoal"`
	LeftHasUrl    string `json:"leftHasUrl"`
	RightId       string `json:"rightId"`
	RightName     string `json:"rightName"`
	RightBadge    string `json:"rightBadge"`
	RightGoal     string `json:"rightGoal"`
	RightHasUrl   string `json:"rightHasUrl"`
	MatchDesc     string `json:"matchDesc"`
	StartTime     string `json:"startTime"`
	Title         string `json:"title"`
	Logo          string `json:"logo"`
	MatchPeriod   string `json:"matchPeriod"` // 0 未开始 1：进行中 2： 结束
	LivePeriod    string `json:"livePeriod"`
	Quarter       string `json:"quarter"`
	QuarterTime   string `json:"quarterTime"`
	LiveType      string `json:"liveType"`
	LiveId        string `json:"liveId"`
	ProgramId     string `json:"programId"`
	IsPay         string `json:"isPay"`
	GroupName     string `json:"groupName"`
	CompetitionId string `json:"competitionId"`
	TvLiveId      string `json:"tvLiveId"`
	IfHasPlayback string `json:"ifHasPlayback"`
	Url           string `json:"url"`
	CategoryId    string `json:"categoryId"`
	ScheduleId    string `json:"scheduleId"`
	RoseNewsId    string `json:"roseNewsId"`
	WebUrl        string `json:"webUrl"`
	IconType      string `json:"iconType"`
	IsFree        string `json:"isFree"`
	LatestNews    string `json:"latestNews"`
	Week          string `json:"week"`
	LeftEnName    string `json:"leftEnName"`
	RightEnName   string `json:"rightEnName"`
}

type ResultItem struct {
	LeftName    string
	RightName   string
	Period      string
	LeftGoal    string
	RightGoal   string
	DataUrl     string
	MatchPeriod string
	WebUrl      string
	Id          string
}

var typeMap = map[string]string{
	"0": "未开始",
	"1": "进行中",
	"2": "已结束",
}

func GetContent() []ResultItem {
	date := time.Now().Format("2006-01-02")
	client := resty.New()
	resp, _ := client.R().Get(fmt.Sprintf("https://matchweb.sports.qq.com/kbs/list?from=NBA_PC&columnId=100000&startTime=%s&endTime=%s", date, date))

	data := gjson.Get(string(resp.Body()), "data")
	today := gjson.Get(data.String(), date)
	var list []TodaySchedule
	err := json.Unmarshal([]byte(today.String()), &list)
	if err != nil {
		return nil
	}
	var results []ResultItem
	for _, item := range list {
		if item.MatchType == "2" || item.MatchType == "3" {
			if len(strings.Split(item.Mid, ":")) <= 0 {
				continue
			}
			code := strings.Split(item.Mid, ":")[1]
			d := "https://sports.qq.com/nba-stats/nbascore?mid=" + code
			q := ""
			if item.MatchPeriod == "1" {
				q = q + item.Quarter
				if item.QuarterTime != "" {
					q = q + fmt.Sprintf("(%s)", item.QuarterTime)
				}
			} else {
				q = typeMap[item.MatchPeriod]
			}
			results = append(results, ResultItem{
				LeftName:    item.LeftName,
				RightName:   item.RightName,
				Period:      q,
				LeftGoal:    item.LeftGoal,
				RightGoal:   item.RightGoal,
				DataUrl:     d,
				MatchPeriod: item.MatchPeriod,
				WebUrl:      item.WebUrl,
				Id:          item.Mid,
			})
		}
	}
	return results
}
