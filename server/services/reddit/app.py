import os
import sys
import jwt
import requests
import psycopg2
import datetime as dt
from dotenv import load_dotenv
from flask import Flask, jsonify, request

sys.stdout.reconfigure(line_buffering=True)
load_dotenv("/usr/mount.d/.env")

## #
##
## ENV VARIABLES
##
##

API_CLIENT_ID_TOKEN = os.environ.get("API_CLIENT_ID_TOKEN")
API_CLIENT_SECRET_TOKEN = os.environ.get("API_CLIENT_SECRET_TOKEN")

REDIRECT_URI = os.environ.get("REDIRECT")
REDDIT_AUTH_URL = "https://www.reddit.com/api/v1/authorize?"
REDDIT_TOKEN_URL = "https://www.reddit.com/api/v1/access_token"

EXPIRATION = 60 * 30
BACKEND_KEY = os.environ.get("BACKEND_KEY")


## #
##
## DATABASE CONNECTION
##
##
while True:
	try:
		db = psycopg2.connect(
			database=os.environ.get("DB_NAME"),
			user=os.environ.get("DB_USER"),
			password=os.environ.get("DB_PASSWORD"),
			host=os.environ.get("DB_HOST"),
			port=os.environ.get("DB_PORT")
		)
		break
	except psycopg2.OperationalError:
		continue


## #
##
## ROUTES
##
##

app = Flask(__name__)

##
## OAUTH2
##
@app.route('/oauth', methods=["GET", "POST"])
def oauth():
	# get the URL of the oauth2 reddit
	if request.method == "GET":
		reddit_auth_url = (
			f"{REDDIT_AUTH_URL}"
			f"client_id={API_CLIENT_ID_TOKEN}&response_type=code&state=random_string"
			f"&redirect_uri={REDIRECT_URI}&duration=permanent&scope=identity"
		)
		return reddit_auth_url

	# get the acces token
	if request.method == "POST":

		# get the reddit access token
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

		# reddit acces_token:
		reddit_access_token = token_data.get("access_token")
		reddit_refresh_token = token_data.get("refresh_token")

		


		# get informations about the reddit user
		user_info_url = "https://oauth.reddit.com/api/v1/me"
		headers.update({"Authorization": f"Bearer {reddit_access_token}"})

		user_info_response = requests.get(user_info_url, headers=headers)
		if user_info_response.status_code != 200:
			return jsonify({"error": "Failed to fetch user info"}), user_info_response.status_code

		user_info = user_info_response.json()
		reddit_user_name = user_info.get("name")
		reddit_user_id = user_info.get("id")
	
		app.logger.info("refresh_token: %s", reddit_refresh_token)

		# create a new user in the bdd
		try:
			with db.cursor() as cur:
				tokenid, ownerid = -1, -1
				cur.execute("SELECT id, owner FROM tokens WHERE userid = %s", (reddit_user_id,))
				rows = cur.fetchone()
				if not rows:
					cur.execute("insert into tokens" \
						"(service, token, refresh, userid)" \
						"values (%s, %s, %s, %s)" \
						"returning id", \
							("reddit", reddit_access_token, reddit_refresh_token, reddit_user_id,)
					)
					r = cur.fetchone()
					if not r:
						raise Exception("could not fetch")
					tokenid = r[0]
					cur.execute("insert into users (tokenid) values (%s) returning id", (tokenid,))
					r = cur.fetchone()
					if not r:
						raise Exception("could not fetch")
					ownerid = r[0]
					db.commit()
				else:
					tokenid, ownerid = rows[0], rows[1]
				db.commit()
				data = jwt.encode({
					"id": ownerid,
					"exp": dt.datetime.now() + dt.timedelta(seconds=EXPIRATION)
				}, BACKEND_KEY, algorithm="HS256")
				return jsonify({ "token": data }), 200
			
		
		except (Exception, psycopg2.Error) as err:
			return jsonify({ "error":  str(err)}), 400
		return jsonify({ "error": "unexpected end of code"}), 500


@app.route('/send', methods=["POST"])
def send():
	return jsonify({"status": "caca"}), 200



##
## INFORMATIONS ABOUT ACTION/REACTION OF REDDIT SERVICE
##
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
