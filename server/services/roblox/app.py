from flask import Flask, jsonify
from flask import request as Request
from dotenv import load_dotenv
import random
import requests
import string
import os
import psycopg2
import time
import jwt
import datetime as dt
from sys import stderr

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

class Command:
    def __init__(self, name: str, data: dict[str, str]):
        self.name = name
        self.extra = " ".join([k + " " + v for k, v in data.items()])
        self.str = self.name + " " + self.extra

    def __str__(self):
        return self.str

def general_reaction(name, data):
    areaid = data.get("userid")
    spices: dict = data.get("spices")

    if not "gameid" in spices:
         return jsonify({ "error": "expected (at least) gameid" }), 400

    gameid = spices.get("gameid")
    spices.pop("gameid")
    command = Command(name, spices)

    try:
        with db.cursor() as cur:
            print(command, gameid, file=stderr)
            cur.execute("update micro_roblox "\
                "set command = %s "\
                "from tokens "\
                "where micro_roblox.robloxid = tokens.userid "\
                "and tokens.owner = %s "\
                "and tokens.service = %s "\
                "and micro_roblox.gameid = %s",
                (str(command), str(areaid), "roblox", str(gameid),)
            )
            nrows = cur.rowcount
            db.commit()
            if nrows <= 0:
                return jsonify({ "error": "awaiting for your game to connect at least once." }), 425
                 
        return jsonify({ "status": "ok" }), 200
    except (Exception, psycopg2.Error) as err:
        return jsonify({ "error":  str(err)}), 400

@app.route("/newpart", methods=["POST"])
def react_newpart():
    return general_reaction("newpart", Request.json)

@app.route("/kill", methods=["POST"])
def react_kill():
    return general_reaction("kill", Request.json)

@app.route("/kick", methods=["POST"])
def react_kick():
    return general_reaction("kick", Request.json)

@app.route("/insert", methods=["POST"])
def react_insert():
    return general_reaction("insert", Request.json)

def try_getting_informations(robloxid, gameid):
    try:
        with db.cursor() as cur:
            cur.execute("""
                insert into micro_roblox (robloxid, gameid)
                values (%s, %s)
                on conflict (gameid)
                do nothing
            """, (str(robloxid), str(gameid)))
            db.commit()
            cur.execute(
                "select command from micro_roblox "\
                "where gameid = %s "\
                "and command is not null",
                (str(gameid),)
            )
            commands = cur.fetchall()
            command_list = [c[0] for c in commands]
            cur.execute(
                "update micro_roblox set command = null where gameid = %s",
                (str(gameid),)
            )
            db.commit()
            return jsonify({ "list": command_list}), 200
    except (psycopg2.Error) as err:
        return jsonify({ "error": "postgres says <(" + str(err) + ")"}), 400
    except Exception as err:
        return jsonify({ "error": str(err)}), 400

@app.route("/webhook", methods=["POST"])
def webhook():
    data = Request.get_json()
    if not "gameid" in data:
        return jsonify({"error": "gameid is required"}), 400
    if not "robloxid" in data:
        return jsonify({"error": "robloxid is required"}), 400
    if not "method" in data:
        return jsonify({"error": "expected a method"}), 400
    robloxid = data.get("robloxid")
    gameid = data.get("gameid")
    method = data.get("method")
    print(robloxid, gameid, method, file=stderr)
    if method == "retrieve":
        return try_getting_informations(robloxid, gameid)
        

@app.route("/", methods=["GET"])
def routes():
    x = {
        "color": "#000000",
        "image": "/assets/roblox.png",
        "areas": [
            {
                "name": "newpart",
                "description": "Creates a new part",
                "spices": [
                    {
                        "name": "color",
                        "type": "text",
                        "title": "The color of the part."
                    },
                    {
                        "name": "position",
                        "type": "text",
                        "title": "3 numbers for the part coordinates."
                    },
                    {
                        "name": "size",
                        "type": "text",
                        "title": "3 numbers for the part size."
                    },
                    {
                        "name": "anchored",
                        "type": "dropdown",
                        "title": "Is the part anchored.",
                        "extra": ["true", "false"]
                    }
                ]
            }
        ]
    }
    return jsonify(x), 200

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=80)
