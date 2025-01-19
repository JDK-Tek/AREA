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

WEATHER = {
    "Clear": {"min": 800, "max": 800},
    "Clouds": {"min": 801, "max": 804},
    "Drizzle": {"min": 300, "max": 321},
    "Rain": {"min": 500, "max": 531},
    "Thunderstorm": {"min": 200, "max": 232},
    "Snow": {"min": 600, "max": 622},
    "Mist": {"min": 701, "max": 701},
    "Smoke": {"min": 711, "max": 711},
    "Haze": {"min": 721, "max": 721},
    "Dust": {"min": 731, "max": 731},
    "Fog": {"min": 741, "max": 741},
    "Dust": {"min": 761, "max": 761},
    "Ash": {"min": 762, "max": 762},
    "Squall": {"min": 771, "max": 771},
    "Tornado": {"min": 781, "max": 781},
    "Extreme": {"min": 900, "max": 906},
    "Additional": {"min": 951, "max": 962}
}


##
## ACTIONS
##

ACTION_WEATHER_AT_CITY = "when-weather-at-city"
oreo.create_area(
	ACTION_WEATHER_AT_CITY,
	NewOreo.TYPE_ACTIONS,
	"When a specific weather is detected in a city",
	[
		{
			"name": "city",
			"type": "input",
			"title": "The name of the city"
		},
		{
			"name": "weather",
			"type": "dropdown",
			"title": "The weater to detect",
			"extra": [key for key in WEATHER.keys()]
		},
		{
			"name": "unit",
			"type": "dropdown",
			"title": "The weather unit",
			"extra": ["kelvin", "imperial", "metric"]
		}
	]
)
@app.route(f'/{ACTION_WEATHER_AT_CITY}', methods=["POST"])
def weather_at_city():
	app.logger.info(f"{ACTION_WEATHER_AT_CITY} endpoint hit")

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
				  ACTION_WEATHER_AT_CITY,
				  json.dumps(spices),
				  ""
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_PRESSURE_AT_CITY = "when-pressure-at-city"
oreo.create_area(
	ACTION_PRESSURE_AT_CITY,
	NewOreo.TYPE_ACTIONS,
	"When a specific pressure is detected in a city",
	[
		{
			"name": "city",
			"type": "input",
			"title": "The name of the city"
		},
		{
			"name": "data",
			"type": "number",
			"title": "The pressure to compare to",
		},
		{
			"name": "comparaison",
			"type": "dropdown",
			"title": "The comparaison",
			"extra": ["inferior", "superior", "inferior or equal", "superior or equal"]
		},
		{
			"name": "unit",
			"type": "dropdown",
			"title": "The weather unit",
			"extra": ["kelvin", "imperial", "metric"]
		}
	]
)
@app.route(f'/{ACTION_PRESSURE_AT_CITY}', methods=["POST"])
def pressure_at_city():
	app.logger.info(f"{ACTION_PRESSURE_AT_CITY} endpoint hit")

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
				  ACTION_PRESSURE_AT_CITY,
				  json.dumps(spices),
				  ""
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_HUMIDITY_AT_CITY = "when-humidity-at-city"
oreo.create_area(
	ACTION_HUMIDITY_AT_CITY,
	NewOreo.TYPE_ACTIONS,
	"When a specific humidity is detected in a city",
	[
		{
			"name": "city",
			"type": "input",
			"title": "The name of the city"
		},
		{
			"name": "data",
			"type": "number",
			"title": "The humidity to compare to",
		},
		{
			"name": "comparaison",
			"type": "dropdown",
			"title": "The comparaison",
			"extra": ["inferior", "superior", "inferior or equal", "superior or equal"]
		},
		{
			"name": "unit",
			"type": "dropdown",
			"title": "The weather unit",
			"extra": ["kelvin", "imperial", "metric"]
		}
	]
)
@app.route(f'/{ACTION_HUMIDITY_AT_CITY}', methods=["POST"])
def humidity_at_city():
	app.logger.info(f"{ACTION_HUMIDITY_AT_CITY} endpoint hit")

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
				  ACTION_HUMIDITY_AT_CITY,
				  json.dumps(spices),
				  ""
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_TEMP_AT_CITY = "when-temp-at-city"
oreo.create_area(
	ACTION_TEMP_AT_CITY,
	NewOreo.TYPE_ACTIONS,
	"When a specific temp is detected in a city",
	[
		{
			"name": "city",
			"type": "input",
			"title": "The name of the city"
		},
		{
			"name": "data",
			"type": "number",
			"title": "The temp to compare to",
		},
		{
			"name": "comparaison",
			"type": "dropdown",
			"title": "The comparaison",
			"extra": ["inferior", "superior", "inferior or equal", "superior or equal"]
		},
		{
			"name": "unit",
			"type": "dropdown",
			"title": "The weather unit",
			"extra": ["kelvin", "imperial", "metric"]
		}
	]
)
@app.route(f'/{ACTION_TEMP_AT_CITY}', methods=["POST"])
def temp_at_city():
	app.logger.info(f"{ACTION_TEMP_AT_CITY} endpoint hit")

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
				  ACTION_TEMP_AT_CITY,
				  json.dumps(spices),
				  ""
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_VISIBILITY_AT_CITY = "when-visibility-at-city"
oreo.create_area(
	ACTION_VISIBILITY_AT_CITY,
	NewOreo.TYPE_ACTIONS,
	"When a specific visibility is detected in a city",
	[
		{
			"name": "city",
			"type": "input",
			"title": "The name of the city"
		},
		{
			"name": "data",
			"type": "number",
			"title": "The visibility to compare to",
		},
		{
			"name": "comparaison",
			"type": "dropdown",
			"title": "The comparaison",
			"extra": ["inferior", "superior", "inferior or equal", "superior or equal"]
		},
		{
			"name": "unit",
			"type": "dropdown",
			"title": "The weather unit",
			"extra": ["kelvin", "imperial", "metric"]
		}
	]
)
@app.route(f'/{ACTION_VISIBILITY_AT_CITY}', methods=["POST"])
def visibility_at_city():
	app.logger.info(f"{ACTION_VISIBILITY_AT_CITY} endpoint hit")

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
				  ACTION_VISIBILITY_AT_CITY,
				  json.dumps(spices),
				  ""
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_WIND_AT_CITY = "when-wind-at-city"
oreo.create_area(
	ACTION_WIND_AT_CITY,
	NewOreo.TYPE_ACTIONS,
	"When a specific wind speed is detected in a city",
	[
		{
			"name": "city",
			"type": "input",
			"title": "The name of the city"
		},
		{
			"name": "data",
			"type": "number",
			"title": "The wind speed to compare to",
		},
		{
			"name": "comparaison",
			"type": "dropdown",
			"title": "The comparaison",
			"extra": ["inferior", "superior", "inferior or equal", "superior or equal"]
		},
		{
			"name": "unit",
			"type": "dropdown",
			"title": "The weather unit",
			"extra": ["kelvin", "imperial", "metric"]
		}
	]
)
@app.route(f'/{ACTION_WIND_AT_CITY}', methods=["POST"])
def wind_at_city():
	app.logger.info(f"{ACTION_WIND_AT_CITY} endpoint hit")

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
				  ACTION_WIND_AT_CITY,
				  json.dumps(spices),
				  ""
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200



WEATHER_DATA_COMP = [
	(ACTION_WIND_AT_CITY, "wind/speed"),
	(ACTION_VISIBILITY_AT_CITY, "visibility"),
	(ACTION_TEMP_AT_CITY, "main/temp"),
	(ACTION_HUMIDITY_AT_CITY, "main/humidity"),
	(ACTION_PRESSURE_AT_CITY, "main/pressure"),
]

def webhook():
	print("Starting newsapi webhook", file=sys.stderr)

	while True:
		print("Checking for new weather", file=sys.stderr)
		time.sleep(10)

		with db.cursor() as cur:
			cur.execute("SELECT userid, bridgeid, triggers, spices, last_weather FROM micro_openweather "
			   "WHERE triggers = %s", (ACTION_WEATHER_AT_CITY,))
	
			rows = cur.fetchall()
			print(f"Checking weather for {len(rows)} users", file=sys.stderr)
	
			for row in rows:
				userid, bridgeid, triggers, spices, last_weather = row
				spices = json.loads(spices)
	
				unit = spices.get("unit")
				city = spices.get("city")
				weather = spices.get("weather")
				if not city or not weather or not unit:
					print(f"Missing required fields: 'city': {city}, 'weather': {weather}, 'unit': {unit}", file=sys.stderr)
					continue

				unit_url = f"&units={quote(unit)}" if (unit != "kelvin") else ""
				url = f"{API_URL}/weather?" \
					f"appid={API_KEY}" \
					f"&q={quote(city)}"	\
					f"{unit_url}" \
	
				print(f"Requesting weather for {city} at {url}", file=sys.stderr)

				res = requests.get(url)
				if res.status_code != 200:
					print(f"Error while fetching weather for {city}", file=sys.stderr)
					continue

				print(f"Got weather for {city}", file=sys.stderr)
				# update the last weather
				cur.execute("UPDATE micro_openweather "
					"SET last_weather = %s "
					"WHERE userid = %s AND bridgeid = %s AND triggers = %s AND spices = %s", (
					weather,
					userid,
					bridgeid,
					triggers,
					json.dumps(spices)
				))

				db.commit()

				data = res.json()
				if not data:
					continue

				weather_id = data["weather"][0]["id"]
				if not weather_id:
					continue

				if (last_weather != weather
				and weather_id >= WEATHER[weather]["min"]
				and weather_id <= WEATHER[weather]["max"]):

					res = requests.put(
						f"http://backend:{BACKEND_PORT}/api/orchestrator",
						json={
							"bridge": bridgeid,
							"userid": userid,
							"ingredients": {
								"weather": str(data.get("weather")[0].get("description")),
								"temperature": str(data.get("main").get("temp")),
								"feels_like": str(data.get("main").get("feels_like")),
								"humidity": str(data.get("main").get("humidity")),
								"pressure": str(data.get("main").get("pressure")),
								"wind_speed": str(data.get("wind").get("speed")),
								"wind_deg": str(data.get("wind").get("deg")),
								"clouds": str(data.get("clouds").get("all")),
								"visibility": str(data.get("visibility")),
							}
						}
					)
					print(f"New weather detected: {weather} in {city}: {res.json()}", file=sys.stderr)

			for (weather_data, path_comp) in WEATHER_DATA_COMP:
				print(f"Checking data weather for {weather_data}", file=sys.stderr)
		
				cur.execute("SELECT userid, bridgeid, triggers, spices, last_weather FROM micro_openweather "
				   "WHERE triggers = %s", (weather_data,))
		
				rows = cur.fetchall()
				print(f"Checking weather for {len(rows)} users", file=sys.stderr)
		
				for row in rows:
					userid, bridgeid, triggers, spices, last_weather = row
					spices = json.loads(spices)
		
					unit = spices.get("unit")
					city = spices.get("city")
					value = spices.get("data")
					comparaison = spices.get("comparaison")
					if not city or not value or not comparaison or not unit:
						print(f"Missing required fields: 'city': {city}, 'value': {value}, 'comparaison': {comparaison}, 'unit': {unit}", file=sys.stderr)
						continue

					unit_url = f"&units={quote(unit)}" if (unit != "kelvin") else ""
					url = f"{API_URL}/weather?" \
						f"appid={API_KEY}" \
						f"&q={quote(city)}"	\
						f"{unit_url}" \
	
					print(f"Requesting weather for {city} at {url}", file=sys.stderr)
		
					res = requests.get(url)
					if res.status_code != 200:
						print(f"Error while fetching weather for {city}", file=sys.stderr)
						continue
		
					print(f"Got weather for {city}", file=sys.stderr)
		
					data = res.json()
					if not data:
						continue
		
					path_splited = path_comp.split("/")
					value_to_compare = data
					for key in path_splited:
						if isinstance(value_to_compare, dict) and key in value_to_compare:
							value_to_compare = value_to_compare[key]
						else:
							value_to_compare = None
							break

					if (last_weather != "true" and (
					(comparaison == "inferior" and value_to_compare < value)
					or (comparaison == "superior" and value_to_compare > value)
					or (comparaison == "inferior or equal" and value_to_compare <= value)
					or (comparaison == "superior or equal" and value_to_compare >= value))):
		
						cur.execute("UPDATE micro_openweather "
							"SET last_weather = %s "
							"WHERE userid = %s AND bridgeid = %s AND triggers = %s AND spices = %s", (
							"true",
							userid,
							bridgeid,
							triggers,
							json.dumps(spices)
						))
			
						db.commit()

						res = requests.put(
							f"http://backend:{BACKEND_PORT}/api/orchestrator",
							json={
								"bridge": bridgeid,
								"userid": userid,
								"ingredients": {
									"weather": str(data.get("weather")[0].get("description")),
									"temperature": str(data.get("main").get("temp")),
									"feels_like": str(data.get("main").get("feels_like")),
									"humidity": str(data.get("main").get("humidity")),
									"pressure": str(data.get("main").get("pressure")),
									"wind_speed": str(data.get("wind").get("speed")),
									"wind_deg": str(data.get("wind").get("deg")),
									"clouds": str(data.get("clouds").get("all")),
									"visibility": str(data.get("visibility")),
								}
							}
						)
						print(f"New weather detected: {value} in {city}: {res.json()}", file=sys.stderr)
					elif (last_weather != "true"):
						# update the last weather
						cur.execute("UPDATE micro_openweather "
							"SET last_weather = %s "
							"WHERE userid = %s AND bridgeid = %s AND triggers = %s AND spices = %s", (
							"false",
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
