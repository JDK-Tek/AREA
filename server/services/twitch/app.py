import requests
import os
from dotenv import load_dotenv
from random import randint
import psycopg2
from flask import Flask, jsonify
from flask import request as Request
import jwt
from sys import stderr
import datetime as dt
import threading
from time import sleep

load_dotenv("/usr/mount.d/.env")
load_dotenv(".env")

CLIENT_ID = os.environ.get("CLIENTID")
CLIENT_SECRET = os.environ.get("CLIENTSECRET")
BACKEND_KEY = os.environ.get("BACKEND_KEY")
BACKEND_PORT = os.environ.get("BACKEND_PORT")
EXPIRATION = eval(os.environ.get("EXPIRATION"))
REDIRECT = os.environ.get("REDIRECT")

GET_CODE_URL = "https://id.twitch.tv/oauth2/authorize"
GET_TOKEN_URL = "https://id.twitch.tv/oauth2/token"
GET_ID_URL = "https://api.twitch.tv/helix/users"
IS_STREAMING_URL = "https://api.twitch.tv/helix/streams"

# the database

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

# general functions

def generate_state(size: int):
    return "".join([chr(randint(ord('a'), ord('z'))) for _ in range(size)])

def myurlencode(x: dict):
    return "&".join("{}={}".format(*i) for i in x.items())

# everything user related

def create_empty_user() -> tuple[bool, any]:
    with db.cursor() as cur:
        cur.execute("insert into users default values returning id")
        rows = cur.fetchone()
        db.commit()
        if not rows:
            return False, "couldnt insert"
        return True, rows[0]
    return False, "couldnt get the cursor"

def retrieve_id_from_atok(atok: str):
	try:
		return jwt.decode(atok, BACKEND_KEY, algorithms=["HS256"])
	except jwt.ExpiredSignatureError:
		return None
	except jwt.InvalidTokenError:
		return None

def push_token_in_database(tok: str, refreshtok: str, userid: str, areaid: int):
    with db.cursor() as cur:
        cur.execute("""
            insert into tokens (service, token, refresh, userid, owner)
            values (%s, %s, %s, %s, %s)
            returning id
        """, ("twitch", tok, refreshtok, userid, areaid))
        rows = cur.fetchone()
        db.commit()
        if not rows:
            return False, "couldnt insert"
        return True, rows[0]
    return False, "couldnt connect"

def get_twitch_user_id(tok: str, name: str|None) -> tuple[bool, any]:
    headers = {
        "Authorization": f"Bearer {tok}",
        "Client-ID": CLIENT_ID
    }
    url = GET_ID_URL if (name is None) else GET_ID_URL + "?" + myurlencode({
        "login": name
    })
    rep = requests.get(url, headers=headers)
    if rep.status_code != 200:
        return False, rep.text
    data = rep.json()
    if "data" not in data or len(data["data"]) == 0:
        return False, "no data in parametters"
    userid = data["data"][0]["id"]
    return True, userid

def get_access_tokens(code: str) -> tuple[bool, str, str]:
    data = {
        "client_id": CLIENT_ID,
        "client_secret": CLIENT_SECRET,
        "code": code,
        "grant_type": "authorization_code",
        "redirect_uri": REDIRECT,
    }
    rep = requests.post(GET_TOKEN_URL, data=data)
    token_data = rep.json()
    if "error" in token_data:
        return False, token_data["error"]
    if not "access_token" in token_data or not "refresh_token" in token_data:
        return False, "cant get access token", None
    print(token_data["access_token"], file=stderr)
    return True, token_data["access_token"], token_data["refresh_token"]

def create_area_token(id: int):
    return jwt.encode({
		"id": id,
		"exp": dt.datetime.now() + dt.timedelta(seconds=EXPIRATION)
	}, BACKEND_KEY, algorithm="HS256")

# everything that touches oauth

def generate_oauth_link():
    state = generate_state(40)
    params = {
        "response_type": "code",
        "client_id": CLIENT_ID,
        "redirect_uri": REDIRECT,
        "scope": "channel%3Amanage%3Apolls+channel%3Aread%3Apolls",
        "state": state
    }
    return GET_CODE_URL + "?" + myurlencode(params)

def exchange_code_for_token(code: str, atok: str|None) -> str:
    success, tok, refreshtok = get_access_tokens(code)
    assert success, tok
    success, userid = get_twitch_user_id(tok, None)
    assert success, userid
    print(userid, file=stderr)
    areaid = None
    if (atok is not None) and (len(atok) != 0):
        areaid = retrieve_id_from_atok(atok)
        assert areaid, "invalid token"
    else:
        success, areaid = create_empty_user()
        assert areaid, areaid
    if areaid is None:
        return "couldnt get areaid"
    print(areaid, file=stderr)
    success, tokid = push_token_in_database(tok, refreshtok, userid, areaid)
    assert success, str(tokid)
    return create_area_token(areaid)

@app.route("/oauth", methods=["GET", "POST"])
def oauth_route():
    if Request.method == "GET":
        return generate_oauth_link(), 200
    data = Request.get_json()
    if not "code" in data:
        return jsonify({ "error": "missing code" }), 400
    try:
        atok = exchange_code_for_token(
            data["code"],
            Request.headers.get("Authorization")
        )
        return jsonify({ "token": atok })
    except psycopg2.Error as err:
        return jsonify({ "error": f"database: {str(err)}" }), 500
    except AssertionError as err:
        return jsonify({ "error": f"assert: {str(err)}" }), 400
    except BaseException as err:
        return jsonify({ "error": f"unknown: {str(err)}" }), 400

# everything that touches to actions/reactions

def is_streamer_streaming(token: str, name: str) -> tuple[bool, any]:
    success, res = get_twitch_user_id(token, name)
    if not success:
        return False, res
    url = IS_STREAMING_URL + "?" + myurlencode({
        "user_id": res
    })
    headers = {
        "Authorization": f"Bearer {token}",
        "Client-ID": CLIENT_ID
    }
    rep = requests.get(url, headers=headers)
    if rep.status_code != 200:
        return False, rep.text
    stream_data = rep.json()
    if stream_data["data"]:
        return True, True
    else:
        return True, False

@app.route("/onstreamstart", methods=["POST"])
def action_isstreaming():
    data = Request.get_json()
    if not "spices" in data:
        return jsonify({ "error": "missing spices" }), 400
    if not "streamer" in data["spices"]:
        return jsonify({ "error": "expected streamer" }), 400
    spices = data["spices"]
    with db.cursor() as cur:
        cur.execute("select userid from tokens where owner = %s",
                    (data["userid"],))
        rows = cur.fetchone()
        assert rows, "can't get rows"
        userid = rows[0]
    with db.cursor() as cur:
        cur.execute("""
            insert into micro_twitch_onstream
            (streamer, userid, areaid, bridge)
            values (%s, %s, %s, %s)
        """, (spices["streamer"], userid, data["userid"], data["bridge"]))
    db.commit()
    return jsonify({ "status": "ok" })

# @app.route("/onstreamend", methods=["POST"])
# def action_isnotstreaming():
#     pass

# the master thread

def call_bridge(bridge: int, areaid: int):
    requests.put(
        f"http://backend:{BACKEND_PORT}/api/orchestrator",
        json = {
            "bridge": bridge,
			"userid": areaid,
            "ingredients": {}
        }
    )

def treat_user(userid: str, streamer: str, connected: bool, bridge: int, areaid: int):
    with db.cursor() as cur:
        cur.execute("select token from tokens where userid = %s", (userid,))
        rows = cur.fetchone()
    if not rows:
        return
    tok = rows[0]
    print(streamer, tok, connected, file=stderr)
    success, now_connected = is_streamer_streaming(tok, streamer)
    if not success:
        return
    print(streamer, tok, connected, now_connected, file=stderr)
    if now_connected and not connected:
        call_bridge(bridge, areaid)
        with db.cursor() as cur:
            cur.execute("update micro_twitch_onstream set connected = %s", (True,))
            db.commit()

def update_onstreeam():
    with db.cursor() as cur:
        cur.execute("select userid, streamer, connected, bridge, areaid from micro_twitch_onstream")
        rows = cur.fetchall()
    assert (rows is not None), "couldnt fetch rows"
    for u in rows:
        print(u, file=stderr)
        treat_user(u[0], u[1], u[2], u[3], u[4])
    
def master_thread():
    while True:
        try:
            update_onstreeam()
        except AssertionError as err:
            print("an error occured (assert): ", str(err), file=stderr)
        except psycopg2.Error as err:
            print("an error occured (database): ", str(err), file=stderr)
        except BaseException as err:
            print("an error occured (unknwon): ", str(err), file=stderr)
        sleep(5)

if __name__ == "__main__":
    threading.Thread(target=master_thread).start()
    app.run(host='0.0.0.0', port=80)

# print(generate_oauth_link())
# exchange_code_for_token("r33bqd54s1q0886yslxe3r155zc1eb")
# print(is_streamer_streaming("g7snbv64e5hzof8uxhofhos8kvhlj7", "potatozytb"))
