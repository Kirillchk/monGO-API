package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func InitHandlers() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./dist"))))

	http.Handle("/user/reg", http.HandlerFunc(handleRegUser))
	http.Handle("/user/log", http.HandlerFunc(handleLogUser))
	http.Handle("/user/jwt", http.HandlerFunc(handleUserJWT))

	http.Handle("/DB/collections", http.HandlerFunc(handleListCollections))
	http.Handle("/DB/collection", http.HandlerFunc(handleCollection))
	http.Handle("/DB/document", http.HandlerFunc(handleDocument))
}

type UserBody struct {
	Login    string
	Password string
}
type UserBodyWJWT struct {
	Login string
	JWT   string
}

// User handlers
func handleRegUser(res http.ResponseWriter, req *http.Request) {
	var данные UserBody
	ошибка := json.NewDecoder(req.Body).Decode(&данные)
	if ошибка != nil {
		return
	}
	длинаЛогина := len(данные.Login)
	длинаПароля := len(данные.Password)
	if длинаЛогина <= 5 {
		http.Error(res, "БРАТЕЦ!!! Твой логин слишком мал!", http.StatusBadRequest)
	} else if длинаПароля <= 5 {
		http.Error(res, "БРАТЕЦ!!! Твой пароль слишком мал!", http.StatusBadRequest)
	}
	jwt, err := AddUser(данные.Login, данные.Password)
	if err != nil {
		http.Error(res, "User already exists", http.StatusBadRequest)
	}
	res.Write([]byte(jwt))
}
func handleLogUser(res http.ResponseWriter, req *http.Request) {
	var данные UserBody
	ошибка := json.NewDecoder(req.Body).Decode(&данные)
	if ошибка != nil {
		http.Error(res, "Malformed data", http.StatusBadRequest)
		return
	}
	user, err := FindUser(данные.Login)
	if err != nil {
		http.Error(res, fmt.Sprintf("User not found:%v", err), http.StatusNotFound)
	}
	if user.Password != данные.Password {
		http.Error(res, "incorrect password", http.StatusUnauthorized)
	}
	res.Write([]byte(user.JWT))
}
func handleUserJWT(res http.ResponseWriter, req *http.Request) {
	var data UserBodyWJWT
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		http.Error(res, "Decoding error", http.StatusBadRequest)
		return
	}
	DBuser, err := FindUser(data.Login)
	if err != nil {
		http.Error(res, "User wasnt found", http.StatusNotFound)
		return
	}
	if DBuser.JWT != data.JWT {
		http.Error(res, "JWT is incorrect", http.StatusUnauthorized)
		return
	}
}

// Database handlers
func handleListCollections(res http.ResponseWriter, req *http.Request) {
	result, _ := ListCollections()
	json.NewEncoder(res).Encode(result)
}
func handleCollection(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	colectionName := query.Get("collection")
	switch req.Method {
	case http.MethodGet:
		bson, _ := FindCollection(colectionName)
		result, _ := json.Marshal(bson)
		res.Write([]byte(result))
	case http.MethodDelete:
		DeleteCollection(colectionName)
	case http.MethodPost:
		AddColletion(colectionName)
	}
}
func handleDocument(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	colectionName := query.Get("collection")
	var doc Document
	if req.Method == http.MethodPatch {
		var body bson.A
		json.NewDecoder(req.Body).Decode(&body)
		fmt.Println(body[0])
		var body1 bson.M = body[0].(map[string]interface{})
		var body2 bson.M = body[1].(map[string]interface{})
		doc = Document{
			collection: colectionName,
			document:   body1,
		}
		doc.Update(body2)
	} else {
		var body bson.M
		json.NewDecoder(req.Body).Decode(&body)
		doc = Document{
			collection: colectionName,
			document:   body,
		}
		switch req.Method {
		case http.MethodPost:
			doc.Add()
		case http.MethodDelete:
			doc.Delete()
		}
	}

}
