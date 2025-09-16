package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

func input() string {
	fmt.Println("请输入完整cURL，结束后输入EOF（大写，单独一行）:")
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "EOF" {
			break
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "读取输入出错:", err)
		return ""
	}

	return strings.Join(lines, "\n")
}

type KBListResponse struct {
	KBList []struct {
		Title   string `json:"kcmc"`
		Teacher string `json:"xm"`
		Time    string `json:"jcs"`
		Week    string `json:"zcd"`
		Day     string `json:"xqj"`
		Place   string `json:"cdmc"`
		PS      string `json:"xkbz"`
	} `json:"kbList"`
}

func getKBList(curlLine string) KBListResponse {
	var (
		urlStr    string
		method    = "GET"
		payload   string
		cookieStr string
		headers   = make(map[string]string)
	)

	reURL := regexp.MustCompile(`curl\s+'([^']+)'`)
	if m := reURL.FindStringSubmatch(curlLine); len(m) > 1 {
		urlStr = m[1]
	}
	reHeader := regexp.MustCompile(`-H\s+'([^:]+):\s?([^']*)'`)
	for _, m := range reHeader.FindAllStringSubmatch(curlLine, -1) {
		headers[m[1]] = m[2]
	}
	reCookie := regexp.MustCompile(`-b\s+'([^']+)'`)
	if m := reCookie.FindStringSubmatch(curlLine); len(m) > 1 {
		cookieStr = m[1]
	}
	reData := regexp.MustCompile(`--data(?:-raw)?\s+'([^']*)'`)
	if m := reData.FindStringSubmatch(curlLine); len(m) > 1 {
		payload = m[1]
		method = "POST"
	}

	var req *http.Request
	var err error
	if method == "POST" {
		req, err = http.NewRequest(method, urlStr, strings.NewReader(payload))
	} else {
		req, err = http.NewRequest(method, urlStr, nil)
	}
	if err != nil {
		panic(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if cookieStr != "" {
		req.Header.Set("Cookie", cookieStr)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var kbListResp KBListResponse
	if err := json.Unmarshal(body, &kbListResp); err != nil {
		panic(err)
	}
	return kbListResp
}

type ClassDetail struct {
	Title     string
	Teacher   string
	Place     string
	PS        string
	Week      []int
	Day       int
	StartTime int
	EndTime   int
}

func PhraseClassInfo(raw KBListResponse) []ClassDetail {
	fmt.Println("从教务管理获取数据中...")
	var ret []ClassDetail
	for _, kb := range raw.KBList {
		class := ClassDetail{
			Title:   kb.Title,
			Teacher: kb.Teacher,
			Place:   kb.Place,
			PS:      kb.PS,
		}

		weeks := strings.Split(kb.Week, ",")
		for _, w := range weeks {
			if strings.Contains(w, "-") {
				var start, end int
				fmt.Sscanf(w, "%d-%d", &start, &end)
				for i := start; i <= end; i++ {
					class.Week = append(class.Week, i)
				}
			} else {
				var week int
				fmt.Sscanf(w, "%d", &week)
				class.Week = append(class.Week, week)
			}
		}
		fmt.Sscanf(kb.Day, "%d", &class.Day)
		if strings.Contains(kb.Time, "-") {
			var start, end int
			fmt.Sscanf(kb.Time, "%d-%d", &start, &end)
			class.StartTime = start
			class.EndTime = end
		} else {
			var time int
			fmt.Sscanf(kb.Time, "%d", &time)
			class.Week = append(class.Week, time)
		}
		ret = append(ret, class)
	}
	fmt.Println("数据获取完毕...")
	return ret
}

const CalendarHeader = `BEGIN:VCALENDAR
CALSCALE:GREGORIAN
PRODID:-//Apple Inc.//macOS 15.2//EN
VERSION:2.0
X-APPLE-CALENDAR-COLOR:#CC73E1
BEGIN:VTIMEZONE
TZID:Asia/Shanghai
BEGIN:STANDARD
DTSTART:19890917T020000
RRULE:FREQ=YEARLY;UNTIL=19910914T170000Z;BYMONTH=9;BYDAY=3SU
TZNAME:GMT+8
TZOFFSETFROM:+0900
TZOFFSETTO:+0800
END:STANDARD
BEGIN:DAYLIGHT
DTSTART:19910414T020000
RDATE:19910414T020000
TZNAME:GMT+8
TZOFFSETFROM:+0800
TZOFFSETTO:+0900
END:DAYLIGHT
END:VTIMEZONE
`

var basinTime time.Time
var startTimes = []time.Time{
	time.Date(0, 0, 0, 0, 0, 0, 0, time.Local),
	time.Date(0, 0, 0, 8, 0, 0, 0, time.Local),
	time.Date(0, 0, 0, 8, 50, 0, 0, time.Local),
	time.Date(0, 0, 0, 9, 50, 0, 0, time.Local),
	time.Date(0, 0, 0, 10, 40, 0, 0, time.Local),
	time.Date(0, 0, 0, 11, 30, 0, 0, time.Local),
	time.Date(0, 0, 0, 14, 05, 0, 0, time.Local),
	time.Date(0, 0, 0, 14, 55, 0, 0, time.Local),
	time.Date(0, 0, 0, 15, 45, 0, 0, time.Local),
	time.Date(0, 0, 0, 16, 40, 0, 0, time.Local),
	time.Date(0, 0, 0, 17, 30, 0, 0, time.Local),
	time.Date(0, 0, 0, 18, 30, 0, 0, time.Local),
	time.Date(0, 0, 0, 19, 20, 0, 0, time.Local),
	time.Date(0, 0, 0, 20, 10, 0, 0, time.Local),
}
var endTimes = []time.Time{
	time.Date(0, 0, 0, 0, 0, 0, 0, time.Local),
	time.Date(0, 0, 0, 8, 45, 0, 0, time.Local),
	time.Date(0, 0, 0, 9, 35, 0, 0, time.Local),
	time.Date(0, 0, 0, 10, 35, 0, 0, time.Local),
	time.Date(0, 0, 0, 11, 25, 0, 0, time.Local),
	time.Date(0, 0, 0, 12, 15, 0, 0, time.Local),
	time.Date(0, 0, 0, 14, 50, 0, 0, time.Local),
	time.Date(0, 0, 0, 15, 40, 0, 0, time.Local),
	time.Date(0, 0, 0, 16, 30, 0, 0, time.Local),
	time.Date(0, 0, 0, 17, 25, 0, 0, time.Local),
	time.Date(0, 0, 0, 18, 15, 0, 0, time.Local),
	time.Date(0, 0, 0, 19, 15, 0, 0, time.Local),
	time.Date(0, 0, 0, 20, 05, 0, 0, time.Local),
	time.Date(0, 0, 0, 20, 55, 0, 0, time.Local),
}

func GetClassTime(week, day, classStart, classEnd int) (start, end time.Time) {
	if day == 7 {
		day = 0
	}
	retDay := basinTime.AddDate(0, 0, (week-1)*7+day)
	startTime := retDay.Add(time.Hour*time.Duration(startTimes[classStart].Hour()) + time.Minute*time.Duration(startTimes[classStart].Minute()))
	endTime := retDay.Add(time.Hour*time.Duration(endTimes[classEnd].Hour()) + time.Minute*time.Duration(endTimes[classEnd].Minute()))
	return startTime, endTime
}

func CreateCalendar(classInfos []ClassDetail) {
	fmt.Println("请输入第一周的星期日的日期，格式为YYYY-MM-DD: ")
	reader := bufio.NewReader(os.Stdin)
	timeInput, _ := reader.ReadString('\n')
	var year, month, day int
	fmt.Sscanf(timeInput, "%d-%d-%d", &year, &month, &day)
	basinTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	cal := ics.NewCalendar()
	for _, classInfo := range classInfos {
		for _, week := range classInfo.Week {
			event := cal.AddEvent(uuid.New().String())
			event.SetSummary(classInfo.Title)
			event.SetLocation(classInfo.Place)
			event.SetDescription(fmt.Sprintf("教师: %s\n备注: %s", classInfo.Teacher, classInfo.PS))
			startTime, endTime := GetClassTime(week, classInfo.Day, classInfo.StartTime, classInfo.EndTime)
			event.SetStartAt(startTime)
			event.SetEndAt(endTime)
		}
	}
	f, err := os.Create("calendar.ics")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := os.WriteFile("calendar.ics", []byte(cal.Serialize()), 777); err != nil {
		panic(err)
	}
}

func main() {
	// Command line flags
	var (
		webMode = flag.Bool("web", false, "启动Web服务器模式")
		port    = flag.Int("port", 8080, "Web服务器端口号")
		help    = flag.Bool("h", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		fmt.Println("WHU课表转日历工具")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  ./whu_class            # 命令行模式（默认）")
		fmt.Println("  ./whu_class -web       # Web服务器模式")
		fmt.Println("  ./whu_class -web -port=9000  # 指定端口的Web服务器模式")
		fmt.Println()
		fmt.Println("选项:")
		flag.PrintDefaults()
		return
	}

	if *webMode {
		// Web server mode
		StartWebServer(*port)
		return
	}

	// CLI mode (original functionality)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("程序出错:", r)
			fmt.Println("请检查输入的cURL是否正确，或在github中反馈问题。")
			fmt.Println("按回车键退出...")
			fmt.Scanln()
		}
	}()
	var curl string
	for {
		curl = input()
		if curl != "" {
			break
		} else {
			fmt.Println("输入失败，请重新输入。")
		}
	}
	kbListResp := getKBList(curl)
	classDetail := PhraseClassInfo(kbListResp)
	CreateCalendar(classDetail)
	fmt.Println("日历文件已生成，文件名为calendar.ics")
	fmt.Println("按回车键退出...")
	fmt.Scanln()
}
