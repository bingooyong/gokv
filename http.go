package main

import "net/http"

func getHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	val, err := srv.Get(key)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(val.(string)))
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	val := r.Form.Get("val")

	if err := srv.Set(key, val); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func delHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	if err := srv.Delete(key); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
