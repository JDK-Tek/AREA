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
## UTILS
##
##

def get_beared_token(request):
	token = request.headers.get("Authorization")
	if not token:
		return None
	if token.startswith("Bearer "):
		return token[7:]
	return None

def retrieve_token(token):
	try:
		data = jwt.decode(token, BACKEND_KEY, algorithms=["HS256"])
		return data
	except jwt.ExpiredSignatureError:
		return None
	except jwt.InvalidTokenError:
		return None

def retrieve_user_token(id):
	try:
		with db.cursor() as cur:
			cur.execute("SELECT token FROM tokens WHERE service = 'reddit' AND owner = %s", (id,))
			rows = cur.fetchone()
			if not rows:
				return None
			return rows[0]
	except (Exception, psycopg2.Error) as err:
		return None

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
			f"&redirect_uri={REDIRECT_URI}&duration=permanent&scope=identity%20read%20submit%20vote%20privatemessages%20save"
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
		headers = {"User-Agent": "area/1.0"}
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

		# create a new user in the bdd
		try:
			with db.cursor() as cur:
				tokenid, ownerid = -1, -1
				cur.execute("SELECT id, owner FROM tokens WHERE userid = %s", (reddit_user_id,))
				rows = cur.fetchone()
				if not rows:
					cur.execute("INSERT INTO tokens" \
						"(service, token, refresh, userid)" \
						"VALUES (%s, %s, %s, %s)" \
						"RETURNING id", \
							("reddit", reddit_access_token, reddit_refresh_token, reddit_user_id,)
					)
					r = cur.fetchone()
					if not r:
						raise Exception("could not fetch")
					tokenid = r[0]
					cur.execute("INSERT INTO users (tokenid) VALUES (%s) RETURNING id", (tokenid,))
					r = cur.fetchone()
					if not r:
						raise Exception("could not fetch")
					ownerid = r[0]
					cur.execute("UPDATE tokens SET owner = %s WHERE id = %s", (ownerid, tokenid))
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

##
## ACTIONS
##

# @app.route('/new-post-save-by-me', methods=["POST"])
# def new_post_save_by_me():
# 	return jsonify({"status": "caca"}), 200

##
## REACTIONS
##

# Submit a new post on a subreddit
@app.route('/submit-new-post', methods=["POST"])
def submit_new_post():
    app.logger.info("submit-new-post endpoint hit")
    user = retrieve_token(get_beared_token(request))
    if not user:
        app.logger.error("Invalid area token")
        return jsonify({"error": "Invalid area token"}), 401

    access_token = retrieve_user_token(user.get("id"))
    if not access_token:
        app.logger.error("Invalid reddit token")
        return jsonify({"error": "Invalid reddit token"}), 401

    if not request.is_json:
        app.logger.error("Request is not valid JSON")
        return jsonify({"error": "Invalid JSON"}), 400

    spices = request.json.get("spices", {})
    subreddit = spices.get("subreddit")
    title = spices.get("title")
    content = spices.get("content")

    if not subreddit or not title or not content:
        app.logger.error("Missing required fields: subreddit=%s, title=%s, content=%s", subreddit, title, content)
        return jsonify({"error": "Missing required fields"}), 400

    reddit_submit_url = "https://oauth.reddit.com/api/submit"
    headers = {
        "User-Agent": "area/1.0",
        "Authorization": f"Bearer {access_token}",
    }
    body = {
        "kind": "self",
        "sr": subreddit,
        "title": title,
        "text": content,
    }

    res = requests.post(reddit_submit_url, headers=headers, data=body)

    if res.status_code != 200:
        app.logger.info("Failed to submit post: %s", res.json())
        return jsonify({
            "error": "Failed to submit post",
            "details": res.json()
        }), res.status_code

    app.logger.info(f"User {user.get('id')} submitted a post on r/{subreddit}: {title}")
    return jsonify({"status": "Post submitted"}), 200

##
## INFORMATIONS ABOUT ACTION/REACTION OF REDDIT SERVICE
##
@app.route('/', methods=["GET"])
def info():
	res = {
		"color": "#ff4500",
		"image": "http://link.com",
		"areas": [
			{
				"name": "submit-new-post",
				"type": "reaction",
				"description": "Submit a new post on a subreddit",
				"spices": [
					{
						"title": "Title of the post",
						"name": "title",
						"type": "input"
					},
					{
						"title": "Subreddit where to post without the r/",
						"name": "subreddit",
						"type": "input"
					},
					{
						"title": "Content of the post",
						"name": "content",
						"type": "text"
					}
				]
			}
		]
	}


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80, debug=True)
