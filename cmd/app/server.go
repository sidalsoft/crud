package app

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sidalsoft/crud/pkg/customers"
	"github.com/sidalsoft/crud/pkg/security"
	"log"
	"net/http"
	"strconv"
	"time"
)

//Server представляет собой логический сервер нашего приложения
type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
	authSvc     *security.AuthService
}

func NewServer(mux *mux.Router, customerSvc *customers.Service, authSvc *security.AuthService) *Server {
	return &Server{mux: mux, customerSvc: customerSvc, authSvc: authSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

//INIt инициализирует сервер (регистрирует все Handlerы)
func (s *Server) Init() {

	s.mux.HandleFunc("/api/customers", s.handleSave).Methods(POST)
	s.mux.HandleFunc("/api/customers/token", s.handleGenerateToken).Methods(POST)
	s.mux.HandleFunc("/api/customers/token/validate", s.handleValidateToken).Methods(POST)

	//s.mux.Use(middleware.Basic(s.authSvc))
	//s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	//s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	//s.mux.HandleFunc("/customers", s.handleGetAllActiveCustomers)
	//s.mux.HandleFunc("/customers.blockById", s.handleBlockByID)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)
	//s.mux.HandleFunc("/customers.unblockById", s.handleUnBlockByID)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockByID).Methods(DELETE)
	//s.mux.HandleFunc("/customers.removeById", s.handleDelete)
	s.mux.HandleFunc("/customers/{id}", s.handleDelete).Methods(DELETE)
	//s.mux.HandleFunc("/customers.save", s.handleSave)
	s.mux.HandleFunc("/customers", s.handleSave).Methods(POST)
}

func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}

	item, err := s.customerSvc.ByID(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetAllCustomers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customerSvc.All(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, items)
}

func (s *Server) handleGetAllActiveCustomers(writer http.ResponseWriter, request *http.Request) {

	items, err := s.customerSvc.AllActive(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	parceJSON(writer, items)
}

func (s *Server) handleBlockByID(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customerSvc.ChangeActive(request.Context(), id, false)

	if errors.Is(err, customers.ErrNotFound) {
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, item)
}

func (s *Server) handleUnBlockByID(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customerSvc.ChangeActive(request.Context(), id, true)

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, item)
}

func (s *Server) handleDelete(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customerSvc.Delete(request.Context(), id)

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, item)
}

func (s *Server) handleSave(writer http.ResponseWriter, request *http.Request) {
	var item *customers.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item.Active = true
	item.Created = time.Now()
	customer, err := s.customerSvc.Save(request.Context(), item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, customer)
}

func (s *Server) handleGenerateToken(writer http.ResponseWriter, request *http.Request) {
	data := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	token, err := s.authSvc.TokenForCustomer(request.Context(), data.Login, data.Password)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, struct {
		Token string `json:"token"`
	}{Token: token})
}

func (s *Server) handleValidateToken(writer http.ResponseWriter, request *http.Request) {
	data := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, expire, err := s.authSvc.AuthenticateCustomer(request.Context(), data.Token)
	if err == security.ErrNoSuchUser {

		parceErrJSON(writer, struct {
			Status string `json:"status"`
			Reason string `json:"reason"`
		}{Status: "fail", Reason: "not found"}, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if expire.Unix() < time.Now().Unix() {
		parceErrJSON(writer, struct {
			Status string `json:"status"`
			Reason string `json:"reason"`
		}{Status: "fail", Reason: "expired"}, http.StatusBadRequest)
		return
	}

	parceJSON(writer, struct {
		Status     string `json:"status"`
		CustomerId int64  `json:"customerId"`
	}{Status: "ok", CustomerId: id})
}

func parceJSON(writer http.ResponseWriter, iData interface{}) {

	data, err := json.Marshal(iData)

	if err != nil {
		log.Println(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func parceErrJSON(writer http.ResponseWriter, iData interface{}, code int) {

	data, err := json.Marshal(iData)

	if err != nil {
		log.Println(writer, http.StatusInternalServerError, err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}
