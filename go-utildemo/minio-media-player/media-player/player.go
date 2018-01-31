package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go"
)

var (
	bucketName = flag.String("b", "", "Bucket name to be used for media assets.")
	endPoint   = flag.String("e", "http://10.39.35.38:9000", "Choose a custom endpoint.")
)

type mediaPlayList struct {
	Key string
	URL string
}

type mediaHandlers struct {
	minioClient *minio.Client
}

var supportedAccessEnvs = []string{
	"ACCESS_KEY",
	"AWS_ACCESS_KEY",
	"AWS_ACCESS_KEY_ID",
}

var supportedSecretEnvs = []string{
	"SECRET_KEY",
	"AWS_SECRET_KEY",
	"AWS_SECRET_ACCESS_KEY",
}

func mustGetAccessKeys() (accessKey, secretKey string) {
	for _, accessKeyEnv := range supportedAccessEnvs {
		accessKey = os.Getenv(accessKeyEnv)
		if accessKey != "" {
			break
		}
	}

	for _, secretKeyEnv := range supportedSecretEnvs {
		secretKey = os.Getenv(secretKeyEnv)
		if secretKey != "" {
			break
		}
	}

	if accessKey == "" {
		accessKey = "PYKBI9S2IH1JV00GH5KS"
	}

	if secretKey == "" {
		secretKey = "ReTLfTEDGI3YWRwD9kbKagUxYFTZTPGkSXpmgQk8"
	}
	return accessKey, secretKey
}

func isSecure(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return u.Scheme == "https"
}

func findHost(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return u.Host
}

func main() {
	flag.Parse()

	if *bucketName == "" {
		log.Fatalln("Bucket name cannot be empty.")
	}

	// Fetch access keys if possible or fail.
	accessKey, secretKey := mustGetAccessKeys()
	minioClient, err := minio.New(findHost(*endPoint), accessKey, secretKey, isSecure(*endPoint))
	if err != nil {
		log.Fatalln(err)
	}

	mediaPlayer := mediaHandlers{
		minioClient: minioClient,
	}
	// Handler to serve the index page.
	http.Handle("/", http.FileServer(assetFS()))

	// End point for list object operations.
	// Called when player in the front end is initialized.
	http.HandleFunc("/list/v1", mediaPlayer.ListObjectsHandler)

	// Given point which receives the object name and returns presigned URL in the response.
	http.HandleFunc("/getpresign/v1", mediaPlayer.GetPresignedURLHandler)

	log.Println("Starting media player, please visit your browser at http://localhost:9001")

	// Port is defaulted to "9001" no need to change this.
	http.ListenAndServe(":9001", nil)
}

// ListObjectsHandler - handler for ListsObjects from the Object Storage server and bucket configured above.
func (api mediaHandlers) ListObjectsHandler(w http.ResponseWriter, r *http.Request) {
	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	var playListEntries []mediaPlayList

	// Tracks if first object presigned.
	var firstObjectPresigned bool

	// Set recursive to list all objects.
	var isRecursive = true

	// List all objects from a bucket-name with a matching prefix.
	for objectInfo := range api.minioClient.ListObjects(*bucketName, "", isRecursive, doneCh) {
		if objectInfo.Err != nil {
			http.Error(w, objectInfo.Err.Error(), http.StatusInternalServerError)
			return
		}
		objectName := objectInfo.Key // object name.
		playListEntry := mediaPlayList{
			Key: objectName,
		}
		if !firstObjectPresigned {
			// Generating presigned url for the first object in the list.
			// presigned URL will be generated on the fly for the
			// other objects when they are played.
			expirySecs := 1000 * time.Second // 1000 seconds.
			presignedURL, err := api.minioClient.PresignedGetObject(*bucketName, objectName, expirySecs, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			playListEntry.URL = presignedURL.String()
		}
		playListEntries = append(playListEntries, playListEntry)
	}
	playListEntriesJSON, err := json.Marshal(playListEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Successfully wrote play list in json.
	w.Write(playListEntriesJSON)
}

// GetPresignedURLHandler - generates presigned access URL for an object.
func (api mediaHandlers) GetPresignedURLHandler(w http.ResponseWriter, r *http.Request) {
	// The object for which the presigned URL has to be generated is sent as a query
	// parameter from the client.
	objectName := r.URL.Query().Get("objName")
	if objectName == "" {
		http.Error(w, "No object name set, invalid request.", http.StatusBadRequest)
		return
	}
	presignedURL, err := api.minioClient.PresignedGetObject(*bucketName, objectName, 1000*time.Second, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(presignedURL.String()))
}
