package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func InitHandlers() {
	http.Handle("/DB/reguser", http.HandlerFunc(handleRegUser))
	http.Handle("/DB/loguser", http.HandlerFunc(handleLogUser))
	http.Handle("/DB/collections", http.HandlerFunc(handleListCollections))
	http.Handle("/DB/collection", http.HandlerFunc(handleCollection))
	http.Handle("/DB/document", http.HandlerFunc(handleDocument))

	http.Handle("/dashboard", http.FileServer(http.Dir("./dist")))

}

type UserBody struct {
	Login    string
	Password string
}

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
		return
	}
	jwt, err := AddUser(данные.Login, данные.Password)
	if err != nil {
		http.Error(res, "User already exists", http.StatusBadRequest)
	}
	res.Write([]byte(jwt))
}
func handleListCollections(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(ListCollections())
}
func handleCollection(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	colectionName := query.Get("collection")
	if req.Method == http.MethodDelete {
		DeleteCollection(colectionName)
	} else {
		bson, _ := FindCollection(colectionName)
		res.Write([]byte(fmt.Sprint(bson)))
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
