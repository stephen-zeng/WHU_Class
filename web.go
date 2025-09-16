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

const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WHU课表转日历</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .section {
            margin-bottom: 25px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: #555;
        }
        textarea, input {
            width: 100%;
            padding: 12px;
            border: 2px solid #ddd;
            border-radius: 5px;
            font-size: 14px;
            box-sizing: border-box;
        }
        textarea {
            height: 150px;
            resize: vertical;
            font-family: monospace;
        }
        input:focus, textarea:focus {
            outline: none;
            border-color: #007bff;
        }
        button {
            background-color: #007bff;
            color: white;
            padding: 12px 24px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            width: 100%;
            margin-top: 10px;
        }
        button:hover {
            background-color: #0056b3;
        }
        button:disabled {
            background-color: #ccc;
            cursor: not-allowed;
        }
        .result {
            margin-top: 20px;
            padding: 15px;
            border-radius: 5px;
            display: none;
        }
        .success {
            background-color: #d4edda;
            border: 1px solid #c3e6cb;
            color: #155724;
        }
        .error {
            background-color: #f8d7da;
            border: 1px solid #f5c6cb;
            color: #721c24;
        }
        .loading {
            text-align: center;
            color: #666;
        }
        .instructions {
            background-color: #e9ecef;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 20px;
            font-size: 14px;
            line-height: 1.5;
        }
        .instructions h3 {
            margin-top: 0;
            color: #495057;
        }
        .class-info {
            margin-top: 15px;
        }
        .class-item {
            background-color: #f8f9fa;
            padding: 10px;
            margin: 5px 0;
            border-radius: 3px;
            border-left: 4px solid #007bff;
        }
        .download-btn {
            background-color: #28a745;
            margin-top: 10px;
            display: none;
        }
        .download-btn:hover {
            background-color: #218838;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎓 WHU课表转日历工具</h1>
        
        <div class="instructions">
            <h3>使用说明：</h3>
            <ol>
                <li>进入本科教务系统的个人课表查询页面</li>
                <li>打开浏览器开发者工具（F12），切换到Network选项卡</li>
                <li>点击查询课表，完成验证码验证</li>
                <li>在Network选项卡找到课表数据请求，右键选择"Copy as cURL"</li>
                <li>将完整的cURL命令粘贴到下方文本框中</li>
                <li>填写第一周星期日的日期</li>
                <li>点击生成日历文件</li>
            </ol>
        </div>

        <form id="calendarForm">
            <div class="section">
                <label for="curlCommand">完整cURL命令：</label>
                <textarea id="curlCommand" name="curlCommand" placeholder="请粘贴从浏览器复制的完整cURL命令..." required></textarea>
            </div>
            
            <div class="section">
                <label for="firstSunday">第一周星期日日期（格式：YYYY-MM-DD）：</label>
                <input type="date" id="firstSunday" name="firstSunday" required>
            </div>
            
            <button type="submit" id="submitBtn">生成日历文件</button>
        </form>

        <div id="result" class="result"></div>
        
        <button id="downloadBtn" class="download-btn" style="display:none;">下载日历文件</button>
    </div>

    <script>
        let calendarData = '';
        
        document.getElementById('calendarForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const curlCommand = document.getElementById('curlCommand').value;
            const firstSunday = document.getElementById('firstSunday').value;
            
            if (!curlCommand.trim() || !firstSunday) {
                showResult('请填写所有必需字段', 'error');
                return;
            }
            
            const submitBtn = document.getElementById('submitBtn');
            const resultDiv = document.getElementById('result');
            
            submitBtn.disabled = true;
            submitBtn.textContent = '处理中...';
            resultDiv.style.display = 'block';
            resultDiv.className = 'result loading';
            resultDiv.innerHTML = '正在处理课表数据，请稍候...';
            
            fetch('/api/generate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    curlCommand: curlCommand,
                    firstSunday: firstSunday
                })
            })
            .then(response => response.json())
            .then(data => {
                submitBtn.disabled = false;
                submitBtn.textContent = '生成日历文件';
                
                if (data.success) {
                    calendarData = data.calendar;
                    let message = '日历文件生成成功！';
                    if (data.classData && data.classData.length > 0) {
                        message += '<div class="class-info"><strong>发现课程：</strong>';
                        data.classData.forEach(cls => {
                            message += '<div class="class-item">';
                            message += '<strong>' + cls.Title + '</strong><br>';
                            message += '教师：' + cls.Teacher + '<br>';
                            message += '地点：' + cls.Place + '<br>';
                            message += '周数：' + cls.Week.join(', ') + '<br>';
                            message += '时间：星期' + cls.Day + ' 第' + cls.StartTime + '-' + cls.EndTime + '节';
                            if (cls.PS) {
                                message += '<br>备注：' + cls.PS;
                            }
                            message += '</div>';
                        });
                        message += '</div>';
                    }
                    showResult(message, 'success');
                    document.getElementById('downloadBtn').style.display = 'block';
                } else {
                    showResult('生成失败: ' + (data.error || data.message), 'error');
                }
            })
            .catch(error => {
                submitBtn.disabled = false;
                submitBtn.textContent = '生成日历文件';
                showResult('请求失败: ' + error.message, 'error');
            });
        });
        
        document.getElementById('downloadBtn').addEventListener('click', function() {
            if (calendarData) {
                const blob = new Blob([calendarData], { type: 'text/calendar' });
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'whu_calendar.ics';
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);
                document.body.removeChild(a);
            }
        });
        
        function showResult(message, type) {
            const resultDiv = document.getElementById('result');
            resultDiv.style.display = 'block';
            resultDiv.className = 'result ' + type;
            resultDiv.innerHTML = message;
        }
        
        // Set default first Sunday to next Sunday
        function setDefaultFirstSunday() {
            const today = new Date();
            const nextSunday = new Date(today);
            const daysUntilSunday = 7 - today.getDay();
            if (daysUntilSunday === 7) {
                nextSunday.setDate(today.getDate());
            } else {
                nextSunday.setDate(today.getDate() + daysUntilSunday);
            }
            
            const year = nextSunday.getFullYear();
            const month = String(nextSunday.getMonth() + 1).padStart(2, '0');
            const day = String(nextSunday.getDate()).padStart(2, '0');
            
            document.getElementById('firstSunday').value = year + '-' + month + '-' + day;
        }
        
        setDefaultFirstSunday();
    </script>
</body>
</html>`

// Handler for the main page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexHTML)
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