package handlers

import (
	"net/http"

	"github.com/diaku/eht-demo/internal/dto"
	"github.com/diaku/eht-demo/internal/usecase"
	"github.com/diaku/eht-demo/pkg/response"
)

type EthHandler struct {
	uc usecase.EthUC
}

func NewEthHandler(uc usecase.EthUC) *EthHandler {
	return &EthHandler{uc: uc}
}

func (eh EthHandler) Chain(w http.ResponseWriter, r *http.Request) {
	chainID, err := eh.uc.ChainInfo(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"chainId": chainID})
}

func (h *EthHandler) Balance(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("address")
	if addr == "" {
		response.Error(w, http.StatusBadRequest, "address query param is required")
		return
	}

	wei, eth, err := h.uc.Balance(r.Context(), addr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, dto.BalanceRes{
		Address: addr,
		Wei:     wei,
		Eth:     eth,
	})
}
