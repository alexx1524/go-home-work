package internalhttp

import (
	"encoding/json"
	"net/http"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) InitializeEventsRoutes() {
	s.Router.HandleFunc("/event", s.createEvent).Methods("POST")
	s.Router.HandleFunc("/event", s.updateEvent).Methods("PUT")
	s.Router.HandleFunc("/event/{id}", s.getEvent).Methods("GET")
	s.Router.HandleFunc("/event/{id}", s.deleteEvent).Methods("DELETE")
	s.Router.HandleFunc("/event/search", s.searchEvents).Methods("POST")
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var event storage.Event
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := s.app.InsertEvent(ctx, event); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, event)
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid event ID")
		return
	}

	if err := s.app.DeleteEvent(ctx, id); err != nil {
		respondWithError(w, 1, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, id)
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var event storage.Event
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := s.app.UpdateEvent(ctx, event); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, event)
}

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid event ID")
		return
	}

	event, err := s.app.GetEventByID(ctx, id)
	if err != nil {
		respondWithError(w, 1, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, event)
}

func (s *Server) searchEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var searchCriteria EventsSearchCriteria
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&searchCriteria); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()

	events, err := s.app.GetEventsForPeriod(ctx, searchCriteria.StartDate, searchCriteria.EndDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
