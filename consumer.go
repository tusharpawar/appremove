package appremove

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HandleAppRemoveEvent handles a Pub/Sub message.
func HandleAppRemoveEvent(ctx context.Context, m PubSubMessage) error {
	log.Println(string(m.Data))
	var event AppRemoveEvent
	json.Unmarshal(m.Data, &event)

	log.Println(event.Device.AdvertisingID)
	ReadDataStore("inApp", event.Device.AdvertisingID)

	return nil
}

//ReadDataStore ..
func ReadDataStore(applicationName, advertisingID string) {
	// func ReadDataStore(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("ReadDatastore")
	// 	advertisingID := "2c96e4e7-593a-4112-82af-7e21d386b397" //testing purpose, to be commented
	// 	applicationName := "inApp"                              //testing purpose, to be commented
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		errorDetails := fmt.Sprintf("Failed to initialize datastore client: %+v", err)
		fmt.Printf("Error: %s", errorDetails)
		return
	}

	query := datastore.NewQuery("uninstall").
		Filter("advertising_id =", advertisingID).
		Filter("application_name =", applicationName).
		Limit(1)

	iter := client.Run(ctx, query)
	for {
		var res UninstallData
		_, err := iter.Next(&res)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		log.Println(res)
		makeHTTPGETRequest(res.URL)
	}

}

//ReadFireStore ...
func ReadFireStore(advertisingID string) {
	//func ReadFireStore(w http.ResponseWriter, r *http.Request) {
	//	advertisingID := "45ba7bc8-c2c4-41a3-97c5-8bd063e81849"
	ctx := context.Background()
	log.Println("Reading Firestore for advID: ", advertisingID)
	opt := option.WithCredentialsFile("./inapp-test-707fe-firebase-adminsdk-ibdxb-7f80f3f4f5.json") //path to the token
	conf := &firebase.Config{ProjectID: "inapp-test-707fe"}
	app, err := firebase.NewApp(ctx, conf, opt) //pass token here
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx) //create client only once
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//iter := client.Collection("uninstall").Documents(ctx)
	iter := client.Collection("uninstall").Where("advertisingID", "==", advertisingID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		log.Println(doc.Data())
		if url, ok := doc.Data()["url"]; ok {
			makeHTTPGETRequest(url.(string))
		}
	}
}

//FeedFireStore to populate values in firestore
// func FeedFireStore(w http.ResponseWriter, r *http.Request) {
// 	var adList []string
// 	adList = []string{"dd3cc431-6697-4d7d-972e-3126dae3fa2a", "08f58c0b-189f-4cc7-8adb-af39d712b039", "27807cff-6ea4-4363-8662-3792bc101d24", "08f58c0b-189f-4cc7-8adb-af39d712b039", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "08f58c0b-189f-4cc7-8adb-af39d712b039", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "08f58c0b-189f-4cc7-8adb-af39d712b039", "08f58c0b-189f-4cc7-8adb-af39d712b039", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "45ba7bc8-c2c4-41a3-97c5-8bd063e81849", "069c0408-38af-45c2-ab49-ec2421e61d10", "ef45fe6b-09d5-4362-b99c-606ecf4e4f39", "a93c5d26-47d3-4e1c-8f12-0bed49141bd2", "bcf175a4-5d81-4155-a74f-074a05f713e6", "e68d28e6-32a1-4781-8aaf-983244a960c6", "dca621f8-5092-498b-811e-04ac20a30505", "2c96e4e7-593a-4112-82af-7e21d386b397", "e6e3cc86-7d20-4f12-9ae7-8f11d142d7f1", "ccc2d2a0-6553-4939-823d-56e83f8b84c8", "9afd5074-7d10-4e6c-8a05-41152a3e3ae6", "0d67c3b3-3107-400a-b3f5-d8799c852a8c", "bb16c1a9-9f27-4685-a46d-0ccac11c318d", "2c96e4e7-593a-4112-82af-7e21d386b397", "70dbe428-e3ad-4bbd-8835-58b83d197aa0", "cec130d7-7a61-48fc-80b6-9fda6233c8d7", "31804649-74bd-417e-94f8-baee272bf24a", "62dc3ca6-6ecc-4e15-8e12-d6e842a79940", "abfd6a08-9280-46e8-9312-7bf5d0102a24", "a8b01450-bbc5-44d7-935b-618e7904d4ab", "ea2fb4e6-57ec-4312-a49b-691734984055", "5f665f1c-e378-4a3e-95dd-f932b696055e", "2da94180-5c39-4170-85fa-2ddd0b4cceec", "f1cc7187-328d-4a6d-9962-fc54e24caded", "908bab90-61ed-4683-83f7-a97289a59991", "3fdbdf83-0515-47e0-9d5c-bc324ee9d400"}
// 	ctx := context.Background()
// 	log.Println("Writing data to firestore")
// 	opt := option.WithCredentialsFile("./inapp-test-707fe-firebase-adminsdk-ibdxb-7f80f3f4f5.json") //path to the token
// 	conf := &firebase.Config{ProjectID: "inapp-test-707fe"}
// 	app, err := firebase.NewApp(ctx, conf, opt) //pass token here
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	client, err := app.Firestore(ctx) //create client only once
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer client.Close()

// 	for _, ad := range adList {
// 		_, _, err := client.Collection("uninstall").Add(ctx, map[string]interface{}{
// 			"advertisingID": ad,
// 			"url":           "https://www.google.com/",
// 		})
// 		if err != nil {
// 			log.Println("Error occured to write data to firestore")
// 		}
// 	}

// }

func makeHTTPGETRequest(url string) error {
	log.Println("Calling url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("GET request failed, err: %+v", err)
		return err
	}

	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("Failed to parse response body, err: %+v", err)
	// 	return err
	// }

	log.Println("status: ", resp.StatusCode)
	return nil
}
