package assetserver

import (
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

type Asset struct {
  gorm.Model
  Path string
  ContentType string
  Data string //TODO: convert this to a byte[] dept on GORM conventions/rules
}

type AssetServer struct {
  
}