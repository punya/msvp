package hello

import (
	"appengine"
	"appengine/search"
	"encoding/json"
	"net/http"
)

type Incident struct {
	appengine.GeoPoint
	Text string
}

type PostResult struct {
	DocId string
}

func init() {
	http.HandleFunc("/posts", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	c := appengine.NewContext(r)
	index, err := search.Open("incidents")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)

	switch r.Method {
	case "GET":
		var incidents []Incident
		for it := index.List(c, nil); ; {
			var incident Incident
			_, err := it.Next(&incident)
			if err == search.Done {
				break
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			incidents = append(incidents, incident)
		}
		encoder.Encode(incidents)

	case "POST":
		var incident Incident
		err := decoder.Decode(&incident)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := index.Put(c, "", &incident)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := PostResult{DocId: id}
		encoder.Encode(result)
	}
}
