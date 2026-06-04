package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)
type herostats struct {
	name    string
	games   int
	wins    int
	winrate float64
	kills   float64
	deaths  float64
	assists float64
	kda     float64
	gpm     float64
	xpm     float64
}

func main() {
	db, err := sql.Open("sqlite", "./stats.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to SQLite!")
	query := `
	select h.name as name,
	count(m.hero_id) as games,
	sum(m.win) as wins,
	avg(m.kills) as kills,
	avg(m.deaths) as deaths,
	avg(m.assists) as assists,
	COALESCE(avg(m.gpm),0) as gpm,
	COALESCE(avg(m.xpm),0) as xpm
	from matches m
	inner join heroes h on h.id = m.hero_id
	group by h.name
	order by wins desc,games desc;
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}	
	defer rows.Close()
	fmt.Println("-----------------------------------------")
	fmt.Println("Hero              | Games | Wins | Winrate% | Kills | Deaths | Assists | KDA     |  GPM |  XPM")
	fmt.Println("-----------------------------------------")
	for rows.Next(){
		var stats herostats
		err := rows.Scan(&stats.name,&stats.games,&stats.wins,&stats.kills,&stats.deaths,&stats.assists,&stats.gpm,&stats.xpm)
		if err != nil{
			log.Fatal(err)
		}
		stats.winrate = float64(stats.wins)/float64(stats.games)*100
		if stats.deaths!=0{
			stats.kda = (stats.kills + stats.assists) / stats.deaths
		}else{
			stats.kda = (stats.kills + stats.assists)		}

		fmt.Printf("%-15s | %3d | %3d | %6.2f%% | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f\n", stats.name,
			stats.games,stats.wins,stats.winrate,stats.kills,stats.deaths,stats.assists,stats.kda,stats.gpm,stats.xpm)
	}
	
		
}
