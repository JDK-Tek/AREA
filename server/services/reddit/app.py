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

def get_token_from_id(id, service):
	try:
		with db.cursor() as cur:
			cur.execute("SELECT token FROM tokens " \
				"WHERE owner = %s AND service = %s", (
					id,
					service,
				)
			)
			rows = cur.fetchone()
			if not rows:
				return None
			return rows[0]
	except (Exception, psycopg2.Error) as err:
		return None
	


## #
##
## INITIALIZATION
##
##

class NewOreo:
	TYPE_REACTIONS = "reaction"
	TYPE_ACTIONS = "action"

	def __init__(self, service, color, image):
		self.service = service
		self.color = color
		self.image = image
		self.areas = []
	
	def create_area(self, name, type, description, spices):
		self.areas.append({
			"name": name,
			"type": type,
			"description": description,
			"spices": spices
		})

## #
##
## ROUTES
##
##

app = Flask(__name__)
oreo = NewOreo(
	service="reddit",
	color="#ff4500",
	image="/assets/services/reddit.webp"
)

PERMISSIONS_REQUIRED = [
	"creddits",
	"modnote",
	"modcontributors",
	"modmail",
	"modconfig",
	"subscribe",
	"structuredstyles",
	"vote",
	"wikiedit",
	"mysubreddits",
	"submit",
	"modlog",
	"modposts",
	"modflair",
	"announcements",
	"save",
	"modothers",
	"read",
	"privatemessages",
	"report",
	"identity",
	"livemanage",
	"account",
	"modtraffic",
	"wikiread",
	"edit",
	"modwiki",
	"modself",
	"history",
	"flair"
]

##
## OAUTH2
##

def generate_beared_token(id):
	return jwt.encode({
		"id": id,
		"exp": dt.datetime.now() + dt.timedelta(seconds=EXPIRATION)
	}, BACKEND_KEY, algorithm="HS256")

@app.route('/oauth', methods=["GET", "POST"])
def oauth():
	# get the URL of the oauth2 reddit
	scopes = "%20".join(PERMISSIONS_REQUIRED)
	if request.method == "GET":
		reddit_auth_url = (
			f"{REDDIT_AUTH_URL}"
			f"client_id={API_CLIENT_ID_TOKEN}&response_type=code&state=random_string"
			f"&redirect_uri={REDIRECT_URI}&duration=permanent&scope={scopes}"
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

		

		#########



		# get informations about the reddit user
		user_info_url = "https://oauth.reddit.com/api/v1/me"
		headers.update({"Authorization": f"Bearer {reddit_access_token}"})

		user_info_response = requests.get(user_info_url, headers=headers)
		if user_info_response.status_code != 200:
			return jsonify({"error": "Failed to fetch user info"}), user_info_response.status_code

		user_info = user_info_response.json()
		reddit_user_id = user_info.get("id")



		#########



		# data treatment
		area_bearer_token = retrieve_token(get_beared_token(request))
		userid = area_bearer_token.get("id", None) if area_bearer_token else None
	
		app.logger.info(f"Header: {request.headers}")
		app.logger.info(f"Bear token: {area_bearer_token}")
		app.logger.info(f"Area user id: {userid}")

		# user is not logged in an area account
		if not area_bearer_token or not userid:
			try:
				with db.cursor() as cur:
					cur.execute("SELECT owner FROM tokens " \
				 		"WHERE userid = %s AND service = %s", (
							reddit_user_id,
							oreo.service,
						)
					)
					rows = cur.fetchone()
			
					# service account not linked with any area account: create new token and new area account
					if not rows:
						# create new area account empty entry
						cur.execute("INSERT INTO users " \
							"DEFAULT VALUES " \
							"RETURNING id"
						)
						userid = cur.fetchone()[0]

						# create new token linked with the new area account
						cur.execute("INSERT INTO tokens " \
							"(service, token, refresh, userid, owner) " \
							"VALUES (%s, %s, %s, %s, %s)", (
								oreo.service,
		 						reddit_access_token,
								reddit_refresh_token,
								reddit_user_id,
								userid,
							)
						)

						db.commit()
						return jsonify({ "token": generate_beared_token(userid) }), 200
				
					# service account already linked with an area account: update token
					else:
						cur.execute(
								"UPDATE tokens " \
								"SET token = %s, refresh = %s " \
								"WHERE userid = %s AND service = %s " \
								"RETURNING owner", (
									reddit_access_token,
									reddit_refresh_token,
									reddit_user_id,
									oreo.service,
								)
						)
						
						userid = cur.fetchone()[0]
			
						db.commit()
						return jsonify({ "token": generate_beared_token(userid) }), 200
			except (Exception, psycopg2.Error) as err:
				return jsonify({ "error":  str(err)}), 400

		
		
		# user is already logged in an area account
		else:
			try:
				with db.cursor() as cur:
					# check if the reddit account is already linked with an area account
					cur.execute("SELECT owner FROM tokens " \
						"WHERE userid = %s AND service = %s", (
							reddit_user_id,
							oreo.service,
						)
					)
					rows = cur.fetchone()

					# reddit account not linked with any area account: create new token
					if not rows:
						cur.execute("INSERT INTO tokens " \
							"(service, token, refresh, userid, owner) " \
							"VALUES (%s, %s, %s, %s, %s)", (
								oreo.service,
								reddit_access_token,
								reddit_refresh_token,
								reddit_user_id,
								userid,
							)
						)

						db.commit()
						return jsonify({ "token": generate_beared_token(userid) }), 200

					# reddit account already linked with an area account (same account): update token
					elif rows[0] == userid:
						cur.execute(
							"UPDATE tokens " \
							"SET token = %s, refresh = %s " \
							"WHERE userid = %s AND service = %s " \
							"RETURNING owner", (
								reddit_access_token,
								reddit_refresh_token,
								reddit_user_id,
								oreo.service,
							)
						)
						userid = cur.fetchone()[0]
						db.commit()
						return jsonify({ "token": generate_beared_token(userid) }), 200

					# reddit account already linked with an area account (different account):forbidden
					else:
						return jsonify({ "error": "reddit account already linked with an area account" }), 403
			except (Exception, psycopg2.Error) as err:
				return jsonify({ "error":  str(err)}), 400

		return jsonify({ "error": "unexpected end of code"}), 500


##
## ACTIONS
##

##
## REACTIONS
##

# Submit a new post on a subreddit
oreo.create_area(
	name="submit-new-post",
	type=oreo.TYPE_REACTIONS,
	description="Submit a new post on a subreddit",
	spices=[
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
)
@app.route('/submit-new-post', methods=["POST"])
def submit_new_post():
    app.logger.info("submit-new-post endpoint hit")

    if not request.is_json:
        app.logger.error("Request is not valid JSON")
        return jsonify({"error": "Invalid JSON"}), 400

    userid = request.json.get("userid")
    if not userid:
        app.logger.error(f"Missing required fields: 'userid'")
        return jsonify({"error": "Missing required fields"}), 400

    access_token = get_token_from_id(request.json.get("userid"), "reddit")

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

    app.logger.info(f"User {userid} submitted a post on r/{subreddit}: {title}")
    return jsonify({"status": "Post submitted"}), 200


# Submit a new link on a subreddit
oreo.create_area(
	name="submit-new-link",
	type=oreo.TYPE_REACTIONS,
	description="Submit a new link on a subreddit",
	spices=[
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
			"title": "The url to post",
			"name": "url",
			"type": "text"
		}
	]
)
@app.route('/submit-new-link', methods=["POST"])
def submit_new_link():
    app.logger.info("submit-new-link endpoint hit")

    if not request.is_json:
        app.logger.error("Request is not valid JSON")
        return jsonify({"error": "Invalid JSON"}), 400

    userid = request.json.get("userid")
    if not userid:
        app.logger.error("Missing required fields: 'userid'")
        return jsonify({"error": "Missing required fields"}), 400

    access_token = get_token_from_id(request.json.get("userid"), "reddit")

    spices = request.json.get("spices", {})
    subreddit = spices.get("subreddit")
    title = spices.get("title")
    url = spices.get("url")

    if not subreddit or not title or not url:
        app.logger.error("Missing required fields: subreddit=%s, title=%s, url=%s", subreddit, title, url)
        return jsonify({"error": "Missing required fields"}), 400

    reddit_submit_url = "https://oauth.reddit.com/api/submit"
    headers = {
        "User-Agent": "area/1.0",
        "Authorization": f"Bearer {access_token}",
    }
    body = {
        "kind": "link",
        "sr": subreddit,
        "title": title,
        "url": url,
    }

    res = requests.post(reddit_submit_url, headers=headers, data=body)

    if res.status_code != 200:
        app.logger.error("Failed to submit link: %s", res.json())
        return jsonify({
            "error": "Failed to submit link",
            "details": res.json()
        }), res.status_code

    app.logger.info(f"User {userid} submitted a link on r/{subreddit}: {title}")
    return jsonify({"status": "Link submitted"}), 200


# Reply to a post
oreo.create_area(
	name="reply-post",
	type=oreo.TYPE_REACTIONS,
	description="Reply to a post on a subreddit",
	spices=[
		{
			"title": "Post id (e.g. t3_<id>)",
			"name": "post_id",
			"type": "input"
		},
		{
			"title": "Reply message",
			"name": "reply_msg",
			"type": "text"
		}
	]
)
@app.route('/reply-post', methods=["POST"])
def reply_post():
    app.logger.info("reply-post endpoint hit")

    if not request.is_json:
        app.logger.error("Request is not valid JSON")
        return jsonify({"error": "Invalid JSON"}), 400

    userid = request.json.get("userid")
    if not userid:
        app.logger.error("Missing required fields: 'userid'")
        return jsonify({"error": "Missing required fields"}), 400

    access_token = get_token_from_id(request.json.get("userid"), "reddit")

    spices = request.json.get("spices", {})
    post_id = spices.get("post_id")
    reply_msg = spices.get("reply_msg")

    if not post_id or not reply_msg:
        app.logger.error("Missing required fields: post_id=%s, reply_msg=%s", post_id, reply_msg)
        return jsonify({"error": "Missing required fields"}), 400

    reddit_submit_url = "https://oauth.reddit.com/api/comment"
    headers = {
        "User-Agent": "area/1.0",
        "Authorization": f"Bearer {access_token}",
    }
    body = {
        "api_type": "json",
		"thing_id": post_id,
		"text": reply_msg
    }

    res = requests.post(reddit_submit_url, headers=headers, data=body)

    if res.status_code != 200:
        app.logger.error("Failed to submit link: %s", res.json())
        return jsonify({
            "error": "Failed to submit link",
            "details": res.json()
        }), res.status_code

    app.logger.info(f"User {userid} reply to the '{post_id}': {reply_msg}")
    return jsonify({"status": "Reply submitted"}), 200


# Reply to a message
oreo.create_area(
	name="reply-message",
	type=oreo.TYPE_REACTIONS,
	description="Reply to a private message",
	spices=[
		{
			"title": "Post id (e.g. t4_<id>)",
			"name": "message_id",
			"type": "input"
		},
		{
			"title": "Reply message",
			"name": "reply_msg",
			"type": "text"
		}
	]
)
@app.route('/reply-message', methods=["POST"])
def reply_message():
    app.logger.info("reply-message endpoint hit")

    if not request.is_json:
        app.logger.error("Request is not valid JSON")
        return jsonify({"error": "Invalid JSON"}), 400

    userid = request.json.get("userid")
    if not userid:
        app.logger.error(f"Missing required fields: 'userid'")
        return jsonify({"error": "Missing required fields"}), 400

    access_token = get_token_from_id(request.json.get("userid"), "reddit")

    spices = request.json.get("spices", {})
    message_id = spices.get("message_id")
    reply_msg = spices.get("reply_msg")

    if not message_id or not reply_msg:
        app.logger.error("Missing required fields: message_id=%s, reply_msg=%s", message_id, reply_msg)
        return jsonify({"error": "Missing required fields"}), 400

    reddit_submit_url = "https://oauth.reddit.com/api/comment"
    headers = {
        "User-Agent": "area/1.0",
        "Authorization": f"Bearer {access_token}",
    }
    body = {
        "api_type": "json",
        "thing_id": message_id,
        "text": reply_msg,
    }

    res = requests.post(reddit_submit_url, headers=headers, data=body)

    if res.status_code != 200:
        app.logger.error("Failed to submit reply: %s", res.json())
        return jsonify({
            "error": "Failed to submit reply",
            "details": res.json()
        }), res.status_code

    app.logger.info(f"User {userid} replied to message '{message_id}': {reply_msg}")
    return jsonify({"status": "Reply submitted"}), 200







##
## INFORMATIONS ABOUT ACTION/REACTION OF REDDIT SERVICE
##
@app.route('/', methods=["GET"])
def info():
	res = {
		"service": oreo.service,
		"color": oreo.color,
		"image": oreo.image,
		"areas": oreo.areas
	}
	return jsonify(res), 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80, debug=True)
