package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

// Metar represents a row in table flight-log
type Metar struct {
	StationID       null.String `db:"STATION_ID" json:"stationId"`
	ObservationTime time.Time   `json:"observationTime"`
	RawText         null.String `json:"rawText"`
	// Version             int64           //`db:"version"`
	// CoPilot             sql.NullString  //`db:"CO_PILOT"`
	// Created             time.Time       //`db:"CREATED"`
	// DayDual             sql.NullFloat64 //`db:"DAY_DUAL"`
	// DaySolo             sql.NullFloat64 //`db:"DAY_SOLO"`
	// FlightDate          time.Time       //`db:"FLIGHT_DATE"`
	// InstrumentFlightSim sql.NullFloat64 //`db:"INSTRUMENT_FLIGHT_SIM"`
	// InstrumentImc       sql.NullFloat64 //`db:"INSTRUMENT_IMC"`
	// InstrumentNoIfrAppr sql.NullInt32   //`db:"INSTRUMENT_NO_IFR_APPR"`
	// InstrumentSimulated sql.NullFloat64 //`db:"INSTRUMENT_SIMULATED"`
	// MakeModel           sql.NullString  //`db:"MAKE_MODEL"`
	// Modified            time.Time       //`db:"MODIFIED"`
	// NightDual           sql.NullFloat64 //`db:"NIGHT_DUAL"`
	// NightSolo           sql.NullFloat64 //`db:"NIGHT_SOLO"`
	// Pic                 sql.NullString  //`db:"PIC"`
	// Registration        sql.NullString  //`db:"REGISTRATION"`
	// Remarks             sql.NullString  //`db:"REMARKS"`
	// RouteFrom           sql.NullString  //`db:"ROUTE_FROM"`
	// RouteTo             sql.NullString  //`db:"ROUTE_TO"`
	// TosLdgsDay          sql.NullInt32   //`db:"TOS_LDGS_DAY"`
	// TosLdgsNight        sql.NullInt32   //`db:"TOS_LDGS_NIGHT"`
	// XCountryDay         sql.NullFloat64 //`db:"X_COUNTRY_DAY"`
	// XCountryNight       sql.NullFloat64 //`db:"X_COUNTRY_NIGHT"`
}

// GetUser function
func GetUser() string {
	var user string
	var err error
	if Db == nil {
		log.Println("db is null")
	}
	err = Db.Get(&user, "select user from dual")
	if err != nil {
		log.Println(err)
		return ""
	}
	log.Println(user)
	return user
}

// SelectMetars function
func SelectMetars(stationIDs []string, fromObservationTime time.Time, toObservationTime time.Time) []Metar {
	log.Printf("stationIDs: %v %T, fromObservationTime: %v %T, toObservationTime: %v %T", stationIDs, stationIDs, fromObservationTime, fromObservationTime, toObservationTime, toObservationTime)
	log.Printf("length of stationIDs: %v", len(stationIDs))
	metarSlice := []Metar{}
	var err error
	//err = Db.Select(&metarSlice, "select station_id, observation_time, raw_text from metar where station_id in (:1) and observation_time >= :2 and observation_time <= :3 order by observation_time", stationIDs, fromObservationTime, toObservationTime)
	query, args, err := sqlx.In("select station_id, observation_time, raw_text from metar where station_id in (?) and observation_time >= ? and observation_time <= ? order by station_id, observation_time", stationIDs, fromObservationTime, toObservationTime)
	if err != nil {
		log.Println(err)
		return nil
	}
	query = Db.Rebind(query)
	err = Db.Select(&metarSlice, query, args...)
	log.Print(len(metarSlice))

	//log.Print(flightLogSlice[len(flightLogSlice)-1].FlightDate)

	var jsonData []byte
	jsonData, err = json.Marshal(metarSlice[0])
	if err != nil {
		log.Println(err)
	}
	log.Println("before ===============================")
	log.Println(string(jsonData))
	log.Println("after ===============================")

	return metarSlice

}
