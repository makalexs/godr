package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/makalexs/godr/config"
	"log"
)

type GetObjectNearPoint struct {
	OwnerInfo     struct {
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		MiddleName string `json:"middleName"`
		PhoneMain  string `json:"phoneMain"`
		PhoneAdd   string `json:"phoneAdd"`
		Email      string `json:"email"`
	}
	SpecInfo      struct {
		Type          string `json:"type"`
		Area          int `json:"area"`
		Rooms      	  int `json:"rooms"`
		SchemeFileUrl string `json:"schemeFileUrl"`
		Floor         int `json:"floor"`
	}
	Price         float64 `json:"price"`
}

func CommonInsert(sql string) uint64 {
	var id uint64
	databaseUrl := "postgres://"+config.GetConfiguration().DatabasePostgres.User+":"+config.GetConfiguration().DatabasePostgres.Pass
	databaseUrl = databaseUrl +"@"+config.GetConfiguration().DatabasePostgres.Url+":"+config.GetConfiguration().DatabasePostgres.Port+"/"+config.GetConfiguration().DatabasePostgres.Database
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalln(err.Error())
		return 0
	}
	defer conn.Close(context.Background())

	row := conn.QueryRow(context.Background(), sql)
	err = row.Scan(&id)
	if err != nil {
		log.Fatalln(err.Error())
		return 0
	}

	return id
}

func CommonSelectNearPoint(sql string) GetObjectNearPoint {
	databaseUrl := "postgres://"+config.GetConfiguration().DatabasePostgres.User+":"+config.GetConfiguration().DatabasePostgres.Pass
	databaseUrl = databaseUrl +"@"+config.GetConfiguration().DatabasePostgres.Url+":"+config.GetConfiguration().DatabasePostgres.Port+"/"+config.GetConfiguration().DatabasePostgres.Database
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalln(err.Error())
		return GetObjectNearPoint{}
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		log.Fatalln(err.Error())
		return GetObjectNearPoint{}
	}

	var object GetObjectNearPoint
	for rows.Next() {
		var firstName string
		var lastName string
		var middleName string
		var phoneMain string
		var phoneAdd string
		var email string
		var objectType string
		var area int
		var rooms int
		var schemeFileUrl string
		var floor int
		var price float64
		if err := rows.Scan(&firstName, &lastName, &middleName, &phoneMain, &phoneAdd, &email, &objectType, &area, &rooms, &schemeFileUrl, &floor, &price); err == nil {
			object = GetObjectNearPoint{}
			object.OwnerInfo.FirstName = firstName
			object.OwnerInfo.LastName = lastName
			object.OwnerInfo.MiddleName = middleName
			object.OwnerInfo.PhoneMain = phoneMain
			object.OwnerInfo.PhoneAdd = phoneAdd
			object.OwnerInfo.Email = email
			object.SpecInfo.Type = objectType
			object.SpecInfo.Area = area
			object.SpecInfo.Rooms = rooms
			object.SpecInfo.SchemeFileUrl = schemeFileUrl
			object.SpecInfo.Floor = floor
			object.Price = price
			return object
		} else {
			log.Fatalln(err.Error())
			return GetObjectNearPoint{}
		}
	}
	return GetObjectNearPoint{}
}