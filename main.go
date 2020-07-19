package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/godror/godror"
	"github.com/magiconair/properties"
	"github.com/thalabi/AvWxServer/model"
	"gopkg.in/natefinch/lumberjack.v2"
)

var err error
var timeLocation, _ = time.LoadLocation("Local")
var prop *properties.Properties

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Credentials Create a struct to read the username and password from the request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type authority struct {
	Authority string `json:"authority"`
}

//User Structure returned to client
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	// Authorities []struct {
	// 	Authority string `json:"authority"`
	// } `json:"authorities"`
	Authorities []authority `json:"authorities"`
	Token       string      `json:"token"`
}

func main() {
	prop = properties.MustLoadFile("application.properties", properties.UTF8)

	setUpLogging()

	log.Println("Starting server ...")

	model.InitDB(prop)

	if prop.GetBool("production-environment", false) {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/ping", ping)
	router.GET("/pingDb", pingDb)
	router.GET("/getStationIds", getStationIds)
	router.GET("/getMetarListInObervationTimeRange", getMetarListInObervationTimeRange)
	router.GET("/getMetarListForLatestNObservations", getMetarListForLatestNObservations)
	router.GET("/getBuildVersionAndTimestamp", getBuildVersionAndTimestamp)
	router.POST("/authenticate", authenticate)

	router.Run(":" + prop.GetString("http-port", "8080"))
}

func upperStationIDs(stationIDs *[]string) {
	x := *stationIDs
	for i, stationID := range x {
		x[i] = strings.ToUpper(stationID)
	}
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
func pingDb(c *gin.Context) {
	stationIDs := model.SelectStationIDs()
	c.JSON(http.StatusOK, stationIDs)
}

func getStationIds(c *gin.Context) {
	stationIDs := model.SelectStationIDs()
	c.JSON(http.StatusOK, stationIDs)
}

func getMetarListInObervationTimeRange(c *gin.Context) {
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
}

func getMetarListForLatestNObservations(c *gin.Context) {
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
}

func getBuildVersionAndTimestamp(c *gin.Context) {
	buildVersion := prop.GetString("build-version", "N/A")
	buildTimestamp := prop.GetString("build-timestamp", "N/A")

	c.String(http.StatusOK, "%v_%v", buildVersion, buildTimestamp)
}

// From https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/#implementation-in-go
func authenticate(c *gin.Context) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(c.Request.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		c.JSON(http.StatusBadRequest, "Malformed header")
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		c.JSON(http.StatusUnauthorized, "User unauthorized")

		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	// user := User{
	// 	Username:  creds.Username,
	// 	Password:  "******",
	// 	FirstName: "Tarif",
	// 	LastName:  "Halabi",
	// 	Authorities: []struct {
	// 		Authority string `json:"authority"`
	// 	}{{Authority: "auth1"}, {Authority: "auth2"}},
	// 	Token: tokenString,
	// }
	user := User{
		Username:    creds.Username,
		Password:    "******",
		FirstName:   "Tarif",
		LastName:    "Halabi",
		Authorities: []authority{{Authority: "auth1"}, {Authority: "auth2"}},
		Token:       tokenString,
	}

	c.JSON(http.StatusOK, user)
	/*
		// Finally, we set the client cookie for "token" as the JWT we just generated
		// we also set an expiry time which is the same as the token itself
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

	*/
}

func setUpLogging() {
	if prop.GetString("log-file", "") != "" {
		// Set up lumberjack logging
		logFile := lumberjack.Logger{
			Filename:   prop.GetString("log-file", ""),
			MaxSize:    500, // megabytes
			MaxBackups: 30,
			MaxAge:     30,   //days
			Compress:   true, // disabled by default
		}
		log.SetOutput(&logFile)

		// Set gin output to lumberjack log
		gin.DefaultWriter = &logFile
	}
}
