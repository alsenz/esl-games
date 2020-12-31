package assetserver

type Asset struct {
  Path string
  ContentType string
  Data string //TODO: convert this to a byte[] dept on GORM conventions/rules
}
