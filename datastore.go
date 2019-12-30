package appremove

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
)

//UninstallData ...
type UninstallData struct {
	ApplicationName string         `datastore:"application_name",json:"applicationName"`
	AdvertisingID   string         `datastore:"advertising_id",json:"advertisingID"`
	URL             string         `datastore:"url",json:"url"`
	K               *datastore.Key `datastore:"__key__"`
}

// func (ud *UninstallData) LoadKey(k *datastore.Key) error {
// 	ud.K = k
// 	return nil
// }

// func (ud *UninstallData) Load(ps []datastore.Property) error {
// 	return datastore.LoadStruct(ud, ps)
// }

// func (ud *UninstallData) Save() ([]datastore.Property, error) {
// 	return datastore.SaveStruct(ud)
// }

//SaveUninstallData ...
func SaveUninstallData(w http.ResponseWriter, r *http.Request) {
	var uninstallData UninstallData

	err := json.NewDecoder(r.Body).Decode(&uninstallData)
	if err != nil {
		errorDetails := fmt.Sprintf("Failed to parse req body: %+v", err)
		log.Printf("Error: %s", errorDetails)
		http.Error(w, errorDetails, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		errorDetails := fmt.Sprintf("Failed to initialize datastore client: %+v", err)
		log.Printf("Error: %s", errorDetails)
		http.Error(w, errorDetails, http.StatusInternalServerError)
		return
	}

	key := datastore.IncompleteKey("uninstall", nil)
	log.Printf("Storing uninstallData : %+v", uninstallData)
	_, err = client.Put(ctx, key, &uninstallData)
	if err != nil {
		errorDetails := fmt.Sprintf("Failed to save uninstall data to datastore: %+v", err)
		log.Printf("Error: %s", errorDetails)
		http.Error(w, errorDetails, http.StatusInternalServerError)
		return
	}
	//}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(uninstallData)
	//fmt.Fprintf(w, "%+v", uninstallData)

}
