package model

import (
	"log"

	"gopkg.in/guregu/null.v3"
)

// MetarStationIDMv represents a row in mv metar_station_id_mv
type MetarStationIDMv struct {
	StationID null.String `db:"STATION_ID" json:"stationId"`
	Name      null.String `json:"name"`
}

// SelectStationIDs function
func SelectStationIDs() []MetarStationIDMv {
	metarSlice := []MetarStationIDMv{}
	sqlStatement := `
		select
			station_id,
			name
		from metar_station_id_mv
	`
	err := Db.Select(&metarSlice, sqlStatement)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Print(len(metarSlice))

	return metarSlice

}
