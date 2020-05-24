package model

import (
	"log"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	"github.com/magiconair/properties"
)

// Db handle exported
var Db *sqlx.DB

// InitDB func
func InitDB(prop *properties.Properties) {
	var err error
	Db, err = sqlx.Open("godror", prop.GetString("username", "")+"/"+prop.GetString("password", "")+"@"+prop.GetString("connection-string", ""))
	if err != nil {
		log.Panic(err)
	}

	if err = Db.Ping(); err != nil {
		log.Panic(err)
	}

	//Db.Mapper = reflectx.NewMapperFunc("json", oracleColumnNameMapper)
	Db.MapperFunc(oracleColumnNameMapper)
}

func oracleColumnNameMapper(columnName string) string {
	if strings.ToUpper(columnName) == columnName { // Do nothing if already uppercase
		return columnName
	}
	const underscore = '_'
	var returnColumnName string
	for i, char := range columnName {
		//log.Printf("%#U", char)
		if unicode.IsUpper(char) && i != 0 {
			returnColumnName += string(underscore)
		}
		returnColumnName += string(char)
	}
	return strings.ToUpper(returnColumnName)
}
