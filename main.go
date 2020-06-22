package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/godror/godror"
	"github.com/magiconair/properties"
	"github.com/thalabi/AvWxServer/model"
)

var err error
var timeLocation, _ = time.LoadLocation("Local")

func main() {
	log.Println("Starting server ...")

	prop := properties.MustLoadFile("application.properties", properties.UTF8)
	model.InitDB(prop)

	if prop.GetBool("production-environment", false) {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	// cors - allows all origins
	router.Use(cors.Default())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/pingDb", func(c *gin.Context) {
		user := model.GetUser()
		c.JSON(200, gin.H{
			"message": user,
		})
	})

	router.GET("/getStationIds", func(c *gin.Context) {
		stationIDs := model.SelectStationIDs()
		c.JSON(http.StatusOK, stationIDs)
	})

	router.GET("/getMetarListInObervationTimeRange", func(c *gin.Context) {
		//timeLocation, _ := time.LoadLocation("Local")

		const (
			stationIDsKey               = "stationId"
			fromObservationTimeKey      = "fromObservationTime"
			toObservationTimeKey        = "toObservationTime"
			lastNumberOfObservationsKey = "lastNumberOfObservations"
		)
		var paramsMap map[string][]string
		paramsMap = c.Request.URL.Query()
		log.Println("paramsMap: ", paramsMap)

		if stationIDs := paramsMap[stationIDsKey]; stationIDs == nil {
			c.JSON(http.StatusBadRequest, "Missing stationId")
			return
		}
		stationIDs := paramsMap[stationIDsKey]
		upperStationIDs(&stationIDs)
		log.Printf("stationIDs: %v %T", stationIDs, stationIDs)
		log.Printf("length of stationIDs: %v", len(stationIDs))

		if paramsMap[fromObservationTimeKey] == nil && paramsMap[lastNumberOfObservationsKey] == nil {
			c.JSON(http.StatusBadRequest, "Request must supply fromObservationTime or lastNumberOfObservations")
			return
		}
		var fromObservationTime time.Time
		layout := "2006-01-02"
		fromObservationTime, err = time.ParseInLocation(layout, paramsMap[fromObservationTimeKey][0], timeLocation)
		if err != nil {
			c.JSON(http.StatusBadRequest, paramsMap[fromObservationTimeKey][0]+" is not a valid date. Format is yyyy-mm-dd")
			return
		}
		var toObservationTime time.Time
		if paramsMap[toObservationTimeKey] != nil {
			t, err := time.ParseInLocation(layout, paramsMap[toObservationTimeKey][0], timeLocation)
			if err != nil {
				c.JSON(http.StatusBadRequest, paramsMap[toObservationTimeKey][0]+" is not a valid date. Format is yyyy-mm-dd")
				return
			}
			toObservationTime = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, t.Nanosecond(), t.Location())
		} else {
			toObservationTime = time.Now().In(timeLocation)
		}
		metarArray := model.SelectMetarListInObervationTimeRange(stationIDs, fromObservationTime, toObservationTime)
		c.JSON(http.StatusOK, metarArray)
	})

	//
	//
	//
	router.GET("/getMetarListForLatestNObservations", func(c *gin.Context) {
		//timeLocation, _ := time.LoadLocation("Local")

		const (
			stationIDsKey           = "stationId"
			latestNumberOfMetarsKey = "latestNumberOfMetars"
		)
		var paramsMap map[string][]string
		paramsMap = c.Request.URL.Query()
		log.Println("paramsMap: ", paramsMap)

		stationIDs, ok := paramsMap[stationIDsKey]
		if ! /* not */ ok {
			c.JSON(http.StatusBadRequest, "Missing stationId")
			return
		}
		upperStationIDs(&stationIDs)
		log.Printf("stationIDs: %v %T", stationIDs, stationIDs)
		log.Printf("length of stationIDs: %v", len(stationIDs))

		latestNumberOfMetars, ok := paramsMap[latestNumberOfMetarsKey]
		if ! /* not */ ok {
			c.JSON(http.StatusBadRequest, "Missing latestNumberOfMetars")
			return
		}
		log.Printf("latestNumberOfMetars: %v %T", latestNumberOfMetars, latestNumberOfMetars)
		log.Printf("latestNumberOfMetars of stationIDs: %v", len(latestNumberOfMetars))
		metarArray := model.SelectMetarListForLatestNObservations(stationIDs, latestNumberOfMetars[0])
		c.JSON(http.StatusOK, metarArray)
	})

	router.GET("/getBuildVersionAndTimestamp", func(c *gin.Context) {
		buildVersion := prop.GetString("build-version", "N/A")
		buildTimestamp := prop.GetString("build-timestamp", "N/A")

		c.JSON(http.StatusOK, buildVersion+"_"+buildTimestamp)
	})

	router.Run(":" + prop.GetString("http-port", "8080"))
}

func upperStationIDs(stationIDs *[]string) {
	x := *stationIDs
	for i, stationID := range x {
		x[i] = strings.ToUpper(stationID)
	}
}
