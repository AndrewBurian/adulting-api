package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AndrewBurian/mediatype"
)

var (
	jsonType, _ = mediatype.ParseSingle("application/json; chatset=utf-8")
)

// DecodeJSON reads an arbitrary JSON object from a request
func DecodeJSON(r *http.Request, obj interface{}) error {
	content := ContentType(r)
	if content == nil || content.SubType != "json" {
		return fmt.Errorf("Content Type not JSON")
	}

	return json.NewDecoder(r.Body).Decode(obj)
}

// WriteResponse attempts to write the JSON object into the response
func WriteResponse(w http.ResponseWriter, r *http.Request, obj interface{}) error {
	if !Accepts(r).SupportsType(jsonType) {
		http.Error(w, "Unsupported response content types", http.StatusNotAcceptable)
		return fmt.Errorf("Doesn't accept JSON")
	}

	w.Header().Add("Content-Type", jsonType.String())

	return json.NewEncoder(w).Encode(obj)
}
