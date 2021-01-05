package main

import (
    "flag"
    "fmt"
    "github.com/alsenz/esl-games/pkg/assetserver"
    "net/http"
)

//TODO let's sub in mux here! TODO TODO

//TODO

func main() {
    pgHost := flag.String("pg-host", "localhost", "Postgresql database host")
    pgPort := flag.String("pg-port", "5432", "Postgresql database port")
    pgUser := flag.String("pg-user", "postgres", "Postgresql user")
    pgPass := flag.String("pg-password", "admin", "Postgresql user password")
    pgDb := flag.String("pg-database", "esl_games", "Postgresql database")
    asAddr := flag.String("listen", "0.0.0.0:8080", "Address for server to listen on")
    flag.Parse()

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", *pgHost, *pgUser, *pgPass, *pgDb, *pgPort)
    fmt.Println(dsn)
    as := assetserver.NewAssetServer(dsn)
    http.HandleFunc("/", as.ServeAsset)
    http.ListenAndServe(*asAddr, nil)
}
