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

@app.route("/print", methods=["POST"])
def roblox_print():
    data = Request.json
    areaid = data.get("userid")
    spices = data.get("spices")
    command = spices.get("command")
    gameid = spices.get("gameid")
    # "select userid from tokens where owner = %s limit 1", areaid
    # "update micro_roblox set command = %s where robloxid = %d"
    # => "update micro_roblox set command = %s from tokens where micro_roblox.robloxid = tokens.userid and tokens.owner = %d and micro_roblox.gameid = %d", command, areaid, gameid
    try:
        with db.cursor() as cur:
            print(command, gameid, file=stderr)
            cur.execute("update micro_roblox "\
                "set command = %s "\
                "from tokens "\
                "where micro_roblox.robloxid = tokens.userid "\
                "and tokens.owner = %s "\
                "and micro_roblox.gameid = %s ",
                (command, str(areaid), gameid,)
            )
            print("foo", file=stderr)
            db.commit()
            print("bar", file=stderr)
            return jsonify({ "status", "ok" }), 200
    except (Exception, psycopg2.Error) as err:
        return jsonify({ "error":  str(err)}), 400

@app.route("/webhook", methods=["POST"])
def webhook():
    data = Request.json
    robloxid = data.get("robloxid")
    gameid = data.get("gameid")
    # => "insert into micro_roblox (robloxid, gameid) values (%d, %d) on conflict (gameid) do nothing", robloxid, gameid
    try:
        with db.cursor() as cur:
            print(robloxid, gameid, file=stderr)
            cur.execute(
                 "insert into micro_roblox (robloxid, gameid) "\
                 "values (%s, %s) "\
                 "on conflict (gameid) "\
                 "do nothing",
                (robloxid, gameid,)
            )
            db.commit()
        while True:
            # "select command from micro_roblox where robloxid = %d and command is not null", robloxid
            # => "update micro_roblox set command = null where robloxid = %d and command is not null returning command", robloxid
            with db.cursor() as cur:
                cur.execute(
                     "update micro_roblox "\
                     "set command = null "\
                     "where gameid = %s "\
                     "and command is not null "\
                     "returning command",
                    (gameid,)
                )
                rows = cur.fetchone()
                db.commit()
                if rows:
                    return jsonify({ "message": str(rows[0])}), 200
            time.sleep(1)
    except (psycopg2.Error) as err:
        return jsonify({ "error": "postgres says <(" + str(err) + ")"}), 400
    except Exception as err:
        return jsonify({ "error": str(err)}), 400

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=80)