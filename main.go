package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)
const(
	fileName ="poke-map-32809-a9da85a709d8.json"
)
var (

	// uploadableBucket is the destination bucket.
	// All users will upload files directly to this bucket by using generated Signed URL.
	uploadableBucket string
)


func generateV4GetObjectSignedURL(w http.ResponseWriter, r *http.Request)  {
	// bucket := "bucket-name"
	// object := "object-name"
	// serviceAccount := "service_account.json"
	jsonKey, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return
	}

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	key := r.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV2,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	u, err := storage.SignedURL(uploadableBucket, key, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Generated GET signed URL:")
	fmt.Fprintf(w, "%q\n", u)
	fmt.Fprintln(w, "You can use this URL with any user agent, for example:")
	fmt.Fprintf(w, "curl %q\n", u)
	return
}

func signHandler(w http.ResponseWriter, r *http.Request) {
	// Accepts only POST method.
	// Otherwise, this handler returns 405.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	ct := r.FormValue("content_type")
	if ct == "" {
		http.Error(w, "content_type must be set", http.StatusBadRequest)
		return
	}

	// Generates an object key for use in new Cloud Storage Object.
	// It's not duplicate with any object keys because of UUID.
	key := "butter-fly"
	if ext := r.FormValue("ext"); ext != "" {
		key += fmt.Sprintf(".%s", ext)
	}

	// bucket := "bucket-name"
	// object := "object-name"
	// serviceAccount := "service_account.json"
	jsonKey, err := ioutil.ReadFile(fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	opts := &storage.SignedURLOptions{
		Scheme: storage.SigningSchemeV4,
		Method: "PUT",
		Headers: []string{
			"Content-Type: image/jpeg",
			"x-goog-acl: public-read",
		},
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	u, err := storage.SignedURL(uploadableBucket, key, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Generated PUT signed URL:")
	fmt.Fprintf(w, "%q\n", u)
	fmt.Fprintln(w, "You can use this URL with any user agent, for example:")
	fmt.Fprintf(w, "curl -X PUT -H 'Content-Type: image/jpeg' -H 'x-goog-acl: public-read' --upload-file my-file %q\n", u)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, u)
}

func main() {
	uploadableBucket = "signed-bucket-sample"


	fmt.Printf("Started....\n")
	http.HandleFunc("/sign", signHandler)
	http.HandleFunc("/signed-get", generateV4GetObjectSignedURL)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":8080"), nil))
}
