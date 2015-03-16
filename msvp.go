package msvp

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Incident struct {
	Lat      float64 `datastore:lat`
	Lng      float64 `datastore:lng`
	Text     string  `datastore:text`
	Verified bool    `datastore:verified`
}

type IncidentWithKey struct {
	Incident
	Key int64 `datastore:"-"`
}

func init() {
	router := mux.NewRouter()
	router.HandleFunc("/incidents", getIncidents).Methods("GET")
	router.HandleFunc("/incidents", addIncident).Methods("POST")
	router.HandleFunc(`/incidents/{id:\d+}`, updateIncident).Methods("PUT")

	http.Handle("/", router)
}

func getIncidents(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	query := datastore.NewQuery("Incident")
	if !user.IsAdmin(c) {
		query = query.Filter("Verified =", true)
	}

	incidents := []IncidentWithKey{}
	keys, err := query.GetAll(c, &incidents)
	if err != nil {
		return
	}

	for i, _ := range keys {
		incidents[i].Key = keys[i].IntID()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incidents)
}

func addIncident(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	decoder := json.NewDecoder(r.Body)

	var incident Incident
	if err := decoder.Decode(&incident); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !user.IsAdmin(c) {
		incident.Verified = false
	}

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Incident", nil), &incident)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

func updateIncident(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !user.IsAdmin(c) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	key, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var incident Incident
	if err := decoder.Decode(&incident); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = datastore.Put(c, datastore.NewKey(c, "Incident", "", key, nil), &incident)
	return
}
