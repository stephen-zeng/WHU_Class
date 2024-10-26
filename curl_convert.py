import uncurl
import json

def curl_convert():
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