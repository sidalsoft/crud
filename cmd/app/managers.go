package app

import (
	"encoding/json"
	"github.com/sidalsoft/crud/cmd/app/middleware"
	"github.com/sidalsoft/crud/pkg/managers"
	"github.com/sidalsoft/crud/pkg/products"
	"github.com/sidalsoft/crud/pkg/salePositions"
	"github.com/sidalsoft/crud/pkg/sales"
	"net/http"
	"time"
)

func (s *Server) handleManagerRegistration(writer http.ResponseWriter, request *http.Request) {
	var data *managers.Managers
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		println(http.StatusBadRequest)
		println(err)
		return
	}
	data.Active = true
	data.Created = time.Now()
	manager, err := s.managerSvc.Save(request.Context(), data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		println(http.StatusInternalServerError)
		println(err)
		return
	}
	token, err := s.managerSvc.TokenForManager(request.Context(), manager.Phone, data.Password)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		println(http.StatusText(http.StatusInternalServerError), err.Error())
		return
	}
	parceJSON(writer, struct {
		Token string `json:"token"`
	}{Token: token})
	//parceJSON(writer, manager)
}

func (s *Server) handleManagerGetToken(writer http.ResponseWriter, request *http.Request) {
	data := struct {
		Login    string `json:"phone"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		println(http.StatusText(http.StatusBadRequest), err.Error())
		return
	}
	token, err := s.managerSvc.TokenForManager(request.Context(), data.Login, data.Password)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		println(http.StatusText(http.StatusInternalServerError), err.Error())
		return
	}
	parceJSON(writer, struct {
		Token string `json:"token"`
	}{Token: token})
}

func (s *Server) handleManagerGetSales(writer http.ResponseWriter, request *http.Request) {
	managerId, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		println(http.StatusText(http.StatusBadRequest), err.Error())
		return
	}
	total, err := s.saleSvc.TotalByManager(request.Context(), managerId)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		println(http.StatusText(http.StatusInternalServerError), err.Error())
		return
	}
	parceJSON(writer, struct {
		ManagerId int64 `json:"manager_id"`
		Total     int   `json:"total"`
	}{ManagerId: managerId, Total: total})
}

func (s *Server) handleManagerMakeSale(writer http.ResponseWriter, request *http.Request) {
	managerId, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		println(http.StatusText(http.StatusBadRequest), err.Error())
		return
	}
	data := struct {
		Id         int64  `json:"id"`
		CustomerId *int64 `json:"customer_id"`
		Positions  []struct {
			Id        int64  `json:"id"`
			ProductId int64  `json:"product_id"`
			Name      string `json:"name"`
			Qty       int    `json:"qty"`
			Price     int    `json:"price"`
		} `json:"positions"`
	}{}
	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		println(http.StatusText(http.StatusBadRequest), err.Error())
		return
	}
	for _, position := range data.Positions {
		product, err := s.productSvc.ByID(request.Context(), position.ProductId)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			println(http.StatusText(http.StatusInternalServerError), err.Error())
			return
		}
		if product.Qty < position.Qty {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	saleData := &sales.Sales{
		ID:         data.Id,
		ManagerId:  managerId,
		CustomerId: data.CustomerId,
	}
	sale, err := s.saleSvc.Save(request.Context(), saleData)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		println(http.StatusText(http.StatusInternalServerError), err.Error())
		return
	}
	for _, position := range data.Positions {
		salePositionData := &salePositions.SalePositions{
			ID:        position.Id,
			SaleId:    sale.ID,
			ProductId: position.ProductId,
			Name:      position.Name,
			Price:     position.Price,
			Qty:       position.Qty,
		}
		_, err := s.salePositionsSvc.Save(request.Context(), salePositionData)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			println(http.StatusText(http.StatusInternalServerError), err.Error())
			return
		}
	}

	parceJSON(writer, struct {
		Id int64 `json:"id"`
	}{Id: sale.ID})
}

func (s *Server) handleManagerGetProducts(writer http.ResponseWriter, request *http.Request) {

}

func (s *Server) handleManagerChangeProduct(writer http.ResponseWriter, request *http.Request) {
	var data *products.Product
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		println(err.Error())
		return
	}
	data.Active = true
	data.Created = time.Now()

	product, err := s.productSvc.Save(request.Context(), data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		println(err.Error())
		return
	}
	parceJSON(writer, product)
}

func (s *Server) handleManagerRemoveProductByID(writer http.ResponseWriter, request *http.Request) {

}

func (s *Server) handleManagerGetCustomers(writer http.ResponseWriter, request *http.Request) {

}

func (s *Server) handleManagerChangeCustomer(writer http.ResponseWriter, request *http.Request) {

}

func (s *Server) handleManagerRemoveCustomerByID(writer http.ResponseWriter, request *http.Request) {

}
