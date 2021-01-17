package signer

import (
	"cloud.google.com/go/storage"
	"fmt"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"io/ioutil"
	"time"
)

type (
	SignedRequest struct {
		FileName string `json:"fileName"`
		ContentType string  `json:"contentType"`
		PublicRead bool `json:"public-read"`
	}

	GcpSigner struct {
		conf *jwt.Config
		bucket string
	}
)

func NewGcpSigner(googleCredentialJsonFile, bucket string ) *GcpSigner{
	jsonKey, err := ioutil.ReadFile(googleCredentialJsonFile)
	if err != nil {
		panic(err)
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
	panic(err)
	}

	return &GcpSigner{conf: conf, bucket: bucket}
}

func (gcp *GcpSigner) UploadSigned( req *SignedRequest) (string ,error) {
	headers:= make([]string,0)

	if req.ContentType== ""{
		headers = append(headers, "Content-Type: application/octet-stream")
	}else{
		headers = append(headers, fmt.Sprintf("Content-Type: %v", req.ContentType))
	}

	if req.PublicRead{
		headers = append(headers, "x-goog-acl: public-read")
	}

	opts := &storage.SignedURLOptions{
		Scheme: storage.SigningSchemeV4,
		Method: "PUT",
		Headers:headers,
		GoogleAccessID: gcp.conf.Email,
		PrivateKey:     gcp.conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	return  storage.SignedURL(gcp.bucket, req.FileName, opts)
}

func (gcp *GcpSigner) GenerateV4GetObjectSignedURL(req *SignedRequest)  (string,error){
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV2,
		Method:         "GET",
		GoogleAccessID: gcp.conf.Email,
		PrivateKey:     gcp.conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	return storage.SignedURL(gcp.bucket, req.FileName, opts)
}
