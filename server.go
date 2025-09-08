package main

import (
	"context"

	"fmt"
	"log"
	"time"

	"github.com/Nnamdichukwu/go-fortune/config"
	"github.com/Nnamdichukwu/go-fortune/database"
	"github.com/Nnamdichukwu/go-fortune/models"
	"github.com/Nnamdichukwu/go-fortune/requests"
)

func main() {
	if err := config.LoadEnvVars(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Env Vars loaded")

	ctx := context.Background()

	if err := database.ConnectPostgresDB(config.PostgresConfig); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected succesfuly to postgres")
	defer database.PostgresDB.Close()
	
	if err := database.CreatePackagesTable(ctx, database.PostgresDB); err != nil{
		log.Fatal(err)

	}
	fmt.Println("Created the packages table")
	request := requests.Request{
		Owner: "snowflakedb", 
		Repo: "snowpark-python"}

	
	release, err := requests.ChangeLog(request)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("The release version is %s\n", release.Name)

	version := release.Name
	response := models.Response{
		Owner: request.Owner,
		Repo: request.Repo,
		Version: version,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	checkVersion, err := database.GetVersionByOwnerAndRepo(ctx, database.PostgresDB, request)
	
	if err != nil {
		
		id, err := database.InsertIntoPostgres(ctx, database.PostgresDB, response)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id)
		checkVersion, err = database.GetVersionByOwnerAndRepo(ctx, database.PostgresDB, request)
		if err != nil {
			log.Fatal(err)
		}

		
	}
	fmt.Println(checkVersion)
	update := models.VersionUpdate{
		Owner: response.Owner,
		Repo: response.Repo,
		Version: response.Version,
		UpdatedAt: time.Now(),
	}
	fmt.Println(update)
	if err := database.UpdateVersion(ctx, database.PostgresDB,checkVersion.Version,update);err != nil {
		log.Fatal(err)

	}
	fmt.Printf("Updated version  owner %s repo %s version %s ", response.Owner, response.Repo, response.Version)


}
