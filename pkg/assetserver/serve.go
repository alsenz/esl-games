package assetserver

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// Here we put the UNAUTHENTICATED serve functions for serving assets from not-necessarily logged in users of the apps

func (as *AssetServer) ServeAsset(w http.ResponseWriter, r *http.Request) {
	// We don't check the user name or auth, we just ensure the sig is correct
	if sig, had_sig := r.URL.Query()["sig"]; had_sig {
		//Get the end of the path
		assetId := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		//Strip off any extensions - these are optional (and in fact shouldn't be provided since the content-type
		//... will just be whatever is in the database
		assetId = strings.Split(assetId, ".")[0]
		if len(assetId) == 0 {
			http.Error(w, "Asset id (query path) must have length > 0", 400)
			return
		}

		asset := &Asset{}
		result:= as.conn.Where("id = ? AND md5sum = ?", assetId, sig).First(&asset)
		if result.Error == nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Either asset not found or sig incorrect", 404)
			} else {
				//TODO go logging -- how to log the unknown error here?
				//TODO check what the controller does... TODO TODO
			}
			return
		}

		//TODO let's write the asset back


	} else {
		http.Error(w, "Request must include a sig query parameter", 400)
	}
}
