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
BACKEND_PORT = os.environ.get("BACKEND_PORT")

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
			cur.execute("SELECT token FROM tokens " \
			   "WHERE service = 'github' AND owner = %s", (
				   id,
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
	
	def create_area(self, name, type, title, spices):
		self.areas.append({
			"name": name,
			"type": type,
			"title": title,
			"spices": spices
		})


## #
##
## ROUTES
##
##

app = Flask(__name__)
# github color and image
oreo = NewOreo(
	service="github",
	color="#ff4500",
	image="https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png"
)


PERMISSIONS_REQUIRED = [
	"user",
	"repo",
	"public_repo",
	"write:issues"
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
	


		#########



		# get informations about the github user
		user_info_url = f"{GITHUB_API_URL}/user"
		headers.update({"Authorization": f"Bearer {github_access_token}"})

		user_info_response = requests.get(user_info_url, headers=headers)
		if user_info_response.status_code != 200:
			return jsonify({"error": "Failed to fetch user info"}), user_info_response.status_code
		
		user_info = user_info_response.json()
		github_user_id = str(user_info.get("id"))



		#########



		# data treatment
		area_bearer_token = retrieve_token(get_beared_token(request))
		area_user_id = area_bearer_token.get("id", None) if area_bearer_token else None
	
		# user is not logged in an area account
		if not area_bearer_token or not area_user_id:
			try:
				with db.cursor() as cur:
					cur.execute("SELECT owner FROM tokens " \
				 		"WHERE userid = %s AND service = %s", (
							github_user_id,
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
						area_user_id = cur.fetchone()[0]

						# create new token linked with the new area account
						cur.execute("INSERT INTO tokens " \
							"(service, token, refresh, userid, owner) " \
							"VALUES (%s, %s, %s, %s, %s)", (
								oreo.service,
		 						github_access_token,
								github_refresh_token,
								github_user_id,
								area_user_id,
							)
						)

						db.commit()
						return jsonify({ "token": generate_beared_token(area_user_id) }), 200
				
					# service account already linked with an area account: update token
					else:
						cur.execute(
								"UPDATE tokens " \
								"SET token = %s, refresh = %s " \
								"WHERE userid = %s AND service = %s " \
								"RETURNING owner", (
									github_access_token,
									github_refresh_token,
									github_user_id,
									oreo.service,
								)
						)
						
						area_user_id = cur.fetchone()[0]
			
						db.commit()
						return jsonify({ "token": generate_beared_token(area_user_id) }), 200
			except (Exception, psycopg2.Error) as err:
				return jsonify({ "error":  str(err)}), 400

		
		
		# user is already logged in an area account
		else:
			try:
				with db.cursor() as cur:
					cur.execute("SELECT owner FROM tokens " \
				 		"WHERE userid = %s AND service = %s", (
							 github_user_id,
							 oreo.service,
						)
					)
					rows = cur.fetchone()
			
					# service account already linked with an other area account: forbiden
					if rows and rows[0] != area_user_id:
						return jsonify({ "error": "forbiden: user already logged in an other account"}), 403
				
					# service account already linked with the same area account: update token
					elif rows and rows[0] == area_user_id:
						cur.execute("UPDATE tokens " \
				  			"SET token = %s, refresh = %s " \
							"WHERE userid = %s AND service = %s", (
								github_access_token,
								github_refresh_token,
								github_user_id,
								oreo.service,
							)
						)
						db.commit()
				
					# service account not linked with any area account: create new token
					else:
						cur.execute("INSERT INTO tokens " \
							"(service, token, refresh, userid, owner)" \
							"VALUES (%s, %s, %s, %s, %s)", (
								oreo.service,
								github_access_token,
								github_refresh_token,
								github_user_id,
								area_user_id,
							)
						)
						db.commit()
					
					return jsonify({ "token": generate_beared_token(area_user_id) }), 200
			except (Exception, psycopg2.Error) as err:
				return jsonify({ "error":  str(err)}), 400

		return jsonify({ "error": "unexpected end of code"}), 500


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
			"title": "The owner of the repository"
		},
		{
			"name": "repo",
			"type": "input",
			"title": "The repository name"
		},
		{
			"name": "title",
			"type": "input",
			"title": "The title of the issue (Markdown supported)"
		},
		{
			"name": "body",
			"type": "text",
			"title": "The body of the issue (Markdown supported)"
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
			"title": "The id of the issue or pull request"
		},
		{
			"name": "owner",
			"type": "input",
			"title": "The owner of the repository"
		},
		{
			"name": "repo",
			"type": "input",
			"title": "The repository name"
		},
		{
			"name": "body",
			"type": "text",
			"title": "The body of the reply (Markdown supported)"
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
## ACTIONS
##

# When you are assigned to an issue
ACTION_ASSIGNED_ISSUE = "assigned-issue"
oreo.create_area(
	ACTION_ASSIGNED_ISSUE,
	NewOreo.TYPE_ACTIONS,
	"When you are assigned to an issue",
	[ ]
)
@app.route('/assigned-issue', methods=["POST"])
def assigned_issue():
	app.logger.info("assigned-issue endpoint hit")

	# get data
	data = request.json
	if not data:
		return jsonify({"error": "Invalid JSON"}), 400

	area_user_id = data.get("id_user", 1)
	bridge = data.get("bridge")
	spices = data.get("spices")
	if not area_user_id or not bridge:
		return jsonify({"error": f"Missing required fields: 'user_id': {area_user_id}, 'spices': {spices}, 'bridge': {bridge}"}), 400


	with db.cursor() as cur:
		cur.execute("SELECT userid FROM tokens " \
			  "WHERE service = 'github' AND owner = %s", (
				  area_user_id,
			  )
		)
		rows = cur.fetchone()
		if not rows:
			return jsonify({"error": "User not found"}), 404
		github_user_id = rows[0]

		cur.execute("INSERT INTO micro_github" \
			  "(areauserid, userid, bridgeid, triggers) " \
			  "VALUES (%s, %s, %s, %s)", (
				  area_user_id,
				  github_user_id,
				  bridge,
				  ACTION_ASSIGNED_ISSUE
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200



##
## WEBHOOKS
##
@app.route('/webhook', methods=["POST"])
def webhook():
	app.logger.info("webhook endpoint hit")
	
	data = request.json
	action = data.get('action')
	if not action:
		app.logger.error("Invalid JSON")
		return jsonify({"error": "Invalid JSON"}), 400
	
	try:
		if action == "assigned":
			assignee = data.get('assignee', {})
			github_userid = assignee.get('id')

			if not github_userid:
				return jsonify({"error": "Invalid JSON"}), 400

			with db.cursor() as cur:
				cur.execute("SELECT bridgeid, areauserid FROM micro_github " \
					"WHERE userid = %s AND triggers = %s", (
						github_userid,
						ACTION_ASSIGNED_ISSUE
					)
				)

				# check if the user has an action assigned
				rows = cur.fetchall()
				if not rows:
					return jsonify({"status": "ok"}), 200
				
				# get the bridge id
				for row in rows:
					bridge = row[0]
					areauserid = row[1]
					requests.put(
						f"http://backend:{BACKEND_PORT}/api/orchestrator",
						json={
							"bridge": bridge,
							"userid": areauserid,
							"ingredients": {}
						}
					)
				return jsonify({"status": "ok"}), 200
				
	except (Exception, psycopg2.Error) as err:
		app.logger.error(str(err))
		return jsonify({"error": str(err)}), 400

	return jsonify({"status": "ok"}), 200



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
