package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/streatcodes/bernard/bernard-server/db"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) ListenHTTP() error {
	fmt.Printf("Listening on %s\n", s.Config.HTTPListenAddr)
	http.HandleFunc("/", s.fileHandler())
	http.HandleFunc("/auth", s.authHandler())
	http.HandleFunc("/connect", s.websocketHandler())

	return http.ListenAndServe(s.Config.HTTPListenAddr, nil)
}

//TODO should return index.html/200 on 404 for client routing
func (s *Server) fileHandler() http.HandlerFunc {
	handler := http.FileServer(http.Dir("www"))

	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

//TODO add throttle
func (s *Server) authHandler() http.HandlerFunc {
	type authReq struct {
		Email    string
		Password string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//Read incoming req
		if r.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var authAttempt authReq
		err := json.NewDecoder(r.Body).Decode(&authAttempt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//Look up email in DB
		user, err := db.FindUserByEmail(s.DB, authAttempt.Email)
		if err != nil {
			fmt.Printf("Error looking up user in DB: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//Verify password
		err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(authAttempt.Password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//Generate key
		key, err := generateNewKey()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//Assign session to user
		err = db.AddUserSession(s.DB, user.ID, key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.WriteString(w, hex.EncodeToString(key))
	}
}

func (s *Server) websocketHandler() http.HandlerFunc {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//Check user is authenticated
		authHeader := r.Header.Get("Authorization")
		key, err := hex.DecodeString(authHeader)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = db.FindUserByKey(s.DB, key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//Upgrade to websocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		conn.Close()
	}
}
