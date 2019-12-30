package appremove

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

//AppRemoveEvent struct for app remove event record
type AppRemoveEvent struct {
	Platform  string     `bigquery:"platform"`
	EventName string     `bigquery:"event_name"`
	Device    DeviceInfo `bigquery:"device"`
	App       AppInfo    `bigquery:"app_info"`
	Timestamp int        `bigquery:"event_timestamp"`
	UserID    string     `bigquery:"user_pseudo_id"`
}

//DeviceInfo ...
type DeviceInfo struct {
	MobileBrandName string `bigquery:"mobile_brand_name"`
	AdvertisingID   string `bigquery:"advertising_id"`
	OS              string `bigquery:"operating_system"`
	OSVersion       string `bigquery:"operating_system_version"`
}

type AppInfo struct {
	ID      string `bigquery:"id"`
	Version string `bigquery:"version"`
}

const (
	projectID = "inapp-infrastructure-190215"
	topicID   = "app_remove_event"
	getQuery  = `SELECT platform,event_name,device FROM ` + "`inapp-infrastructure-190215.app_remove_events.events`" + `LIMIT 1000`
)

var (
	pubsubClient *pubsub.Client
	once         sync.Once
)

// func ReadFireStore1(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()
// 	fmt.Println("testing")
// 	fmt.Fprintf(w, "ReadFireStore")
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

// 	//iter := client.Collection("uninstall").Documents(ctx)
// 	iter := client.Collection("uninstall").Where("advertisingID", "==", "dd3cc431-6697-4d7d-972e-3126dae3fa2a").Documents(ctx)
// 	for {
// 		doc, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalf("Failed to iterate: %v", err)
// 		}
// 		fmt.Println(doc.Data())
// 		if url, ok := doc.Data()["url"]; ok {
// 			//makeHTTPGETRequest(url.(string))
// 			fmt.Println(url.(string))
// 		}
// 	}
// }

//GetAppRemoveEvents retrieve all app remove events from bigquery table
func GetAppRemoveEvents(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		errorDetails := fmt.Sprintf("Failed to initialize bigqeury client: %+v", err)
		log.Printf("Error: %s", errorDetails)
		http.Error(w, errorDetails, http.StatusInternalServerError)
		return
	}

	q := client.Query(fmt.Sprintf("%s", getQuery)) //TODO: create query according to date
	q.Location = "US"

	it, err := q.Read(ctx)
	if err != nil {
		errorDetails := fmt.Sprintf("Failed to read query: %v", err)
		log.Printf("Error: %s", errorDetails)
		http.Error(w, errorDetails, http.StatusInternalServerError)
		return
	}
	var wg sync.WaitGroup
	for {
		var event AppRemoveEvent
		err := it.Next(&event)
		if err == iterator.Done {
			break
		}
		if err != nil {
			errorDetails := fmt.Sprintf("Failed to iterate result: %v", err)
			log.Printf("Error: %s", errorDetails)
			//http.Error(w, errorDetails, http.StatusInternalServerError)
			continue
		}
		//log.Printf("Info: publising app_remove_event - %+v \n \n", event.Device.AdvertisingID)
		fmt.Println(event.Device.AdvertisingID)
		wg.Add(1)
		go publishAppRemoveEvent(w, event, &wg) //handle retries
	}
	wg.Wait()
	fmt.Fprintf(w, "Published message for AppRemoveEvents")
}

//To create pubsub client only once
func getPubSubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	var err error
	once.Do(func() {
		pubsubClient, err = pubsub.NewClient(ctx, projectID)
	})
	return pubsubClient, err
}

func publishAppRemoveEvent(w http.ResponseWriter, event AppRemoveEvent, wg *sync.WaitGroup) {
	var err error
	defer func() {
		wg.Done()
	}()
	ctx := context.Background()
	pubsubClient, err = getPubSubClient(ctx, projectID)
	if err != nil {
		log.Printf("Error: pubsub.NewClient: %v", err)
		return
	}

	t := pubsubClient.Topic(topicID)
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error: publish event, marshal error: %v", err)
		return
	}
	//fmt.Println(data)
	result := t.Publish(ctx, &pubsub.Message{
		Data: data,
		//Add attributes if required
		// Attributes: map[string]string{

		// },
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		fmt.Printf("Get: %v", err)
		return

	}
	log.Printf("Info: Published message with custom attributes; msg ID: %v\n", id)
	return
}
