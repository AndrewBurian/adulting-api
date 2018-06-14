package main

import (
	"net/http"

	"github.com/andrewburian/adulting-api/data"
	"github.com/andrewburian/powermux"
	"github.com/sirupsen/logrus"
)

// ActivityHandler is the router for user requests
type ActivityHandler struct {
	db     data.ActivityDAL
	logger *logrus.Entry
}

type ActivitiesResponse struct {
	Debits  []*data.Activity `json:"debits"`
	Credits []*data.Activity `json:"credits"`
}

// Setup creates routes for the user handler
func (h *ActivityHandler) Setup(r *powermux.Route) {
	r.GetFunc(h.GetAll)
}

func (h *ActivityHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	log := h.logger

	activities, err := h.db.GetActivites()
	if err != nil {
		log.WithError(err).Error("Database error")
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	resp := &ActivitiesResponse{
		Debits:  make([]*data.Activity, 0, 5),
		Credits: make([]*data.Activity, 0, 5),
	}

	for i := range activities {
		if activities[i].PointValue >= 0 {
			resp.Debits = append(resp.Debits, activities[i])
		} else {
			resp.Credits = append(resp.Credits, activities[i])
		}
	}

	if err := WriteResponse(w, r, resp); err != nil {
		log.WithError(err).Error("Unable to send response")
	}
}
