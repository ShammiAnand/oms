package kitchen

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/shammianand/oms/services/common/auth"
	"github.com/shammianand/oms/services/common/config"
	"github.com/shammianand/oms/services/common/genproto/orders"
	"github.com/shammianand/oms/services/common/util"
	"google.golang.org/grpc"
)

type Handler struct {
	conn  *grpc.ClientConn
	store auth.UserStore
}

func NewHandler(conn *grpc.ClientConn, store auth.UserStore) *Handler {
	return &Handler{conn: conn, store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /login", h.handleLogin)
	router.HandleFunc("POST /signup", h.handleSignup)
	router.HandleFunc("GET /orders", auth.WithJWTAuth(h.getOrdersByCustomerID, h.store))
	router.HandleFunc("POST /orders", auth.WithJWTAuth(h.createOrder, h.store))
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

	payload := LoginUserPayload{}
	if err := util.ParseJSON(r, &payload); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		util.WriteError(
			w,
			http.StatusBadRequest,
			errors.New("invalid email or password"),
		)
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		util.WriteError(
			w,
			http.StatusBadRequest,
			errors.New("incorrect password"),
		)
		return
	}

	token, err := auth.CreateJWT([]byte(config.Envs.Secret), u.ID)
	if err != nil {
		util.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleSignup(w http.ResponseWriter, r *http.Request) {

	payload := RegisterUserPayload{}
	if err := util.ParseJSON(r, &payload); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		util.WriteError(
			w,
			http.StatusBadRequest,
			errors.New("user with the given email already exists"),
		)
		return
	}

	hashedPassword, err := auth.HashPasswords(payload.Password)
	if err != nil {
		util.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	newUser := &auth.User{
		FirstName: payload.FirstName,
		Email:     payload.Email,
		LastName:  payload.LastName,
		Password:  hashedPassword,
	}
	h.store.CreateUser(newUser)

	util.WriteJSON(
		w,
		http.StatusCreated,
		map[string]any{
			"id":    newUser.ID,
			"email": newUser.Email,
		},
	)

}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	c := orders.NewOrderServiceClient(h.conn)
	var createOrderReq orders.CreateOrderRequest
	err := util.ParseJSON(r, &createOrderReq)
	if err != nil {
		util.WriteError(
			w, http.StatusFailedDependency, err,
		)
	}
	createOrderReq.CustomerID = int32(auth.GetUserIdFromContext(r.Context()))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	createOrderRes, err := c.CreateOrder(ctx, &createOrderReq)
	if err != nil {
		util.WriteError(
			w, http.StatusFailedDependency, err,
		)
	}
	util.WriteJSON(w, http.StatusOK, createOrderRes)
}

func (h *Handler) getOrdersByCustomerID(w http.ResponseWriter, r *http.Request) {
	custID := auth.GetUserIdFromContext(r.Context())
	c := orders.NewOrderServiceClient(h.conn)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	res, err := c.GetOrders(ctx, &orders.GetOrdersRequest{
		CustomerID: int32(custID),
	})
	if err != nil {
		log.Fatalf("client error: %v", err)
	}
	util.WriteJSON(w, http.StatusOK, res)

}
