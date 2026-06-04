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

print(response.json())
