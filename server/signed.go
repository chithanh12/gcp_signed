package server

import (
	"fmt"
	"github.com/chithanh12/gcp_signed/signer"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"path/filepath"
)

type (
	GetRequest struct {
		Uuid string `uuid`
	}

	PutRequest struct {
		ContentType string `json:"contentType"`
		FileName string `json:"fileName"`
		PublicRead bool `json:"publicRead"`
	}
	
	FileResponse struct {
		Uuid string `json:"uuid"`
		Privacy string `json:"privacy"`
		FileName string `json:"fileName"` //original filename --> for display in UI
		Key string `json:"key"`
		SignedUrl string `json:"url"`
		ContentType string `json:"contentType"`
	}
)

func (s *Server) SignedUpload(e echo.Context) error{
	var req PutRequest
	if err:= e.Bind(&req);err != nil{
		return  e.JSON(http.StatusBadRequest, map[string]interface{}{"code": "badRequest", "message": "Invalid request"})
	}


	if req.FileName == ""{
		return  e.JSON(http.StatusBadRequest, map[string]interface{}{"code": "badRequest", "message": "Missing filename"})
	}

	containerId:= uuid.New().String()
	fileUuid:= uuid.New().String()
	ext:= filepath.Ext(req.FileName)

	newFilePath:= fmt.Sprintf("%s/%s%s", containerId, fileUuid, ext)

	signedReq:= &signer.SignedRequest{
		Key:         newFilePath,
		ContentType: req.ContentType,
		PublicRead:  req.PublicRead,
	}

	u, err:= s.signer.UploadSigned(signedReq)
	if err != nil{
		return e.JSON(http.StatusInternalServerError, map[string]interface{}{"code":"internalServerError", "message": "Can not signedurl"} )
	}
	privacy:= "private"
	if req.PublicRead{
		privacy = "public"
	}

	resp:= FileResponse{
		Uuid:     fileUuid   ,
		Privacy:     privacy,
		FileName:    req.FileName,
		Key: newFilePath,
		SignedUrl:   u,
		ContentType: req.ContentType,
	}
	// TODO: save info db
	//
	
	return e.JSON(http.StatusOK, resp)
}

func (s *Server) SignedGet(e echo.Context) error{
	var req GetRequest
	if err := e.Bind(e); err !=nil{
		return  e.JSON(http.StatusBadRequest, map[string]interface{}{"code": "badRequest", "message": "Invalid request"})
	}

	//TODO: get from db to get fileKey from req.Uuid
	db:=  &FileResponse{
		Key:         "....",
	}
	//

	signedReq:= &signer.SignedRequest{
		Key: db.Key,
	}
	u, err:= s.signer.GenerateV4GetObjectSignedURL(signedReq)
	if err != nil{
		return e.JSON(http.StatusInternalServerError, map[string]interface{}{"code":"internalServerError", "message": "Can not signedurl"} )
	}

	return e.JSON(http.StatusOK, map[string]interface{}{"uuid": req.Uuid, "url":u })
}