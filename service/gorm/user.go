package gorm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	logger "github.com/anhdt-vnpay/f5_fulltext_search/lib/log"
	"github.com/anhdt-vnpay/f5_fulltext_search/model"
	"github.com/anhdt-vnpay/f5_fulltext_search/runtime"
)

var (
	userLogger = logger.NewLogger("api.service.user")
)

type Request struct {
	Table string     `json:"table"`
	User  model.User `json:"user"`
}

type Response struct {
	Code    uint64       `json:"code"`
	Message string       `json:"message"`
	Items   []model.User `json:"items"`
}

type UserManager struct {
	dbfs runtime.DbFullTextSearch
}

func NewUserManager(dbfs runtime.DbFullTextSearch) *UserManager {
	return &UserManager{
		dbfs: dbfs,
	}
}

func (re *UserManager) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not support", http.StatusBadRequest)
		return
	}

	tableName := r.URL.Query().Get("table")
	condition := r.URL.Query().Get("condition")

	if len(tableName) == 0 {
		userLogger.Errorf("Missing table name")
		http.Error(w, "missing table name", http.StatusBadRequest)
		return
	}

	userLogger.Infof("[GET] table name: %s, condition: %s", tableName, condition)

	var users []model.User
	err := re.dbfs.Get(tableName, condition, &users)
	if err != nil {
		userLogger.Errorf("[GET] error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// users = append(users, user)
	rs := Response{
		Code:    200,
		Message: "Successfully",
		Items:   users,
	}

	brs, err := json.Marshal(rs)
	if err != nil {
		userLogger.Errorf("[GET] error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(brs)
}

func (re *UserManager) Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not support", http.StatusBadRequest)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	bodyObj := Request{}
	if err = json.Unmarshal(reqBody, &bodyObj); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userLogger.Infof("[Insert] request: %s", bodyObj)

	user := bodyObj.User
	if err = re.dbfs.Insert(bodyObj.Table, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []model.User
	users = append(users, user)
	rs := Response{
		Code:    200,
		Message: "Successfully",
		Items:   users,
	}

	brs, err := json.Marshal(rs)
	if err != nil {
		userLogger.Errorf("[Insert] error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(brs)
}

func (re *UserManager) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not support", http.StatusBadRequest)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	bodyObj := Request{}
	if err = json.Unmarshal(reqBody, &bodyObj); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userLogger.Infof("[Update] request: %s", bodyObj)

	user := bodyObj.User
	if err = re.dbfs.Update(bodyObj.Table, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []model.User
	users = append(users, user)
	rs := Response{
		Code:    200,
		Message: "Successfully",
		Items:   users,
	}

	brs, err := json.Marshal(rs)
	if err != nil {
		userLogger.Errorf("[Update] error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(brs)
}

func (re *UserManager) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not support", http.StatusBadRequest)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	bodyObj := Request{}
	if err = json.Unmarshal(reqBody, &bodyObj); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userLogger.Infof("[Delete] request: %s", bodyObj)

	user := bodyObj.User
	if err = re.dbfs.Delete(bodyObj.Table, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []model.User
	users = append(users, user)
	rs := Response{
		Code:    200,
		Message: "Successfully",
		Items:   users,
	}

	brs, err := json.Marshal(rs)
	if err != nil {
		userLogger.Errorf("[Delete] error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(brs)
}

func (re *UserManager) SearchLite(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not support", http.StatusBadRequest)
		return
	}

	query := r.URL.Query().Get("query")

	userLogger.Infof("[SearchLite] query: %s", query)

	searchRs, err := re.dbfs.SearchLite(query)
	if err != nil {
		userLogger.Errorf("[SearchLite] error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userLogger.Infof("DEBUG search RS: %v >>>>> %T", searchRs, searchRs) // []map[string]interface{}

	var users []model.User

	srs := searchRs.([]map[string]interface{})

	for _, s := range srs {
		b, _ := json.Marshal(s)
		var user model.User
		err := json.Unmarshal(b, &user)
		if err != nil {
			userLogger.Errorf("[SearchLite] error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	rs := Response{
		Code:    200,
		Message: "Successfully",
		Items:   users,
	}

	brs, err := json.Marshal(rs)
	if err != nil {
		userLogger.Errorf("[GET] error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(brs)
}
