from flask import Flask, request, jsonify, Response
from dotenv import load_dotenv
from urllib.parse import urlencode
# import json
import os

load_dotenv("/usr/mount.d/.env")

CLIENT_ID = os.environ.get("CLIENT_ID")

API_URL = "https://authorize.roblox.com/"

app = Flask(__name__)

def error(something, code):
	return jsonify(something), code

@app.route('/oauth', methods=["GET", "POST"])
def oauth():
	if request.method == "GET":
		params = {
			"client_id": CLIENT_ID,
			"response_type": "code",
			"redirect_uri": os.environ.get("REDIRECT"),
			"scope": "openid",
			"state": 6789,
			"nonce": 12345,
			"setp": "accountConfirm"
		}
		return API_URL + "?" + urlencode(params)

	if request.method == "POST":
		req = request.get_json()
		if not "code" in req:
			return error({ "error": "missing code" }, 400)
		print("test")
		print(req)
		return req

@app.route('/')
def hello():
	return "Hello World!"

if __name__ == '__main__':
	app.run(host='0.0.0.0', port=80)