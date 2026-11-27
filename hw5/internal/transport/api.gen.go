package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	StatusCancelled  Status = "cancelled"
	StatusCompleted  Status = "completed"
	StatusNew        Status = "new"
	StatusProcessing Status = "processing"
)

func (e Status) Valid() bool {
	switch e {
	case StatusCancelled:
		return true
	case StatusCompleted:
		return true
	case StatusNew:
		return true
	case StatusProcessing:
		return true
	default:
		return false
	}
}

type CreateOrder struct {
	Customer *string `json:"customer,omitempty"`
	Items    []Item  `json:"items"`
}

type Error struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

type Item struct {
	Name string `json:"name"`
	Price    int64 `json:"price"`
	Quantity int64 `json:"quantity"`
}

type Order struct {
	CreatedAt  time.Time      `json:"created_at"`
	Customer   *string        `json:"customer,omitempty"`
	History    []StatusChange `json:"history"`
	Id         int            `json:"id"`
	Items      []Item         `json:"items"`
	Status     Status         `json:"status"`
	TotalPrice int64          `json:"total_price"`
}

type Status string

type StatusChange struct {
	ChangedAt time.Time `json:"changed_at"`
	Status    Status    `json:"status"`
}

// OrderID defines model for OrderID.
type OrderID = int

// BadRequest defines model for BadRequest.
type BadRequest = Error

// Conflict defines model for Conflict.
type Conflict = Error

// InternalError defines model for InternalError.
type InternalError = Error

// NotFound defines model for NotFound.
type NotFound = Error

// UpdateOrderStatusJSONBody defines parameters for UpdateOrderStatus.
type UpdateOrderStatusJSONBody struct {
	Status Status `json:"status"`
}

// CreateOrderJSONRequestBody defines body for CreateOrder for application/json ContentType.
type CreateOrderJSONRequestBody = CreateOrder

// UpdateOrderStatusJSONRequestBody defines body for UpdateOrderStatus for application/json ContentType.
type UpdateOrderStatusJSONRequestBody UpdateOrderStatusJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List orders
	// (GET /api/v1/orders)
	ListOrders(w http.ResponseWriter, r *http.Request)
	// Create an order
	// (POST /api/v1/orders)
	CreateOrder(w http.ResponseWriter, r *http.Request)
	// Get an order
	// (GET /api/v1/orders/{id})
	GetOrder(w http.ResponseWriter, r *http.Request, id OrderID)
	// Get chronological status history
	// (GET /api/v1/orders/{id}/history)
	GetOrderHistory(w http.ResponseWriter, r *http.Request, id OrderID)
	// Change order status
	// (PATCH /api/v1/orders/{id}/status)
	UpdateOrderStatus(w http.ResponseWriter, r *http.Request, id OrderID)
	// Check service availability
	// (GET /health)
	Health(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// List orders
// (GET /api/v1/orders)
func (_ Unimplemented) ListOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create an order
// (POST /api/v1/orders)
func (_ Unimplemented) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get an order
// (GET /api/v1/orders/{id})
func (_ Unimplemented) GetOrder(w http.ResponseWriter, r *http.Request, id OrderID) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get chronological status history
// (GET /api/v1/orders/{id}/history)
func (_ Unimplemented) GetOrderHistory(w http.ResponseWriter, r *http.Request, id OrderID) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Change order status
// (PATCH /api/v1/orders/{id}/status)
func (_ Unimplemented) UpdateOrderStatus(w http.ResponseWriter, r *http.Request, id OrderID) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Check service availability
// (GET /health)
func (_ Unimplemented) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// ListOrders operation middleware
func (siw *ServerInterfaceWrapper) ListOrders(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListOrders(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateOrder operation middleware
func (siw *ServerInterfaceWrapper) CreateOrder(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateOrder(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetOrder operation middleware
func (siw *ServerInterfaceWrapper) GetOrder(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id OrderID

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "integer", Format: ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOrder(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetOrderHistory operation middleware
func (siw *ServerInterfaceWrapper) GetOrderHistory(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id OrderID

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "integer", Format: ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOrderHistory(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// UpdateOrderStatus operation middleware
func (siw *ServerInterfaceWrapper) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id OrderID

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "integer", Format: ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateOrderStatus(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// Health operation middleware
func (siw *ServerInterfaceWrapper) Health(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Health(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/v1/orders", wrapper.ListOrders)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/v1/orders", wrapper.CreateOrder)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/v1/orders/{id}", wrapper.GetOrder)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/v1/orders/{id}/history", wrapper.GetOrderHistory)
	})
	r.Group(func(r chi.Router) {
		r.Patch(options.BaseURL+"/api/v1/orders/{id}/status", wrapper.UpdateOrderStatus)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/health", wrapper.Health)
	})

	return r
}

type BadRequestJSONResponse Error

type ConflictJSONResponse Error

type InternalErrorJSONResponse Error

type NotFoundJSONResponse Error

type ListOrdersRequestObject struct {
}

type ListOrdersResponseObject interface {
	VisitListOrdersResponse(w http.ResponseWriter) error
}

type ListOrders200JSONResponse []Order

func (response ListOrders200JSONResponse) VisitListOrdersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ListOrders500JSONResponse struct{ InternalErrorJSONResponse }

func (response ListOrders500JSONResponse) VisitListOrdersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type CreateOrderRequestObject struct {
	Body *CreateOrderJSONRequestBody
}

type CreateOrderResponseObject interface {
	VisitCreateOrderResponse(w http.ResponseWriter) error
}

type CreateOrder201JSONResponse Order

func (response CreateOrder201JSONResponse) VisitCreateOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type CreateOrder400JSONResponse struct{ BadRequestJSONResponse }

func (response CreateOrder400JSONResponse) VisitCreateOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type CreateOrder500JSONResponse struct{ InternalErrorJSONResponse }

func (response CreateOrder500JSONResponse) VisitCreateOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetOrderRequestObject struct {
	Id OrderID `json:"id"`
}

type GetOrderResponseObject interface {
	VisitGetOrderResponse(w http.ResponseWriter) error
}

type GetOrder200JSONResponse Order

func (response GetOrder200JSONResponse) VisitGetOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetOrder400JSONResponse struct{ BadRequestJSONResponse }

func (response GetOrder400JSONResponse) VisitGetOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetOrder404JSONResponse struct{ NotFoundJSONResponse }

func (response GetOrder404JSONResponse) VisitGetOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetOrder500JSONResponse struct{ InternalErrorJSONResponse }

func (response GetOrder500JSONResponse) VisitGetOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetOrderHistoryRequestObject struct {
	Id OrderID `json:"id"`
}

type GetOrderHistoryResponseObject interface {
	VisitGetOrderHistoryResponse(w http.ResponseWriter) error
}

type GetOrderHistory200JSONResponse []StatusChange

func (response GetOrderHistory200JSONResponse) VisitGetOrderHistoryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetOrderHistory400JSONResponse struct{ BadRequestJSONResponse }

func (response GetOrderHistory400JSONResponse) VisitGetOrderHistoryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetOrderHistory404JSONResponse struct{ NotFoundJSONResponse }

func (response GetOrderHistory404JSONResponse) VisitGetOrderHistoryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetOrderHistory500JSONResponse struct{ InternalErrorJSONResponse }

func (response GetOrderHistory500JSONResponse) VisitGetOrderHistoryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type UpdateOrderStatusRequestObject struct {
	Id   OrderID `json:"id"`
	Body *UpdateOrderStatusJSONRequestBody
}

type UpdateOrderStatusResponseObject interface {
	VisitUpdateOrderStatusResponse(w http.ResponseWriter) error
}

type UpdateOrderStatus200JSONResponse Order

func (response UpdateOrderStatus200JSONResponse) VisitUpdateOrderStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type UpdateOrderStatus400JSONResponse struct{ BadRequestJSONResponse }

func (response UpdateOrderStatus400JSONResponse) VisitUpdateOrderStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type UpdateOrderStatus404JSONResponse struct{ NotFoundJSONResponse }

func (response UpdateOrderStatus404JSONResponse) VisitUpdateOrderStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type UpdateOrderStatus409JSONResponse struct{ ConflictJSONResponse }

func (response UpdateOrderStatus409JSONResponse) VisitUpdateOrderStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type UpdateOrderStatus500JSONResponse struct{ InternalErrorJSONResponse }

func (response UpdateOrderStatus500JSONResponse) VisitUpdateOrderStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type HealthRequestObject struct {
}

type HealthResponseObject interface {
	VisitHealthResponse(w http.ResponseWriter) error
}

type Health200JSONResponse struct {
	Status *string `json:"status,omitempty"`
}

func (response Health200JSONResponse) VisitHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// List orders
	// (GET /api/v1/orders)
	ListOrders(ctx context.Context, request ListOrdersRequestObject) (ListOrdersResponseObject, error)
	// Create an order
	// (POST /api/v1/orders)
	CreateOrder(ctx context.Context, request CreateOrderRequestObject) (CreateOrderResponseObject, error)
	// Get an order
	// (GET /api/v1/orders/{id})
	GetOrder(ctx context.Context, request GetOrderRequestObject) (GetOrderResponseObject, error)
	// Get chronological status history
	// (GET /api/v1/orders/{id}/history)
	GetOrderHistory(ctx context.Context, request GetOrderHistoryRequestObject) (GetOrderHistoryResponseObject, error)
	// Change order status
	// (PATCH /api/v1/orders/{id}/status)
	UpdateOrderStatus(ctx context.Context, request UpdateOrderStatusRequestObject) (UpdateOrderStatusResponseObject, error)
	// Check service availability
	// (GET /health)
	Health(ctx context.Context, request HealthRequestObject) (HealthResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// ListOrders operation middleware
func (sh *strictHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	var request ListOrdersRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ListOrders(ctx, request.(ListOrdersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListOrders")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListOrdersResponseObject); ok {
		if err := validResponse.VisitListOrdersResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// CreateOrder operation middleware
func (sh *strictHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request CreateOrderRequestObject

	var body CreateOrderJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateOrder(ctx, request.(CreateOrderRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateOrder")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateOrderResponseObject); ok {
		if err := validResponse.VisitCreateOrderResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetOrder operation middleware
func (sh *strictHandler) GetOrder(w http.ResponseWriter, r *http.Request, id OrderID) {
	var request GetOrderRequestObject

	request.Id = id

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetOrder(ctx, request.(GetOrderRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetOrder")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetOrderResponseObject); ok {
		if err := validResponse.VisitGetOrderResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetOrderHistory operation middleware
func (sh *strictHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request, id OrderID) {
	var request GetOrderHistoryRequestObject

	request.Id = id

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetOrderHistory(ctx, request.(GetOrderHistoryRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetOrderHistory")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetOrderHistoryResponseObject); ok {
		if err := validResponse.VisitGetOrderHistoryResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// UpdateOrderStatus operation middleware
func (sh *strictHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request, id OrderID) {
	var request UpdateOrderStatusRequestObject

	request.Id = id

	var body UpdateOrderStatusJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.UpdateOrderStatus(ctx, request.(UpdateOrderStatusRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "UpdateOrderStatus")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(UpdateOrderStatusResponseObject); ok {
		if err := validResponse.VisitUpdateOrderStatusResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Health operation middleware
func (sh *strictHandler) Health(w http.ResponseWriter, r *http.Request) {
	var request HealthRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Health(ctx, request.(HealthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Health")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(HealthResponseObject); ok {
		if err := validResponse.VisitHealthResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
