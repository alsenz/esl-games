package main

import (
    "bytes"
    "flag"
    "fmt"
    "github.com/alsenz/esl-games/pkg/account"
    "io"
    "net/http"

    "github.com/alsenz/esl-games/pkg/assetserver"
)

//TODO how to do the request type?

//TODO this probably needs to be moved around.
func PostAsset(w http.ResponseWriter, r *http.Request) {
    user, err := account.CheckAuth(r) //TODO wanna pass a string array of required groups...
    if err != nil {
        //TODO we wanna add a logger here for the error...
        //TODO write an error... TODO TODO

        //TODO wanna check if it's an unauthorised.
    }
    r.ParseMultipartForm(32 << 20) // limit your max input length!
    var buf bytes.Buffer
    file, header, err := r.FormFile("asset")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    io.Copy(&buf, file)

    // do something with the contents...
    // I normally have a struct defined and unmarshal into a struct, but this will
    // work as an example
    contents := buf.String()
    fmt.Println(contents)
    // I reset the buffer in case I want to use it again
    // reduces memory allocations in more intense projects
    buf.Reset()
    // do something else
    // etc write header
    return
}

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
