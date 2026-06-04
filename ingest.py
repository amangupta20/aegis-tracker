import os
import sqlite3

import requests
from dotenv import load_dotenv

load_dotenv()

steam_id = os.getenv("STEAM_ID")

url = f"https://api.opendota.com/api/players/{steam_id}/recentMatches"

response = requests.get(url)

db = sqlite3.connect("stats.db")
cursor = db.cursor()
cursor.execute("PRAGMA foreign_keys = ON;")

for match in response.json():
    if match["player_slot"] < 128:
        side = "Radiant"
    else:
        side = "Dire"
    if side == "Radiant":
        if match["radiant_win"] == True:
            win = 1
        else:
            win = 0
    else:
        if match["radiant_win"] == True:
            win = 0
        else:
            win = 1

    cursor.execute(
        """
        INSERT INTO matches (match_id, win, hero_id, duration, kills, deaths, assists, hero_damage, tower_damage, gpm, xpm, last_hits, side)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    """,
        (
            match["match_id"],
            win,
            match["hero_id"],
            match["duration"],
            match["kills"],
            match["deaths"],
            match["assists"],
            match["hero_damage"],
            match["tower_damage"],
            match["gold_per_min"],
            match["xp_per_min"],
            match["last_hits"],
            side,
        ),
    )
    print("ingested match ", match["match_id"])

db.commit()
