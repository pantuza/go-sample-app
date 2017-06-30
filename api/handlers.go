package api

import (
	"encoding/json"
	"github.com/dimfeld/httptreemux"
	"github.com/pantuza/go-sample-app/db"
	"github.com/pantuza/go-sample-app/music"
	"log"
	"net/http"
)

var repository *db.MusicRepository

type GetMusic struct{}

func (m *GetMusic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := httptreemux.ContextParams(r.Context())
	log.Printf("GET music: %s", params["name"])

	music_obj, err := repository.Get(params["name"])

	switch err {
	case db.ErrMusicNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case db.ErrMusicDecode:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(music_obj)
}

type CreateMusic struct{}

func (m *CreateMusic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := httptreemux.ContextParams(r.Context())
	log.Printf("PUT music: %s", params["name"])

	music := &music.Music{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(music)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	repository.Create(music)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(music)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

type UpdateMusic struct{}

func (m *UpdateMusic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := httptreemux.ContextParams(r.Context())
	log.Printf("POST music: %s", params["name"])

	music := &music.Music{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(music)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = repository.Update(music)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(music)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type DeleteMusic struct{}

func (m *DeleteMusic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := httptreemux.ContextParams(r.Context())
	log.Printf("DELETE music: %s", params["name"])

	err := repository.Delete(params["name"])
	if err != nil {
		switch err {
		case db.ErrMusicNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type ListMusic struct{}

func (m *ListMusic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET musics")

	musics := repository.List()
	encoder := json.NewEncoder(w)
	encoder.Encode(musics)
}

func RunServer() {

	repository = db.New()

	addr := "127.0.0.1:8081"
	router := httptreemux.NewContextMux()

	router.Handler(http.MethodGet, "/music/:name", &GetMusic{})
	router.Handler(http.MethodPut, "/music/:name", &CreateMusic{})
	router.Handler(http.MethodPost, "/music/:name", &UpdateMusic{})
	router.Handler(http.MethodDelete, "/music/:name", &DeleteMusic{})
	router.Handler(http.MethodGet, "/music/", &ListMusic{})

	log.Printf("Running web server on: http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
