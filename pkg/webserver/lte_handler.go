package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type lteModuleHandler struct {
	c *APIComponents
}

func NewLTEModuleHandler(c *APIComponents) *lteModuleHandler {
	return &lteModuleHandler{
		c: c,
	}
}

func (h *lteModuleHandler) GetAccountBalance(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error   string  `json:"error"`
		Balance float64 `json:"balance"`
	}
	response := &Response{}

	b, err := h.c.MC.GetCellAccBalance()

	if err != nil {
		fmt.Println("Failed to get cell account balance: ", err)
	}

	response.Balance = b
	response.Error = err.Error()

	json.NewEncoder(rw).Encode(response)
}

func (h *lteModuleHandler) GetInetBalance(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error   string  `json:"error"`
		Balance float64 `json:"balance"`
	}
	response := &Response{}

	b, err := h.c.MC.GetCellInetBalance()

	if err != nil {
		fmt.Println("Failed to get cell internet balance: ", err)
	}

	response.Balance = b
	response.Error = err.Error()

	json.NewEncoder(rw).Encode(response)
}

func (h *lteModuleHandler) GetTariff(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Error  string `json:"error"`
		Tariff string `json:"tariff"`
	}
	response := &Response{}

	t, err := h.c.MC.GetCellTariff()

	if err != nil {
		fmt.Println("Failed to get cell internet balance: ", err)
	}

	response.Tariff = t
	response.Error = err.Error()

	json.NewEncoder(rw).Encode(response)
}
