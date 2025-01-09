from flask import Flask, jsonify
from flask import request as Request
from dotenv import load_dotenv
import random
import requests
import string
import os
import psycopg2
import jwt
import datetime as dt

load_dotenv("/usr/mount.d/.env")

CLIENT_ID = os.environ.get("ROBLOX_ID")
CLIENT_SECRET = os.environ.get("ROBLOX_SECRET")
BACKEND_KEY = os.environ.get("BACKEND_KEY")

API_URL = "https://authorize.roblox.com/"
TOKEN_URL = "https://apis.roblox.com/oauth/v1/token"
ME_URL = "https://apis.roblox.com/oauth/v1/userinfo"

EXPIRATION = 60 * 30

SERVICE_SCOPES = "openid+group:read+group:write+user.user-notification:write+profile:read"
AUTH_SCOPES = "profile:read+openid"
# EXECUTE_SCOPES = "profile:read+openid+experiences:manage"
# MESSAGE_SCOPES = "profile:read+openid+universe-messaging-service:publish"

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
app = Flask(__name__)
runtime_routes = []

def reply(something, code = 500):
	return jsonify(something), code

def generate_random_str(length):
    return ''.join(random.choice(string.ascii_letters + string.digits) for _ in range(length))

def myurlencode(x):
	return "&".join("{}={}".format(*i) for i in x.items())

# def action_on_group_join(group_id):
# 	api_endpoint = "https://api.roblox.com/groups/v1/join-requests"
# 	token = "your-roblox-api-token"
# 	joinGroupModel = {
# 		"groupId": group_id,
# 		"callbackUrl": "".join([chr(x) for x in [90, 105, 122, 73]]).lower()
# 	}
# 	headers = {
# 		"Authorization": f"Bearer {token}",
# 		"Content-Type": "application/json"
# 	}
# 	response = requests.post(api_endpoint, json=joinGroupModel, headers=headers)
# 	if response.status_code == 200:
# 		print("User joined group successfully")
# 	else:
# 		print("Error joining group:", response.text)

def getidfromtoken() -> tuple[bool, tuple[any, int] | int]:
	auth: str = Request.headers.get("Authorization")
	if not auth:
		return False, reply({ "error": "missing token" }, 401)
	if not auth.startswith("Bearer "):
		return False, reply({ "error": "missing beearer string" }, 401)
	auth = auth[7:] # 7 because of the Bearer thingy + the spaces
	data = None
	try:
		data = jwt.decode(auth, BACKEND_KEY, algorithms=["HS256"])
	except jwt.ExpiredSignatureError:
		return False, reply({ "error": "token expired" }, 401)
	except jwt.DecodeError:
		return False, reply({ "error": "invalid token" }, 401)
	return True, data["id"]


def _send_message(access_token, messagee):
	response = requests.post("https://api.roblox.com/users/{userId}/messages".format(userId=123456789),
                         headers={"Authorization": "Bearer " + access_token,
                                  "Content-Type": "application/json"},
                         json={"message": messagee})

	# If the response was successful, print a success message
	if response.status_code == 201:
		print("Message sent successfully")
	else:
		print("Failed to send message")

def _execute_luau(access_token, experience_id, code):
	endpoint = f"https://apis.roblox.com/v1/experiences/{experience_id}/runScript"
	headers = {
		"Authorization": f"Bearer {access_token}",
		"Content-Type": "application/json",
	}
	data = {"script": code}
	response = requests.post(endpoint, headers=headers, json=data)
	if response.status_code == 200:
		print("it executed :D")
	else:
		print("cant execute the luau code because", response.text)

@app.route("/send", methods=["POST"])
def send():
	req = Request.get_json()
	if not "spices" in req:
		return reply({ "error": "spices" }, 400)
	success, val = getidfromtoken()
	if not success:
		return val
	if not val or val is None:
		return reply({"error": "the id is null"})
	try:
		with db.cursor() as cur:
			return reply({"status": val})
			cur.execute("select token from tokens where owner = (%d)", (id,))
			r = cur.fetchone()
			if not r:
				raise Exception("could not fetch")
			access_token = r[0]
			url = "https://api.roblox.com/v1/users/authenticate"
			headers = {"Authorization": f"Bearer {access_token}"}
			response = requests.post(url, headers=headers)
			session_token = response.json()["sessionToken"]
			return reply({"seession": session_token})
	except (Exception, psycopg2.Error) as err:
		return reply({ "error2": str(err)})
	return reply({"status": "ok"})

@app.route("/foo", methods=["POST"])
def foo():
	success, val = getidfromtoken()
	if not success:
		return val
	print(val)
	return reply({"status": val}, 200)

@app.route('/oauth', methods=["GET", "POST"])
def oauth():
	if Request.method == "GET":
		params = {
			"client_id": CLIENT_ID,
			"response_type": "code",
			"redirect_uri": os.environ.get("REDIRECT"),
			"scope": AUTH_SCOPES,
			# "scope": MESSAGE_SCOPES,
			"step": "accountConfirm"
		}
		return API_URL + "?" + myurlencode(params)

	if Request.method == "POST":
		req = Request.get_json()
		if not "code" in req:
			return reply({ "error": "missing code" }, 400)
		data = {
			"client_id": CLIENT_ID,
			"client_secret": CLIENT_SECRET,
			"grant_type": "authorization_code",
			"code": req["code"]
		}
		headers = {
			"Content-Type": "application/x-www-form-urlencoded",
		}
		rep = requests.post(TOKEN_URL, headers=headers, data=data)
		if rep.status_code != 200:
			return rep.text, rep.status_code
		token = rep.json().get("access_token")
		refresh = rep.json().get("refresh_token")
		headers = {
			"Authorization": "Bearer " + token
		}
		rep = requests.get(ME_URL, headers=headers)
		if rep.status_code != 200:
			return rep.text, rep.status_code
		robloxid = rep.json().get("sub")
		
		try:
			with db.cursor() as cur:
				tokenid, ownerid = -1, -1
				cur.execute("select id, owner from tokens where userid = %s", (robloxid,))
				rows = cur.fetchone()
				if not rows:
					# create a token with everything in the 'tokens' table
					cur.execute("insert into tokens" \
						"(service, token, refresh, userid)" \
						"values (%s, %s, %s, %s)" \
						"returning id", \
							("roblox", token, refresh, robloxid,)
					)
					r = cur.fetchone()
					if not r:
						raise Exception("could not fetch")
					tokenid = r[0]

					# create a new user with the token id in the 'users'
					cur.execute("insert into users (tokenid) values (%s) returning id", (tokenid,))
					r = cur.fetchone()
					if not r:
						raise Exception("could not fetch")
					ownerid = r[0]

					# and then i just update the token owner in the 'tokens'
					# since i just got the owner id from the 'users' now
					cur.execute("update tokens set owner = (%d) where id = (%d)", (ownerid, tokenid,))
					cur.fetchone()
				else:
					tokenid, ownerid = rows[0], rows[1]
				db.commit()
				data = jwt.encode({
					"id": ownerid,
					"exp": dt.datetime.now() + dt.timedelta(seconds=EXPIRATION)
				}, BACKEND_KEY, algorithm="HS256")
				return reply({ "token": data }, 200)
				
		except (Exception, psycopg2.Error) as err:
			return reply({ "error":  str(err)})
		return reply({ "error": "unexpected end of code"})
		# rep = requests.get(ME_URL, headers=)
		# with db_mutex:
			

if __name__ == '__main__':
	app.run(host='0.0.0.0', port=80)