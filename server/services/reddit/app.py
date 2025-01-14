import os
import sys
import requests
from flask import Flask, jsonify, request
from dotenv import load_dotenv

sys.stdout.reconfigure(line_buffering=True)


load_dotenv("/usr/mount.d/.env")

API_CLIENT_ID_TOKEN = os.environ.get("API_CLIENT_ID_TOKEN")
API_CLIENT_SECRET_TOKEN = os.environ.get("API_CLIENT_SECRET_TOKEN")

REDIRECT_URI = os.environ.get("REDIRECT")
REDDIT_AUTH_URL = "https://www.reddit.com/api/v1/authorize?"
REDDIT_TOKEN_URL = "https://www.reddit.com/api/v1/access_token"

app = Flask(__name__)

@app.route('/oauth', methods=["GET", "POST"])
def oauth():
	if request.method == "GET":
		reddit_auth_url = (
			f"{REDDIT_AUTH_URL}"
			f"client_id={API_CLIENT_ID_TOKEN}&response_type=code&state=random_string"
			f"&redirect_uri={REDIRECT_URI}&duration=permanent&scope=identity"
		)
		return reddit_auth_url

	if request.method == "POST":
		code = request.json.get('code')
		if not code:
			return jsonify({"error": "Missing code"}), 400

		auth = (API_CLIENT_ID_TOKEN, API_CLIENT_SECRET_TOKEN)
		data = {
			"grant_type": "authorization_code",
			"code": code,
			"redirect_uri": REDIRECT_URI,
		}
		headers = {"User-Agent": "YourApp/1.0"}

		response = requests.post(REDDIT_TOKEN_URL, data=data, headers=headers, auth=auth)

		if response.status_code != 200:
			return jsonify({"error": "Failed to obtain token"}), response.status_code

		token_data = response.json()
		access_token = token_data.get("access_token")
		



		# ###### TEST #####
		user_info_url = "https://oauth.reddit.com/api/v1/me"
		headers.update({"Authorization": f"Bearer {access_token}"})

		user_info_response = requests.get(user_info_url, headers=headers)
		if user_info_response.status_code != 200:
			return jsonify({"error": "Failed to fetch user info"}), user_info_response.status_code

		user_info = user_info_response.json()
		app.logger.info("user_name: %s", user_info.get("name"))
		app.logger.info("user_id: %s", user_info.get("id"))
		app.logger.info("user_token: %s", access_token)


		return jsonify({"access_token": access_token}), 200

@app.route('/send', methods=["POST"])
def send():
	return jsonify({"status": "caca"}), 200

@app.route('/', methods=["GET"])
def info():
	res = {
		"color": "#3243423",
		"image": "http://link.com",
		"areas": [
			{
				"name": "send-message",
				"type": "reaction",
				"description": "Send a message",
				"spices": [
					{
						"title": "The title of the elem",
						"name": "the id name",
						"type": "the type"
					}
				]
			}
		]
	}


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80, debug=True)
