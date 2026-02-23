package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/diaku/eht-demo/internal/dto"
	"github.com/diaku/eht-demo/internal/usecase"
	"github.com/diaku/eht-demo/pkg/response"
)

type AuthHanlder struct {
	uc usecase.AuthUC
}

func NewAuthHandler(uc usecase.AuthUC) *AuthHanlder {
	return &AuthHanlder{uc: uc}
}

func (ah AuthHanlder) Nouce(w http.ResponseWriter, r *http.Request) {
	nonce, err := ah.uc.NewNonce(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"nonce": nonce})
}

func (ah AuthHanlder) VerifySIWE(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}
	if req.Message == "" || req.Signature == "" {
		response.Error(w, http.StatusBadRequest, "message and signature are required")
		return
	}

	token, addr, err := ah.uc.VerifySIWE(r.Context(), req.Message, req.Signature)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, dto.VerifyRes{
		Token:   token,
		Address: addr,
	})
}
