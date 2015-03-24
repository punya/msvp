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
	Lat      float64 `datastore:",noindex" json:"lat"`
	Lng      float64 `datastore:",noindex" json:"lng"`
	Text     string  `datastore:",noindex" json:"text"`
	Verified bool    `json:"verified"`
}

type IncidentWithKey struct {
	Incident
	Key int64 `datastore:"-" json:"key"`
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

	verified := true
	if r.FormValue("verified") == "false" && user.IsAdmin(c) {
		verified = false
	}

	query := datastore.NewQuery("Incident").Filter("Verified =", verified)

	incidents := []IncidentWithKey{}
	keys, err := query.GetAll(c, &incidents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
