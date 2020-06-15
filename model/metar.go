package model

import (
	"strconv"
	"strings"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

// Metar represents a row in table flight-log
type Metar struct {
	StationID           null.String `db:"STATION_ID" json:"stationId"`
	ObservationTime     time.Time   `json:"observationTime"`
	RawText             null.String `json:"rawText"`
	WindDirDegrees      null.Float  `json:"windDirDegrees"`
	WindSpeedKt         null.Float  `json:"windSpeedKt"`
	WindGustKt          null.Float  `json:"windGustKt"`
	VisibilityStatuteMi null.Float  `json:"visibilityStatuteMi"`
	WxString            null.String `json:"wxString"`
	Auto                null.String `json:"auto"`
	SkyCover1           null.String `db:"SKY_COVER_1" json:"skyCover1"`
	CloudBaseFtAgl1     null.Float  `db:"CLOUD_BASE_FT_AGL_1" json:"cloudBaseFtAgl1"`
	SkyCover2           null.String `db:"SKY_COVER_2" json:"skyCover2"`
	CloudBaseFtAgl2     null.Float  `db:"CLOUD_BASE_FT_AGL_2" json:"cloudBaseFtAgl2"`
	SkyCover3           null.String `db:"SKY_COVER_3" json:"skyCover3"`
	CloudBaseFtAgl3     null.Float  `db:"CLOUD_BASE_FT_AGL_3" json:"cloudBaseFtAgl3"`
	SkyCover4           null.String `db:"SKY_COVER_4" json:"skyCover4"`
	CloudBaseFtAgl4     null.Float  `db:"CLOUD_BASE_FT_AGL_4" json:"cloudBaseFtAgl4"`
	VertVisFt           null.Float  `json:"vertVisFt"`
	TempC               null.Float  `json:"tempC"`
	DewpointC           null.Float  `json:"dewpointC"`
	AltimInHg           null.Float  `json:"altimInHg"`
	// Version             int64           //`db:"version"`
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

// SelectMetarListInObervationTimeRange function
func SelectMetarListInObervationTimeRange(stationIDs []string, fromObservationTime time.Time, toObservationTime time.Time) []Metar {

	log.Printf("stationIDs: %v %T, fromObservationTime: %v %T, toObservationTime: %v %T", stationIDs, stationIDs, fromObservationTime, fromObservationTime, toObservationTime, toObservationTime)
	log.Printf("length of stationIDs: %v", len(stationIDs))

	metarSlice := []Metar{}

	const sqlStatement = `
		select
			station_id,
			observation_time,
			auto,
			raw_text,
			wind_dir_degrees,
			wind_speed_kt,
			wind_gust_kt,
			visibility_statute_mi,
			wx_string, sky_cover_1,
			cloud_base_ft_agl_1,
			sky_cover_2,
			cloud_base_ft_agl_2,
			sky_cover_3,
			cloud_base_ft_agl_3,
			sky_cover_4,
			cloud_base_ft_agl_4,
			vert_vis_ft,
			temp_c,
			dewpoint_c,
			altim_in_hg
		from metar
		where
			station_id in (?)
			and observation_time >= ? and observation_time <= ?
			order by station_id, observation_time`

	query, args, err := sqlx.In(sqlStatement, stationIDs, fromObservationTime, toObservationTime)
	if err != nil {
		log.Println("Error in excuting sqlx.In()")
		log.Println(err)
		return nil
	}
	log.Println("After excuting sqlx.In()")
	//query = Db.Rebind(query)
	query = oracleRebind(query)
	err = Db.Select(&metarSlice, query, args...)
	if err != nil {
		log.Println("Error in excuting Db.Select()")
		log.Println(err)
		return nil
	}
	log.Println("After excuting Db.Select()")
	return metarSlice
}

// SelectMetarListForLatestNObservations function
func SelectMetarListForLatestNObservations(stationIDs []string, latestNumberOfMetars string) []Metar {

	log.Printf("stationIDs: %v %T, latestNumberOfMetars: %v %T", stationIDs, stationIDs, latestNumberOfMetars, latestNumberOfMetars)
	log.Printf("length of stationIDs: %v", len(stationIDs))

	metarSlice := []Metar{}

	const sqlStatement = `
		select
			station_id,
			observation_time,
			auto,
			raw_text,
			wind_dir_degrees,
			wind_speed_kt,
			wind_gust_kt,
			visibility_statute_mi,
			wx_string,
			sky_cover_1,
			cloud_base_ft_agl_1,
			sky_cover_2,
			cloud_base_ft_agl_2,
			sky_cover_3,
			cloud_base_ft_agl_3,
			sky_cover_4,
			cloud_base_ft_agl_4,
			vert_vis_ft,
			temp_c,
			dewpoint_c,
			altim_in_hg
		from (select
				station_id,
					observation_time,
					auto,
					raw_text,
					wind_dir_degrees,
					wind_speed_kt,
					wind_gust_kt,
					visibility_statute_mi,
					wx_string,
					sky_cover_1,
					cloud_base_ft_agl_1,
					sky_cover_2,
					cloud_base_ft_agl_2,
					sky_cover_3,
					cloud_base_ft_agl_3,
					sky_cover_4,
					cloud_base_ft_agl_4,
					vert_vis_ft,
					temp_c,
					dewpoint_c,
					altim_in_hg,
					rank() over (partition by station_id order by observation_time desc) as observation_time_rank
				from metar
				where station_id in (?)
					and observation_time >= sysdate-3
			)
	where observation_time_rank <= ? order by station_id, observation_time`

	query, args, err := sqlx.In(sqlStatement, stationIDs, latestNumberOfMetars)
	if err != nil {
		log.Println("Error in excuting sqlx.In()")
		log.Println(err)
		return nil
	}
	log.Println("After excuting sqlx.In()")

	//query = Db.Rebind(query)
	query = oracleRebind(query)

	err = Db.Select(&metarSlice, query, args...)
	if err != nil {
		log.Println("Error in excuting Db.Select()")
		log.Println(err)
		return nil
	}
	log.Println("After excuting Db.Select()")
	log.Print(len(metarSlice))
	return metarSlice
}

// Replaces ? with :n bind placeholder
func oracleRebind(sqlStatement string) string {
	var i = 0
	for strings.Index(sqlStatement, "?") > -1 {
		i++
		sqlStatement = strings.Replace(sqlStatement, "?", ":"+strconv.Itoa(i), 1)
	}
	return sqlStatement
}