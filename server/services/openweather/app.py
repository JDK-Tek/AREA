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
# newsapi color and image
oreo = NewOreo(
	service="openweather",
	color="#E96E4C",
	image="/assets/openweather.png",
	oauth=False
)

## #
##
## ENV VARIABLES
##
##

API_KEY = os.environ.get("API_KEY")
BACKEND_PORT = os.environ.get("BACKEND_PORT")

API_URL = "https://api.openweathermap.org/data/2.5"

METEO = {
	"Clear": {"min": 800, "max": 800},
	"Clouds": {"min": 801, "max": 804},
	"Drizzle": {"min": 300, "max": 321},
	"Rain": {"min": 500, "max": 531},
	"Thunderstorm": {"min": 200, "max": 232},
	"Snow": {"min": 600, "max": 622},
	"Extreme": {"min": 900, "max": 906},
	"Additional": {"min": 951, "max": 962}
}

##
## ACTIONS
##

ACTION_METEO_AT_CITY = "when-meteo-at-city"
oreo.create_area(
	ACTION_METEO_AT_CITY,
	NewOreo.TYPE_ACTIONS,
	"When a specific meteo is detected in a city",
	[
		{
			"name": "city",
			"type": "input",
			"title": "The name of the city"
		},
		{
			"name": "meteo",
			"type": "dropdown",
			"title": "The meteo to detect",
			"extra": [key for key in METEO.keys()]
		},
	]
)
@app.route(f'/{ACTION_METEO_AT_CITY}', methods=["POST"])
def new_article_general():
	app.logger.info(f"{ACTION_METEO_AT_CITY} endpoint hit")

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
		cur.execute("INSERT INTO micro_openweather" \
			  "(userid, bridgeid, triggers, spices, last_weather) " \
			  "VALUES (%s, %s, %s, %s, %s)", (
				  userid,
				  bridge,
				  ACTION_METEO_AT_CITY,
				  json.dumps(spices),
				  ""
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200



def webhook():
	print("Starting newsapi webhook", file=sys.stderr)

	while True:
		print("Checking for new meteo", file=sys.stderr)
		time.sleep(10)

		with db.cursor() as cur:
			cur.execute("SELECT userid, bridgeid, triggers, spices, last_weather FROM micro_openweather "
			   "WHERE triggers = %s", (ACTION_METEO_AT_CITY,))
	
			for row in cur.fetchall():
				userid, bridgeid, triggers, spices, last_weather = row
				spices = json.loads(spices)

				print(f"Checking meteo for {userid} user", file=sys.stderr)
				# get the weather
				city = spices.get("city")
				meteo = spices.get("meteo")
				if not city or not meteo:
					continue

				url = f"{API_URL}/weather?q={quote(city)}&appid={API_KEY}"
				res = requests.get(url)
				if res.status_code != 200:
					print(f"Error while fetching weather for {city}", file=sys.stderr)
					continue

				weather = res.json()
				if not weather:
					continue

				weather_id = weather["weather"][0]["id"]
				if not weather_id:
					continue

				if (last_weather != meteo
				and weather_id >= METEO[meteo]["min"]
				and weather_id <= METEO[meteo]["max"]):
					if last_weather != meteo:

						res = requests.put(
							f"http://backend:{BACKEND_PORT}/api/orchestrator",
							json={
								"bridge": bridgeid,
								"userid": userid,
								"ingredients": {
									# "meteo": weather.get("weather")[0].get("main"),
									# "temperature": weather.get("main").get("temp"),
									# "feels_like": weather.get("main").get("feels_like"),
									# "humidity": weather.get("main").get("humidity"),
									# "pressure": weather.get("main").get("pressure"),
									# "wind_speed": weather.get("wind").get("speed"),
									# "wind_deg": weather.get("wind").get("deg"),
									# "clouds": weather.get("clouds").get("all"),
									# "visibility": weather.get("visibility"),
								}
							}
						)
						print(f"New meteo detected: {meteo} in {city}: {res.json()}", file=sys.stderr)

					# update the last weather
					cur.execute("UPDATE micro_openweather "
						"SET last_weather = %s "
						"WHERE userid = %s AND bridgeid = %s AND triggers = %s AND spices = %s", (
						meteo,
						userid,
						bridgeid,
						triggers,
						json.dumps(spices)
					))

					db.commit()

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
