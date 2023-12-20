package main

import (
	"fmt"
	"log"
	"os"

	db "github.com/pal-paul/pg-sql/pkg/database"
	env "github.com/pal-paul/pg-sql/pkg/env"
	utils "github.com/pal-paul/pg-sql/pkg/utils"
)

type Environment struct {
	Database struct {
		DBUser     string `env:"INPUT_DB_USER,required=true"`
		DBPassword string `env:"INPUT_DB_PASSWORD,required=true"`
		DBHost     string `env:"INPUT_DB_HOST,required=true"`
		DBPort     int    `env:"INPUT_DB_PORT,required=true"`
		DB         string `env:"INPUT_DB,required=true"`
	}
	Debug      bool   `env:"DEBUG"`
	ScriptsDir string `env:"INPUT_SCRIPTS_DIR,required=true"`
}

var err error
var envVar Environment

var sqlDb db.DatabaseInterface
var dbCredentials db.DBCredentials
var util utils.UtilsInterface

// initialize environment variables
func init() {
	_, err := env.Unmarshal(&envVar)
	if err != nil {
		log.Fatal(err)
	}
	util = utils.New()
}

func main() {
	sqlDb = db.New()
	dbCredentials.DBHost = envVar.Database.DBHost
	dbCredentials.DBPort = envVar.Database.DBPort
	dbCredentials.DBUser = envVar.Database.DBUser
	dbCredentials.DBPassword = envVar.Database.DBPassword
	dbCredentials.DB = envVar.Database.DB

	err = sqlDb.Connect(dbCredentials)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	sqlFiles, err := util.FilePathWalkDir(envVar.ScriptsDir, ".sql")
	if err != nil {
		log.Fatal(err)
	}
	for _, sqlFile := range sqlFiles {
		if envVar.Debug {
			log.Println(fmt.Println("executing sql file: ", sqlFile))
		}
		scripts, err := os.ReadFile(sqlFile)
		if err != nil {
			log.Fatalf("failed to read file: %s", sqlFile)
		}
		sql := string(scripts)
		err = sqlDb.Exec(sql)
		if err != nil {
			log.Fatalf("failed to execute sql file: %s", sqlFile)
		}
	}
}
