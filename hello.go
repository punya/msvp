package hello

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile("/incidents/(\\d+)")

type Incident struct {
	appengine.GeoPoint
	Text     string
	Verified bool
	Key      int64 `datastore:"-"`
}

type PostResult struct {
	DocId string
}

func init() {
	http.HandleFunc("/incidents", dispatch)
	http.HandleFunc("/incidents/", dispatch)
}

func dispatch(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	var err error

	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path == "/incidents" {
		switch r.Method {
		case "GET":
			if incidents, err := getAllIncidents(c); err == nil {
				encoder.Encode(incidents)
				return
			}
		case "POST":
			var incident Incident
			if err := decoder.Decode(&incident); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			addIncident(c, incident)
			return
		default:
			http.Error(w, "Invalid method", http.StatusBadRequest)
			return
		}
	}
	if match := re.FindStringSubmatch(r.URL.Path); match != nil {
		key, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		switch r.Method {
		case "PUT":
			var incident Incident
			if err = decoder.Decode(&incident); err != nil {
				break
			}
			incident.Key = key
			err = saveIncident(c, incident)
		default:
			http.Error(w, "Invalid method", http.StatusBadRequest)
			return
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

func getAllIncidents(c appengine.Context) (incidents []Incident, err error) {
	query := datastore.NewQuery("Incident")
	if !user.IsAdmin(c) {
		query = query.Filter("Verified =", true)
	}
	keys, err := query.GetAll(c, &incidents)
	if err != nil {
		return
	}
	for i, _ := range keys {
		incidents[i].Key = keys[i].IntID()
	}
	return
}

func addIncident(c appengine.Context, incident Incident) (err error) {
	if !user.IsAdmin(c) {
		incident.Verified = false
	}

	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Incident", nil), &incident)
	return
}

func saveIncident(c appengine.Context, incident Incident) (err error) {
	if !user.IsAdmin(c) {
		return errors.New("unauthorized")
	}
	_, err = datastore.Put(c, datastore.NewKey(c, "Incident", "", incident.Key, nil), &incident)
	return
}
