package main

import (
  "fmt"

  "github.com/alsenz/esl-games/pkg/assetserver"
)

func main() {
    fmt.Println("hello world")
    var asset = assetserver.Asset{}
    fmt.Println("Default asset content type is:", asset.ContentType)
}
