import os
import sys
import jwt
import requests
import psycopg2
import datetime as dt
from dotenv import load_dotenv
from flask import Flask, json, jsonify, request
from urllib.parse import quote

sys.stdout.reconfigure(line_buffering=True)
load_dotenv("/usr/mount.d/.env")

## #
##
## ENV VARIABLES
##
##

ACCESS_TOKEN = os.environ.get("ACCESS_TOKEN")
REFRESH_TOKEN = os.environ.get("REFRESH_TOKEN")

REDIRECT_URI = os.environ.get("REDIRECT")
BACKEND_PORT = os.environ.get("BACKEND_PORT")

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
			   "WHERE service = 'slack' AND owner = %s", (
				   id,
				)
			)
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


	def __init__(self, service, color, image, oauth):
		self.service = service
		self.color = color
		self.image = image
		self.oauth = oauth
		self.areas = []
	
	def create_area(self, name, type, title, spices):
		self.areas.append({
			"name": name,
			"type": type,
			"description": title,
			"spices": spices
		})


## #
##
## ROUTES
##
##

app = Flask(__name__)
# slack color and image
oreo = NewOreo(
	service="slack",
	color="#5D0358",
	image="/assets/slack.webp",
	oauth=True
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

##
## INFO
##
@app.route('/', methods=["GET"])
def info():
	res = {
		"color": oreo.color,
		"image": oreo.image,
		"oauth": oreo.oauth,
		"areas": oreo.areas,
	}
	return jsonify(res), 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80, debug=True)
