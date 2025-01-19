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
BACKEND_PORT = os.environ.get("BACKEND_PORT")

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

@app.route("/statupdate", methods=["POST"])
def react_update_stats():
    return general_reaction("statupdate", Request.json)

@app.route("/givebadge", methods=["POST"])
def react_givebadge():
    return general_reaction("givebadge", Request.json)

@app.route("/giveitem", methods=["POST"])
def react_giveitem():
    return general_reaction("giveitem", Request.json)

@app.route("/changeprop", methods=["POST"])
def react_changeprop():
    return general_reaction("changeprop", Request.json)

@app.route("/copy", methods=["POST"])
def react_copy():
    return general_reaction("copy", Request.json)

@app.route("/sendmessage", methods=["POST"])
def react_sendmessage():
    return general_reaction("sendmessage", Request.json)

def get_robloxid_from_areaid(areaid) -> str|None:
    print("c", file=stderr)
    with db.cursor() as cur:
        cur.execute(
            "select userid from tokens where owner = %s",
            (int(areaid),)
        )
        rows = cur.fetchone()
        print("d", file=stderr)
        if not rows:
            return None
        return rows[0]
    return None

def general_action(action_name):
    json = Request.get_json()
    bridge = json.get("bridge")
    areaid = json.get("userid")
    spices = json.get("spices")
    print("a", file=stderr)
    if not bridge or not areaid or not spices:
        return jsonify({ "error": "missing bridge, userid or spices" }), 500
    if not "gameid" in spices:
        return jsonify({ "error": "expected a gameid" }), 500
    print("b", file=stderr)
    try:
        robloxid = get_robloxid_from_areaid(areaid)
        with db.cursor() as cur:
            cur.execute("""
                insert into micro_robloxactions
                (bridge, userid, gameid, robloxid, action)
                values (%s, %s, %s, %s, %s)
                on conflict (bridge) do nothing
            """, (
                int(bridge),
                int(areaid),
                str(spices["gameid"]),
                str(robloxid),
                str(action_name)
            ))
            print("e", file=stderr)
            db.commit()
            print("f", file=stderr)
        return jsonify({ "status": "ok"}), 400
    except (Exception, psycopg2.Error) as err:
        return jsonify({ "error": str(err)}), 400

@app.route("/onprompt", methods=["POST"])
def action_onprompt():
    return general_action("onprompt")

@app.route("/onclick", methods=["POST"])
def action_onclick():
    return general_action("onclick")

@app.route("/ontouch", methods=["POST"])
def action_ontouch():
    return general_action("ontouch")

@app.route("/onplayeradded", methods=["POST"])
def action_onplayeradded():
    return general_action("onplayeradded")

@app.route("/onplayerremoved", methods=["POST"])
def action_onplayerremoved():
    return general_action("onplayerremoved")

@app.route("/onchat", methods=["POST"])
def action_onchat():
    return general_action("onchat")

@app.route("/oninput", methods=["POST"])
def action_oninput():
    return general_action("oninput")

def try_getting_informations(robloxid, gameid):
    if robloxid is None:
        print(1, "no roblox id ?", file=stderr)
        return jsonify({"error": "robloxid is required"}), 400
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
        print(2, str(err), file=stderr)
        return jsonify({ "error": "postgres says <(" + str(err) + ")"}), 401
    except Exception as err:
        print(3, str(err), file=stderr)
        return jsonify({ "error": str(err)}), 402

def get_userid_bridge_from_action(gameid: str, action_name: str):
    with db.cursor() as cur:
        cur.execute("""
            select userid, bridge from micro_robloxactions
            where gameid = %s and action = %s
        """, (str(gameid), str(action_name)))
        rows = cur.fetchone()
        if not rows or len(rows) != 2:
            return None, None
        return rows[0], rows[1]
    return None, None

def on_action(data):
    if not "action" in data:
        return jsonify({"error": "expected an action"}), 400
    gameid = data["gameid"]
    action = data["action"]
    userid, bridge = get_userid_bridge_from_action(gameid, action)
    if userid is None:
        return jsonify({"message": "data doesnt exists in database"}), 100
    ingredients = data.get("ingredients")
    requests.put(
        f"http://backend:{BACKEND_PORT}/api/orchestrator",
        json={
            "bridge": bridge,
            "userid": userid,
            "ingredients": ingredients if ingredients is not None else {}
        }
    )
    return jsonify({"status": "ok"}), 200

@app.route("/webhook", methods=["POST"])
def webhook():
    data = Request.get_json()
    if not "gameid" in data:
        return jsonify({"error": "gameid is required"}), 400
    if not "method" in data:
        return jsonify({"error": "expected a method"}), 400
    robloxid = data.get("robloxid")
    gameid = data["gameid"]
    method = data["method"]
    print(gameid, method, file=stderr)
    if method == "retrieve":
        return try_getting_informations(data.get("robloxid"), gameid)
    elif method == "trigger":
        return on_action(data)
    return jsonify({ "error": "invalid method" }), 400

@app.route("/", methods=["GET"])
def routes():
    x = {
        "color": "#000000",
        "image": "/assets/roblox.png",
        "areas": [
            {
                "name": "newpart",
                "description": "Creates a new part",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
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
            },
            {
                "name": "kill",
                "description": "Kill a player",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                    {
                        "name": "userid",
                        "type": "text",
                        "title": "The player userid you want to kill."
                    },
                ]
            },
            {
                "name": "kick",
                "description": "Kick a player",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                    {
                        "name": "message",
                        "type": "text",
                        "title": "The message you want to send. (optionnal)"
                    },
                    {
                        "name": "userid",
                        "type": "text",
                        "title": "The player userid you want to kick."
                    },
                ]
            },
            {
                "name": "insert",
                "description": "Insert an asset from toolbox.",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                    {
                        "name": "parent",
                        "type": "text",
                        "title": "Where to insert it."
                    },
                    {
                        "name": "id",
                        "type": "text",
                        "title": "The ID of the asset."
                    },
                ]
            },
            {
                "name": "statupdate",
                "description": "Update leaderstats.",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                    {
                        "name": "userid",
                        "type": "text",
                        "title": "The player userid you want to update."
                    },
                    {
                        "name": "stat",
                        "type": "text",
                        "title": "The field to update."
                    },
                    {
                        "name": "add",
                        "type": "text",
                        "title": "How much to add (optionnal)."
                    },
                    {
                        "name": "dec",
                        "type": "text",
                        "title": "How much to decrement (optionnal)."
                    },
                ]
            },
            {
                "name": "givebadge",
                "description": "Gives a badge to a user.",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                    {
                        "name": "badgeid",
                        "type": "text",
                        "title": "The badge id."
                    },
                    {
                        "name": "userid",
                        "type": "text",
                        "title": "The player userid whom you want to give it."
                    },
                ]
            },
            {
                "name": "giveitem",
                "description": "Give an item to a player",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                    {
                        "name": "userid",
                        "type": "text",
                        "title": "The player userid whom you want to give it."
                    },
                    {
                        "name": "id",
                        "type": "text",
                        "title": "The asset id of the tool."
                    },
                ]
            },
            {
                "name": "changeprop",
                "description": "Change an Instance property",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                    {
                        "name": "instance",
                        "type": "text",
                        "title": "The Instance path."
                    },
                    {
                        "name": "property",
                        "type": "text",
                        "title": "The property you want to change."
                    },
                    {
                        "name": "value",
                        "type": "text",
                        "title": "The new value."
                    },
                ]
            },
            {
                "name": "copy",
                "description": "Copy an instance.",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                    {
                        "name": "from",
                        "type": "text",
                        "title": "The source instance path."
                    },
                    {
                        "name": "to",
                        "type": "text",
                        "title": "The parent instance path."
                    },
                ]
            },
            {
                "name": "sendmessage",
                "description": "Send a message",
                "type": "reaction",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                    {
                        "name": "message",
                        "type": "text",
                        "title": "The message you want to send."
                    },
                ]
            },
            {
                "name": "onchat",
                "description": "When a player chats.",
                "type": "action",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                ]
            },
            {
                "name": "onprompt",
                "description": "When a prompt is triggered.",
                "type": "action",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                ]
            },
            {
                "name": "onclick",
                "description": "When a stuff is clicked.",
                "type": "action",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                ]
            },
            {
                "name": "ontouch",
                "description": "When a basepart is touched.",
                "type": "action",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                ]
            },
            {
                "name": "onplayeradded",
                "description": "When a player is added.",
                "type": "action",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                ]
            },
            {
                "name": "onplayerremoved",
                "description": "When a player is removed.",
                "type": "action",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                ]
            },
            {
                "name": "oninput",
                "description": "When a key is pressed.",
                "type": "action",
                "spices": [
                    {
                        "name": "gameid",
                        "type": "text",
                        "title": "The game ID."
                    },
                ]
            },
        ]
    }
    return jsonify(x), 200

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=80)
