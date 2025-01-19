import os
import sys
import jwt
import requests
import threading
import psycopg2
import time
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

API_KEY = os.environ.get("API_KEY")
BACKEND_PORT = os.environ.get("BACKEND_PORT")

API_URL = "https://calendarific.com/api/v2"


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
# calendarific color and image
oreo = NewOreo(
	service="calendarific",
	color="#35008B",
	image="/assets/calendarific.webp",
	oauth=False
)


##
## ACTIONS
##


SUBSCRIBE_FRENCH_EVENT = "when-french-events"
oreo.create_area(
	SUBSCRIBE_FRENCH_EVENT,
	NewOreo.TYPE_ACTIONS,
	"When a French events is coming",
	[]
)
@app.route(f'/{SUBSCRIBE_FRENCH_EVENT}', methods=["POST"])
def new_article_general():
	app.logger.info(f"{SUBSCRIBE_FRENCH_EVENT} endpoint hit")

	# get data
	data = request.json
	if not data:
		return jsonify({"error": "Invalid JSON"}), 400

	userid = data.get("userid")
	bridge = data.get("bridge")
	spices = data.get("spices", {})
	if not userid or not bridge:
		return jsonify({"error": f"Missing required fields: 'userid': {userid}, 'spices': {spices}, 'bridge': {bridge}"}), 400

	with db.cursor() as cur:
		cur.execute("INSERT INTO micro_calendarific" \
			  "(userid, bridgeid, triggers, spices, events) " \
			  "VALUES (%s, %s, %s, %s, %s)", (
				  userid,
				  bridge,
				  SUBSCRIBE_FRENCH_EVENT,
				  json.dumps(spices),
				  ""
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

EVENTS = [
	(SUBSCRIBE_FRENCH_EVENT, "FR"),
]

##
## WEBHOOKS
##

def webhook():
	print("Starting calendarific webhook", file=sys.stderr)

	while True:
		with db.cursor() as cur:
			cur.execute("SELECT userid, bridgeid, triggers, spices, events FROM micro_calendarific")
			
			rows = cur.fetchall()
			if not rows:
				time.sleep(10)
				continue

			year = time.strftime("%Y")
			for row in rows:
				userid, bridgeid, triggers, spices, last_event = row

				for (events, country) in EVENTS:
					if events == triggers:
						url = f"{API_URL}/holidays?api_key={API_KEY}&country={country}&year={year}"
						res = requests.get(url)
						if res.status_code != 200:
							print(f"Error: {res.status_code}", file=sys.stderr)
							continue
						
						data = res.json()
						if not data or not data.get("response") or not data.get("response").get("holidays"):
							continue
						holidays = data.get("response").get("holidays")
						todays_date = time.strftime("%Y-%m-%d")
						# todays_date = "2025-03-01" # for testing
						for holiday in holidays:
							if (holiday.get("date").get("iso") == todays_date and last_event == ""):
								cur.execute("UPDATE micro_calendarific SET events=%s WHERE userid=%s AND bridgeid=%s", (str(holiday.get("name")), userid, bridgeid))
								db.commit()

								res = requests.put(
									f"http://backend:{BACKEND_PORT}/api/orchestrator",
									json={
										"bridge": bridgeid,
										"userid": userid,
										"ingredients": {
											"name": str(holiday.get("name")),
											"description": str(holiday.get("description")),
											"type": str(holiday.get("primary_type"))
										}
									}
								)
								print(f"Sent reaction: {res.json()}", file=sys.stderr)

								break
							elif (holiday.get("date").get("iso") != todays_date and last_event == holiday.get("name")):
								cur.execute("UPDATE micro_calendarific SET events=%s WHERE userid=%s AND bridgeid=%s", ("", userid, bridgeid))
								db.commit()

		time.sleep(10)


##
## INFO
##
@app.route('/', methods=["GET"])
def info():
	res = {
		"color": oreo.color,
		"image": oreo.image,
		"oauth": oreo.oauth,
		"areas": oreo.areas
	}
	return jsonify(res), 200


if __name__ == '__main__':
    threading.Thread(target=webhook, daemon=True).start()
    app.run(host='0.0.0.0', port=80, debug=True)
