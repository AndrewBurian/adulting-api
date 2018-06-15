package main

import (
	"net/http"

	"github.com/AndrewBurian/adulting-api/data"
	"github.com/AndrewBurian/adulting-api/middlewares"
	"github.com/AndrewBurian/powermux"
	"github.com/sirupsen/logrus"
)

// ActivityHandler is the router for user requests
type ActivityHandler struct {
	activityDB data.ActivityDAL
	userDB     data.UserDAL
	logger     *logrus.Entry
}

// ActivitiesResponse is the response for fetching all activites
type ActivitiesResponse struct {
	Debits  []*data.Activity `json:"debits"`
	Credits []*data.Activity `json:"credits"`
}

type ActivityDoneRepsonse struct {
	Points int `json:"points"`
}

// Setup creates routes for the user handler
func (h *ActivityHandler) Setup(r *powermux.Route) {
	r.GetFunc(h.GetAll)
	r.Route("/do/:id").
		MiddlewareFunc(middlewares.RequireAuth).
		PostFunc(h.DoActivity)
}

// GetAll returns the list of all activities available to the user
func (h *ActivityHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	log := h.logger

	activities, err := h.activityDB.GetActivites()
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

	if err := WriteResponse(w, r, resp, http.StatusOK); err != nil {
		log.WithError(err).Error("Unable to send response")
	}
}

// DoActivity Submits an action to be done
func (h *ActivityHandler) DoActivity(w http.ResponseWriter, r *http.Request) {

	act := &data.Activity{
		ID: powermux.PathParam(r, "id"),
	}

	log := h.logger.WithField("act_id", act.ID)

	err := h.activityDB.GetActivity(act)
	if err == data.ErrNotFound {
		http.Error(w, "Activity not found", http.StatusNotFound)
		log.Warn("Activity not found")
		return
	}
	if err != nil {
		http.Error(w, "Error retrieving activity", http.StatusInternalServerError)
		log.WithError(err).Error("Database Error")
		return
	}

	// user guarenteed by auth requirement
	user := middlewares.GetUser(r)

	/*
		We allow the user to dip into negative points, so long as they don't try to
		start with 0 or less points.

		Why? Not sure, sounded cool. Fun to watch people get into adult point debt
	*/
	if user.Points <= 0 && act.PointValue < 0 {
		resp := &ActivityDoneRepsonse{
			Points: user.Points,
		}
		log.WithField("points", user.Points).Debug("Attempt to overspend points")
		if err := WriteResponse(w, r, resp, http.StatusForbidden); err != nil {
			log.WithError(err).Error("Error writing response")
		}
		return
	}

	user.Points += act.PointValue
	if err := h.userDB.UpdatePoints(user); err != nil {
		http.Error(w, "Error updating points", http.StatusInternalServerError)
		log.WithError(err).Error("Database Error")
		return
	}

	resp := &ActivityDoneRepsonse{
		Points: user.Points,
	}

	log.WithField("points", user.Points).Debug("Updated user points")
	if err := WriteResponse(w, r, resp, http.StatusOK); err != nil {
		log.WithError(err).Error("Error writing response")
	}

}
