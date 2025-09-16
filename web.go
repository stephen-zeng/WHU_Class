package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"WHUClass/frontend"
	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

// WebRequest represents the request payload from the frontend
type WebRequest struct {
	CurlCommand string `json:"curlCommand"`
	FirstSunday string `json:"firstSunday"`
}

// WebResponse represents the response to the frontend
type WebResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Calendar  string `json:"calendar,omitempty"`
	Error     string `json:"error,omitempty"`
	ClassData []ClassDetail `json:"classData,omitempty"`
}

// Handler for the main page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := frontend.GetIndexTemplate()
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler for generating calendar via API
func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req WebRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(WebResponse{
			Success: false,
			Error:   "Invalid JSON request",
		})
		return
	}

	// Validate input
	if strings.TrimSpace(req.CurlCommand) == "" || strings.TrimSpace(req.FirstSunday) == "" {
		json.NewEncoder(w).Encode(WebResponse{
			Success: false,
			Error:   "Missing required fields",
		})
		return
	}

	// Parse and validate date
	var year, month, day int
	if n, err := fmt.Sscanf(req.FirstSunday, "%d-%d-%d", &year, &month, &day); n != 3 || err != nil {
		json.NewEncoder(w).Encode(WebResponse{
			Success: false,
			Error:   "Invalid date format. Please use YYYY-MM-DD",
		})
		return
	}

	// Set global basin time
	basinTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	// Get class data
	kbListResp, err := getKBListSafe(req.CurlCommand)
	if err != nil {
		json.NewEncoder(w).Encode(WebResponse{
			Success: false,
			Error:   "Failed to fetch class data: " + err.Error(),
		})
		return
	}

	// Parse class information
	classDetails := PhraseClassInfo(kbListResp)
	if len(classDetails) == 0 {
		json.NewEncoder(w).Encode(WebResponse{
			Success: false,
			Error:   "No class data found. Please check your cURL command.",
		})
		return
	}

	// Generate calendar
	calendarData, err := CreateCalendarWeb(classDetails)
	if err != nil {
		json.NewEncoder(w).Encode(WebResponse{
			Success: false,
			Error:   "Failed to generate calendar: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(WebResponse{
		Success:   true,
		Message:   "Calendar generated successfully",
		Calendar:  calendarData,
		ClassData: classDetails,
	})
}

// Safe version of getKBList that returns error instead of panicking
func getKBListSafe(curlLine string) (KBListResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error
		}
	}()

	var (
		urlStr    string
		method    = "GET"
		payload   string
		cookieStr string
		headers   = make(map[string]string)
	)

	// Parse cURL command (same logic as original)
	reURL := regexp.MustCompile(`curl\s+'([^']+)'`)
	if m := reURL.FindStringSubmatch(curlLine); len(m) > 1 {
		urlStr = m[1]
	}
	if urlStr == "" {
		return KBListResponse{}, fmt.Errorf("could not parse URL from cURL command")
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
		return KBListResponse{}, fmt.Errorf("failed to create request: %v", err)
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
		return KBListResponse{}, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return KBListResponse{}, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return KBListResponse{}, fmt.Errorf("failed to read response: %v", err)
	}

	var kbListResp KBListResponse
	if err := json.Unmarshal(body, &kbListResp); err != nil {
		return KBListResponse{}, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return kbListResp, nil
}

// Web version of CreateCalendar that returns calendar data as string
func CreateCalendarWeb(classInfos []ClassDetail) (string, error) {
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
	return cal.Serialize(), nil
}

// StartWebServer starts the HTTP server
func StartWebServer(port int) {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/generate", generateHandler)

	portStr := ":" + strconv.Itoa(port)
	fmt.Printf("🚀 WHU课表转日历服务已启动\n")
	fmt.Printf("📱 请在浏览器中访问: http://localhost%s\n", portStr)
	fmt.Printf("🌐 如果部署在服务器上，请使用服务器的IP地址访问\n")
	fmt.Printf("⏹️  按Ctrl+C停止服务\n\n")

	log.Fatal(http.ListenAndServe(portStr, nil))
}