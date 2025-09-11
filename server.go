package main

import (
	"context"

	"sync"

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	
	defer cancel()

	if err := database.ConnectPostgresDB(config.PostgresConfig); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected succesfuly to postgres")

	defer database.PostgresDB.Close()

	if err := database.CreatePackagesTable(ctx, database.PostgresDB); err != nil {
		log.Fatal(err)

	}
	fmt.Println("Created the packages table")

	repos := []requests.Request{
		{Owner: "snowflakedb", Repo: "snowpark-python"},
		{Owner: "hashicorp", Repo: "terraform-provider-aws"},
		{Owner: "dlt-hub", Repo: "dlt"},
		{Owner:"snowflakedb", Repo: "terraform-provider-snowflake"},
		{},
		
	
	}

	jobs := make(chan requests.Request, len(repos))
	githubChan := make(chan requests.GithubReleaseWithRepo, len(repos))
	var wg sync.WaitGroup

	// Start fetcher workers
	numFetchers := 5
	for i := 0; i < numFetchers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for req := range jobs {
				release, err := requests.ChangeLog(req)
				if err != nil {
					log.Println("Error fetching release:", err)
					continue
				}
				githubChan <- requests.GithubReleaseWithRepo{
					Owner: req.Owner,
					Repo: req.Repo,
					Release: release,
				}
			}
		}()
	}

	// Start DB writer workers
	numDBWriters := 3
	for i := 0; i < numDBWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for release := range githubChan {
				processRelease(ctx, release)
			}
		}()
	}

	
	for _, repo := range repos {
		jobs <- repo
	}
	close(jobs)

	// Wait for fetchers to finish, then close githubChan
	go func() {
		wg.Wait()
		close(githubChan)
	}()

	// Wait for all workers
	wg.Wait()
}


func processRelease(ctx context.Context, release requests.GithubReleaseWithRepo) {
	response := models.Response{
		Owner:     release.Owner,
		Repo:      release.Repo,
		Version:   release.Release.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	request := requests.Request{
		Owner: release.Owner,
		Repo: release.Repo,
	}
	checkVersion, err := database.GetVersionByOwnerAndRepo(ctx, database.PostgresDB,request )
	if err != nil {
		id, err := database.InsertIntoPostgres(ctx, database.PostgresDB, response)
		if err != nil {
			log.Println("Insert error:", err)
			return
		}
		fmt.Println("Inserted ID:", id)
		checkVersion, err = database.GetVersionByOwnerAndRepo(ctx, database.PostgresDB, request)
		if err != nil {
			log.Println("Error fetching version after insert:", err)
			return
		}
	}

	update := models.VersionUpdate{
		Owner:     response.Owner,
		Repo:      response.Repo,
		Version:   response.Version,
		UpdatedAt: time.Now(),
	}

	if err := database.UpdateVersion(ctx, database.PostgresDB, checkVersion.Version, update); err != nil {
		fmt.Printf("Update error on repo %s: %s",update.Repo, err)
		return
	}

	fmt.Printf("Updated version owner %s repo %s version %s\n",
		response.Owner, response.Repo, response.Version)
}