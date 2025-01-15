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

API_APP_ID=os.environ.get("API_APP_ID")
API_CLIENT_ID_TOKEN=os.environ.get("API_CLIENT_ID_TOKEN")
API_CLIENT_SECRET_TOKEN=os.environ.get("API_CLIENT_SECRET_TOKEN")

REDIRECT_URI = os.environ.get("REDIRECT")
GITHUB_API_OAUTH_URL = "https://github.com/login/oauth"
GITHUB_API_URL = "https://api.github.com"

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
			cur.execute("SELECT token FROM tokens WHERE service = 'github' AND owner = %s", (id,))
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

	def __init__(self, color="#ff4500", image="http://link.com"):
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

oreo = NewOreo()

PERMISSIONS_REQUIRED = [
	"user",
	"repo",
	"public_repo",
	"write:issues"
]

##
## OAUTH2
##
@app.route('/oauth', methods=["GET", "POST"])
def oauth():
	# get the URL of the oauth2 github
	if request.method == "GET":
		scopes = "%20".join(PERMISSIONS_REQUIRED)
		github_auth_url = (
			f"{GITHUB_API_OAUTH_URL}/authorize"
			f"?client_id={API_CLIENT_ID_TOKEN}&response_type=code&state=random_string"
			f"&redirect_uri={REDIRECT_URI}&duration=permanent&scope={scopes}"
		)
		return github_auth_url
	
	# get the token from the code
	if request.method == "POST":
		# get the github access token
		code = request.json.get('code')
		if not code:
			return jsonify({"error": "Missing code"}), 400

		auth = (API_CLIENT_ID_TOKEN, API_CLIENT_SECRET_TOKEN)
		data = {
			"client_id": API_CLIENT_ID_TOKEN,
			"client_secret": API_CLIENT_SECRET_TOKEN,
			"code": code,
			"redirect_uri": REDIRECT_URI
		}
		headers = {
			"User-Agent": "area/1.0",
			"Accept": "application/json"
		}
		response = requests.post(f"{GITHUB_API_OAUTH_URL}/access_token", data=data, headers=headers, auth=auth)

		if response.status_code != 200:
			return jsonify({"error": "Failed to obtain token"}), response.status_code

		token_data = response.json()

		# github acces_token:
		github_access_token = token_data.get("access_token")
		github_refresh_token = token_data.get("refresh_token")
	
		# get informations about the github user
		user_info_url = f"{GITHUB_API_URL}/user"
		headers.update({"Authorization": f"Bearer {github_access_token}"})

		user_info_response = requests.get(user_info_url, headers=headers)
		if user_info_response.status_code != 200:
			return jsonify({"error": "Failed to fetch user info"}), user_info_response.status_code
		
		user_info = user_info_response.json()
		github_user_id = str(user_info.get("id"))

		# create a new user in the bdd
		try:
			with db.cursor() as cur:
				tokenid, ownerid = -1, -1
				cur.execute("SELECT id, owner FROM tokens WHERE userid = %s", (github_user_id,))
				rows = cur.fetchone()
				if not rows:
					cur.execute("INSERT INTO tokens" \
						"(service, token, refresh, userid)" \
						"VALUES (%s, %s, %s, %s)" \
						"RETURNING id", \
							("github", github_access_token, github_refresh_token, github_user_id,)
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

##
## REACTIONS
##

# Create a new issue
oreo.create_area(
	"create-issue",
	NewOreo.TYPE_REACTIONS,
	"Create a new issue in a repository",
	[
		{
			"name": "owner",
			"type": "input",
			"description": "The owner of the repository"
		},
		{
			"name": "repo",
			"type": "input",
			"description": "The repository name"
		},
		{
			"name": "title",
			"type": "input",
			"description": "The title of the issue (Markdown supported)"
		},
		{
			"name": "body",
			"type": "text",
			"description": "The body of the issue (Markdown supported)"
		}
	]
)
@app.route('/create-issue', methods=["POST"])
def create_issue():
    app.logger.info("create-issue endpoint hit")
    user = retrieve_token(get_beared_token(request))
    if not user:
        app.logger.error("Invalid area token")
        return jsonify({"error": "Invalid area token"}), 401

    access_token = retrieve_user_token(user.get("id"))
    if not access_token:
        app.logger.error("Invalid github token")
        return jsonify({"error": "Invalid github token"}), 401

    if not request.is_json:
        app.logger.error("Request is not valid JSON")
        return jsonify({"error": "Invalid JSON"}), 400

    spices = request.json.get("spices", {})
    owner = spices.get("owner")
    repo = spices.get("repo")
    title = spices.get("title")
    body = spices.get("body", "")

    if not owner or not repo:
        app.logger.error("Missing required fields: owner=%s, repo=%s", owner, repo)
        return jsonify({"error": "Missing required fields"}), 400

    github_submit_url = f"{GITHUB_API_URL}/repos/{owner}/{repo}/issues"
    headers = {
        "User-Agent": "area/1.0",
        "Authorization": f"Bearer {access_token}",
        "Accept": "application/vnd.github+json"
    }
    body = {
        "title": title,
		"body": body,
    }

    res = requests.post(github_submit_url, headers=headers, json=body)


    if res.status_code != 201:
        app.logger.info("Failed to create issue: %s", res.json())
        return jsonify({
            "error": "Failed to create issue",
            "details": res.json()
        }), res.status_code

    app.logger.info(f"User {user.get('id')} created a new issue in {owner}/{repo}: {title}")
    return jsonify({"status": "Issue created"}), 200




# Create a new reply to an issue / pull request
oreo.create_area(
	"create-reply",
	NewOreo.TYPE_REACTIONS,
	"Create a new reply to an issue or pull request",
	[
		{
			"name": "id",
			"type": "number",
			"description": "The id of the issue or pull request"
		},
		{
			"name": "owner",
			"type": "input",
			"description": "The owner of the repository"
		},
		{
			"name": "repo",
			"type": "input",
			"description": "The repository name"
		},
		{
			"name": "body",
			"type": "text",
			"description": "The body of the reply (Markdown supported)"
		}
	]
)
@app.route('/create-reply', methods=["POST"])
def create_reply():
    app.logger.info("create-reply endpoint hit")
    user = retrieve_token(get_beared_token(request))
    if not user:
        app.logger.error("Invalid area token")
        return jsonify({"error": "Invalid area token"}), 401

    access_token = retrieve_user_token(user.get("id"))
    if not access_token:
        app.logger.error("Invalid github token")
        return jsonify({"error": "Invalid github token"}), 401

    if not request.is_json:
        app.logger.error("Request is not valid JSON")
        return jsonify({"error": "Invalid JSON"}), 400

    spices = request.json.get("spices", {})
    id = spices.get("id")
    owner = spices.get("owner")
    repo = spices.get("repo")
    body = spices.get("body", "")

    if not owner or not repo:
        app.logger.error("Missing required fields: owner=%s, repo=%s", owner, repo)
        return jsonify({"error": "Missing required fields"}), 400

    github_submit_url = f"{GITHUB_API_URL}/repos/{owner}/{repo}/issues/{id}/comments"
    headers = {
        "User-Agent": "area/1.0",
        "Authorization": f"Bearer {access_token}",
        "Accept": "application/vnd.github+json"
    }
    body = {
		"body": body,
    }

    res = requests.post(github_submit_url, headers=headers, json=body)


    if res.status_code != 201:
        app.logger.info("Failed to create issue: %s", res.json())
        return jsonify({
            "error": "Failed to create issue",
            "details": res.json()
        }), res.status_code

    app.logger.info(f"User {user.get('id')} created a new reply in {owner}/{repo}, to the '{id}' issue/pr")
    return jsonify({"status": "Reply created"}), 200


##
## INFO
##
@app.route('/', methods=["GET"])
def info():
	res = {
		"color": oreo.color,
		"image": oreo.image,
		"areas": oreo.areas
	}
	return jsonify(res), 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80, debug=True)
