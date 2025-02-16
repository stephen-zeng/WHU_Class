import uncurl
import json
import re

def parse_cookies_from_curl(curl):
    # 使用正则表达式提取 -b 参数后的 cookies
    match = re.search(r"-b '([^']+)'", curl)
    if match:
        cookies_str = match.group(1)
        cookies = dict(item.split("=", 1) for item in cookies_str.split("; "))
        return cookies
    return {}

def curl_convert():
	print("请输入curl：")

	curl = ""

	while True:
		s = input()[:-1]

		if s == "":
			break
		curl += s 

	curl += "'"
	context = uncurl.parse_context(curl)

	headers = json.loads(json.dumps(context.headers))
	cookies = json.loads(json.dumps(context.cookies))
	if cookies == {}:
		cookies = parse_cookies_from_curl(curl)
	return headers, cookies