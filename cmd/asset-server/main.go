package main

import (
    "flag"
    "fmt"

    "github.com/alsenz/esl-games/pkg/assetserver"
)

func main() {
    pgHost := flag.String("pg-host", "localhost", "Postgresql database host")
    pgPort := flag.String("pg-port", "5432", "Postgresql database port")
    pgUser := flag.String("pg-user", "postgres", "Postgresql user")
    pgPass := flag.String("pg-password", "admin", "Postgresql user password")
    pgDb := flag.String("pg-database", "esl_games", "Postgresql database")
    flag.Parse()

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", *pgHost, *pgUser, *pgPass, *pgDb, *pgPort)
    fmt.Println(dsn)
    _ = assetserver.NewAssetServer(dsn)
}
