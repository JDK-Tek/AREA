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

# "technology"
# "entertainment",
# "buisness",
# "general",
# "science",
# "health",
# "sports",

# When a new article is published in the 'technology' category
ACTION_NEW_ART_TECH = "new-art-tech"
oreo.create_area(
	ACTION_NEW_ART_TECH,
	NewOreo.TYPE_ACTIONS,
	"Any new article in the 'technology' category",
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

ACTION_NEW_ART_BUISNESS = "new-art-buisness"
oreo.create_area(
	ACTION_NEW_ART_BUISNESS,
	NewOreo.TYPE_ACTIONS,
	"Any new article in the 'buisness' category",
	[]
)
@app.route(f'/{ACTION_NEW_ART_BUISNESS}', methods=["POST"])
def new_article_buisness():
	app.logger.info(f"{ACTION_NEW_ART_BUISNESS} endpoint hit")

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
				  ACTION_NEW_ART_BUISNESS,
				  json.dumps(spices)
			  )
		)

		db.commit()

	return jsonify({"status": "ok"}), 200


CATEGORIES = {
	ACTION_NEW_ART_TECH: "technology",
	ACTION_NEW_ART_BUISNESS: "buisness"
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


def webhook():
	print("Starting newsapi webhook", file=sys.stderr)

	while True:

		time.sleep(20)

		# get all subscribed users to the action
		with db.cursor() as cur:
			cur.execute("SELECT * FROM micro_newsapi")
			users_sub = cur.fetchall()

			if not users_sub or len(users_sub) == 0:
				print("No users subscribed to the newsapi actions", file=sys.stderr)
				continue
			
			print(f"Found {len(users_sub)} users subscribed to the newsapi actions", file=sys.stderr)

	
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
