package assetserver

import (
  "github.com/alsenz/esl-games/pkg/account"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

type Asset struct {
  account.UserObject
  Name string         `gorm:"index"`
  Md5sum string       `gorm:"index"`
  ContentType string  `gorm:"index"`
  Data []byte
}

// This is a thread safe asset server- no state that isn't already thread safe.
type AssetServer struct {
  conn *gorm.DB
}

//TODO note we're just going to re-use the asset server to serve css too..

func NewAssetServer(pgConnection string) *AssetServer {
  assetServer := AssetServer{}
  db, err := gorm.Open(postgres.Open(pgConnection), &gorm.Config{})
  if err != nil {
    panic("Failed to connect to the database")
  }
  db.AutoMigrate(&account.User{})
  db.AutoMigrate(&account.Group{})
  db.AutoMigrate(&Asset{})
  assetServer.conn = db

  return &assetServer
}