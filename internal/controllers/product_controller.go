package controllers

import (
	"encoding/json"
	"golang-sqlserver/internal/models"
	"golang-sqlserver/internal/services"
	"net/http"
	"strconv"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		sendResponse(w, http.StatusBadRequest, msgInvalidPayload, nil)
		return
	}

	if err := c.service.Create(r.Context(), &product); err != nil {
		sendResponse(w, http.StatusInternalServerError, msgFailedCreate, nil)
		return
	}

	sendResponse(w, http.StatusCreated, msgCreateSuccess, product)
}

func (c *ProductController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		sendResponse(w, http.StatusBadRequest, msgInvalidID, nil)
		return
	}

	product, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, msgFailedGet, nil)
		return
	}

	if product == nil {
		sendResponse(w, http.StatusNotFound, msgProductNotFound, nil)
		return
	}

	sendResponse(w, http.StatusOK, msgGetSuccess, product)
}

func (c *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := c.service.GetAll(r.Context())
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, msgFailedGetAll, nil)
		return
	}

	sendResponse(w, http.StatusOK, msgGetAllSuccess, products)
}

func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		sendResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if err := c.service.Update(r.Context(), &product); err != nil {
		sendResponse(w, http.StatusInternalServerError, msgFailedUpdate, nil)
		return
	}

	sendResponse(w, http.StatusOK, msgUpdateSuccess, product)
}

func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		sendResponse(w, http.StatusBadRequest, msgInvalidID, nil)
		return
	}

	if err := c.service.Delete(r.Context(), id); err != nil {
		sendResponse(w, http.StatusInternalServerError, msgFailedDelete, nil)
		return
	}

	sendResponse(w, http.StatusOK, msgDeleteSuccess, nil)
}

func (c *ProductController) Search(w http.ResponseWriter, r *http.Request) {
	var params services.SearchParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		sendResponse(w, http.StatusBadRequest, msgInvalidPayload, nil)
		return
	}

	products, err := c.service.Search(r.Context(), params)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, msgFailedSearchProduct, nil)
		return
	}

	if len(products) == 0 {
		sendResponse(w, http.StatusNotFound, msgProductNotFound, nil)
		return
	}

	sendResponse(w, http.StatusOK, msgProductFound, products)
}

func sendResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

const (
	msgInvalidPayload  = "Invalid request payload"
	msgInvalidID       = "Invalid ID parameter"
	msgProductNotFound = "Product not found"
	msgCreateSuccess   = "Product created successfully"
	msgUpdateSuccess   = "Product updated successfully"
	msgDeleteSuccess   = "Product deleted successfully"
	msgGetSuccess      = "Product retrieved successfully"
	msgGetAllSuccess   = "Products retrieved successfully"

	msgFailedCreate        = "Failed to create product"
	msgFailedUpdate        = "Failed to update product"
	msgFailedDelete        = "Failed to delete product"
	msgFailedGet           = "Failed to get product"
	msgFailedGetAll        = "Failed to get products"
	msgFailedSearchProduct = "Failed to search products"
	msgProductFound        = "Products found"
)

type SearchParams struct {
	Keyword         string `json:"keyword"`
	ProductCategory string `json:"product_category,omitempty"`
	Tier            string `json:"tier,omitempty"`
}
