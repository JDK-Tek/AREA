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

API_URL = "https://newsapi.org/v2"


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


	def __init__(self, service, color, image):
		self.service = service
		self.color = color
		self.image = image
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
	service="newsapi",
	color="#d4bd13",
	image="/assets/services/newsapi.png"
)


##
## ACTIONS
##

COUNTRIES = [
	"fr",
	"us",
	"pt",
	"ae",
	"ar",
	"at",
	"au",
	"be",
	"bg",
	"br",
	"ca",
	"ch",
	"cn",
	"co",
	"cu",
	"cz",
	"de",
	"eg",
	"gb",
	"gr",
	"hk",
	"hu",
	"id",
	"ie",
	"il",
	"in",
	"it",
	"jp",
	"kr",
	"lt",
	"lv",
	"ma",
	"mx",
	"my",
	"ng",
	"nl",
	"no",
	"nz",
	"ph",
	"pl",
	"ro",
	"rs",
	"ru",
	"sa",
	"se",
	"sg",
	"si",
	"sk",
	"th",
	"tr",
	"tw",
	"ua",
	"ve",
	"za"
]

# country
# key word

ACTION_NEW_ART_COUNTRY = "new-art-country"
oreo.create_area(
	ACTION_NEW_ART_COUNTRY,
	NewOreo.TYPE_ACTIONS,
	"When a new article is published in a specific country",
	[
		{
			"name": "country",
			"type": "dropdown",
			"extra": COUNTRIES,
			"title": "Select the country"
		}
	]
)
@app.route(f'/{ACTION_NEW_ART_COUNTRY}', methods=["POST"])
def new_article_country():
	app.logger.info(f"{ACTION_NEW_ART_COUNTRY} endpoint hit")

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
		cur.execute("INSERT INTO micro_newsapi" \
			  "(userid, bridgeid, triggers, spices) " \
			  "VALUES (%s, %s, %s, %s)", (
				  userid,
				  bridge,
				  ACTION_NEW_ART_COUNTRY,
				  json.dumps(spices)
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200


ACTION_NEW_ART_TECH = "new-art-technology"
oreo.create_area(
	ACTION_NEW_ART_TECH,
	NewOreo.TYPE_ACTIONS,
	"When a new article is published in Technology category",
	[]
)
@app.route(f'/{ACTION_NEW_ART_TECH}', methods=["POST"])
def new_article_tech():
	app.logger.info(f"{ACTION_NEW_ART_TECH} endpoint hit")

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
		cur.execute("INSERT INTO micro_newsapi" \
			  "(userid, bridgeid, triggers, spices) " \
			  "VALUES (%s, %s, %s, %s)", (
				  userid,
				  bridge,
				  ACTION_NEW_ART_TECH,
				  json.dumps(spices)
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_NEW_ART_BUSINESS = "new-art-business"
oreo.create_area(
	ACTION_NEW_ART_BUSINESS,
	NewOreo.TYPE_ACTIONS,
	"When a new article is published in Business category",
	[]
)
@app.route(f'/{ACTION_NEW_ART_BUSINESS}', methods=["POST"])
def new_article_business():
	app.logger.info(f"{ACTION_NEW_ART_BUSINESS} endpoint hit")

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
		cur.execute("INSERT INTO micro_newsapi" \
			  "(userid, bridgeid, triggers, spices) " \
			  "VALUES (%s, %s, %s, %s)", (
				  userid,
				  bridge,
				  ACTION_NEW_ART_BUSINESS,
				  json.dumps(spices)
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_NEW_ART_ENTERTAINEMENT = "new-art-entertainment"
oreo.create_area(
	ACTION_NEW_ART_ENTERTAINEMENT,
	NewOreo.TYPE_ACTIONS,
	"When a new article is published in Entertainment category",
	[]
)
@app.route(f'/{ACTION_NEW_ART_ENTERTAINEMENT}', methods=["POST"])
def new_article_entertainment():
	app.logger.info(f"{ACTION_NEW_ART_ENTERTAINEMENT} endpoint hit")

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
		cur.execute("INSERT INTO micro_newsapi" \
			  "(userid, bridgeid, triggers, spices) " \
			  "VALUES (%s, %s, %s, %s)", (
				  userid,
				  bridge,
				  ACTION_NEW_ART_ENTERTAINEMENT,
				  json.dumps(spices)
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_NEW_ART_SCIENCE = "new-art-science"
oreo.create_area(
	ACTION_NEW_ART_SCIENCE,
	NewOreo.TYPE_ACTIONS,
	"When a new article is published in Science category",
	[]
)
@app.route(f'/{ACTION_NEW_ART_SCIENCE}', methods=["POST"])
def new_article_science():
	app.logger.info(f"{ACTION_NEW_ART_SCIENCE} endpoint hit")

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
		cur.execute("INSERT INTO micro_newsapi" \
			  "(userid, bridgeid, triggers, spices) " \
			  "VALUES (%s, %s, %s, %s)", (
				  userid,
				  bridge,
				  ACTION_NEW_ART_SCIENCE,
				  json.dumps(spices)
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_NEW_ART_HEALTH = "new-art-health"
oreo.create_area(
	ACTION_NEW_ART_HEALTH,
	NewOreo.TYPE_ACTIONS,
	"When a new article is published in Health category",
	[]
)
@app.route(f'/{ACTION_NEW_ART_HEALTH}', methods=["POST"])
def new_article_health():
	app.logger.info(f"{ACTION_NEW_ART_HEALTH} endpoint hit")

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
		cur.execute("INSERT INTO micro_newsapi" \
			  "(userid, bridgeid, triggers, spices) " \
			  "VALUES (%s, %s, %s, %s)", (
				  userid,
				  bridge,
				  ACTION_NEW_ART_HEALTH,
				  json.dumps(spices)
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_NEW_ART_SPORTS = "new-art-sports"
oreo.create_area(
	ACTION_NEW_ART_SPORTS,
	NewOreo.TYPE_ACTIONS,
	"When a new article is published in Sports category",
	[]
)
@app.route(f'/{ACTION_NEW_ART_SPORTS}', methods=["POST"])
def new_article_sports():
	app.logger.info(f"{ACTION_NEW_ART_SPORTS} endpoint hit")

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
		cur.execute("INSERT INTO micro_newsapi" \
			  "(userid, bridgeid, triggers, spices) " \
			  "VALUES (%s, %s, %s, %s)", (
				  userid,
				  bridge,
				  ACTION_NEW_ART_SPORTS,
				  json.dumps(spices)
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200

ACTION_NEW_ART_GENERAL = "new-art-general"
oreo.create_area(
	ACTION_NEW_ART_GENERAL,
	NewOreo.TYPE_ACTIONS,
	"When a new article is published in General category",
	[]
)
@app.route(f'/{ACTION_NEW_ART_GENERAL}', methods=["POST"])
def new_article_general():
	app.logger.info(f"{ACTION_NEW_ART_GENERAL} endpoint hit")

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
		cur.execute("INSERT INTO micro_newsapi" \
			  "(userid, bridgeid, triggers, spices) " \
			  "VALUES (%s, %s, %s, %s)", (
				  userid,
				  bridge,
				  ACTION_NEW_ART_GENERAL,
				  json.dumps(spices)
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200



CATEGORIES = {
	ACTION_NEW_ART_TECH: "technology",
	ACTION_NEW_ART_BUSINESS: "buisness",
	ACTION_NEW_ART_ENTERTAINEMENT: "entertainment",
	ACTION_NEW_ART_SCIENCE: "science",
	ACTION_NEW_ART_HEALTH: "health",
	ACTION_NEW_ART_SPORTS: "sports",
	ACTION_NEW_ART_GENERAL: "general"
}

##
## WEBHOOKS
##

def get_news_from_category(category):
	res = requests.get(f"{API_URL}/top-headlines?category={category}&apiKey={API_KEY}")
	if res.status_code != 200:
		print(f"Error getting news from {category} category: {res.text}", file=sys.stderr)
		return None

	data = res.json()
	articles = data.get("articles", [])

	if not articles:
		return []

	res_articles = []
	for article in articles:
		url = article.get("url")

		with db.cursor() as cur:
			cur.execute("SELECT url FROM micro_newsapi_articles "
				"WHERE url = %s", (url,))

			exist = cur.fetchone()
			if exist is None:
				res_articles.append(article)

				cur.execute("INSERT INTO micro_newsapi_articles "
					"(url) "
					"VALUES (%s)", (article.get("url"),)
				)
				db.commit()

			else:
				print(f"Article {url} already exists in the database", file=sys.stderr)
	
	print(f"Returning {len(res_articles)} new articles", file=sys.stderr)
	return res_articles

def get_news_from_country(country):
	res = requests.get(f"{API_URL}/top-headlines?country={country}&apiKey={API_KEY}")
	if res.status_code != 200:
		print(f"Error getting news from {country} country: {res.text}", file=sys.stderr)
		return None

	data = res.json()
	articles = data.get("articles", [])

	if not articles:
		return []

	res_articles = []
	for article in articles:
		url = article.get("url")

		with db.cursor() as cur:
			cur.execute("SELECT url FROM micro_newsapi_articles "
				"WHERE url = %s", (url,))

			exist = cur.fetchone()
			if exist is None:
				res_articles.append(article)

				cur.execute("INSERT INTO micro_newsapi_articles "
					"(url) "
					"VALUES (%s)", (article.get("url"),)
				)
				db.commit()

			else:
				print(f"Article {url} already exists in the database", file=sys.stderr)
	
	print(f"Returning {len(res_articles)} new articles", file=sys.stderr)
	return res_articles

def get_subbed_countries():
	with db.cursor() as cur:
		cur.execute("SELECT spices FROM micro_newsapi "
			"WHERE triggers = %s", (ACTION_NEW_ART_COUNTRY,))
		countries = cur.fetchall()

		if not countries or len(countries) == 0:
			return []

		countries = [json.loads(c[0]).get("country") for c in countries]
		return countries

def webhook():
	print("Starting newsapi webhook", file=sys.stderr)

	while True:

		time.sleep(60)

		# get all subscribed users to the action
		with db.cursor() as cur:
			cur.execute("SELECT * FROM micro_newsapi")
			users_sub = cur.fetchall()

			if not users_sub or len(users_sub) == 0:
				print("No users subscribed to the newsapi actions", file=sys.stderr)
				continue
			
			print(f"Found {len(users_sub)} users subscribed to the newsapi actions", file=sys.stderr)


			# check if there are users subscribed to a category news
			for (action, category) in CATEGORIES.items():

				cur.execute("SELECT userid, bridgeid FROM micro_newsapi "
				"WHERE triggers = %s", (action,))
				users_sub = cur.fetchall()

				# check if there are users subscribed to the technology news
				if users_sub and len(users_sub) > 0:
					print(f"Found {len(users_sub)} users subscribed to the {category} news", file=sys.stderr)
					articles = get_news_from_category(f"{category}")

					if articles and len(articles) > 0:
						print(f"Found {len(articles)} new articles in the {category} category", file=sys.stderr)

						# trigger the action for each user
						for article in articles:
							for user in users_sub:
								userid = user[0]
								bridge = user[1]


								requests.put(
									f"http://backend:{BACKEND_PORT}/api/orchestrator",
									json={
										"bridge": bridge,
										"userid": userid,
										"ingredients": {
											"source": article.get("source").get("name"),
											"author": article.get("author"),
											"publishedAt": article.get("publishedAt"),
											"title": article.get("title"),
											"description": article.get("description"),
											"url": article.get("url"),
											"content": article.get("content")
										}
									}
								)
		
					else:
						print(f"No new articles in the {category} category", file=sys.stderr)
			
			# check if there are users subscribed to the country news
			cur.execute("SELECT userid, bridgeid FROM micro_newsapi "
				"WHERE triggers = %s", (ACTION_NEW_ART_COUNTRY,))
			users_sub = cur.fetchall()

			if users_sub and len(users_sub) > 0:
				print(f"Found {len(users_sub)} users subscribed to the country news", file=sys.stderr)
				subbed_countries = get_subbed_countries()
		
				for country in subbed_countries:
					articles = get_news_from_country(country)

					if articles and len(articles) > 0:
						print(f"Found {len(articles)} new articles in the {country} country", file=sys.stderr)

						# trigger the action for each user
						for article in articles:
							for user in users_sub:
								userid = user[0]
								bridge = user[1]

								requests.put(
									f"http://backend:{BACKEND_PORT}/api/orchestrator",
									json={
										"bridge": bridge,
										"userid": userid,
										"ingredients": {
											"source": article.get("source").get("name"),
											"author": article.get("author"),
											"publishedAt": article.get("publishedAt"),
											"title": article.get("title"),
											"description": article.get("description"),
											"url": article.get("url"),
											"content": article.get("content")
										}
									}
								)
		
					else:
						print(f"No new articles in the {country} country", file=sys.stderr)


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
    threading.Thread(target=webhook, daemon=True).start()
    app.run(host='0.0.0.0', port=80, debug=True)
