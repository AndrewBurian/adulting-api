package main

import (
	"context"
	"net/http"

	"github.com/AndrewBurian/mediatype"
	log "github.com/sirupsen/logrus"
)

type contentTypeCtx string

var (
	contentTypeCtxKey = contentTypeCtx("content-type")
	acceptsTypeCtxKey = contentTypeCtx("accepts")
)

// ContentTypeDetect is a middleware that adds content types to the reqest
func ContentTypeDetect(w http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
	contentType, accepts, err := mediatype.ParseRequest(r)
	if err != nil {
		log.WithField("middleware", "content-type-detect").WithError(err).Warn("Invalid Content Type")
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	ctx = context.WithValue(ctx, contentTypeCtxKey, contentType)
	ctx = context.WithValue(ctx, acceptsTypeCtxKey, accepts)
	r = r.WithContext(ctx)

	next(w, r)
}

// ContentType returns the content type of the request body
func ContentType(r *http.Request) *mediatype.ContentType {
	return r.Context().Value(contentTypeCtxKey).(*mediatype.ContentType)
}

// Accepts returns the content type of the request body
func Accepts(r *http.Request) mediatype.ContentTypeList {
	return r.Context().Value(acceptsTypeCtxKey).(mediatype.ContentTypeList)
}
