// Package trangrp maintains the group of handlers for transaction example.
package trangrp

import (
	"context"
	"net/http"

	"github.com/ahrav/service/business/core/product"
	"github.com/ahrav/service/business/core/user"
	v1 "github.com/ahrav/service/business/web/v1"
	"github.com/ahrav/service/foundation/web"
)

// Handlers manages the set of product endpoints.
type Handlers struct {
	user    *user.Core
	product *product.Core
}

// New constructs a handlers for route access.
func New(user *user.Core, product *product.Core) *Handlers {
	return &Handlers{
		user:    user,
		product: product,
	}
}

// Create adds a new user and product at the same time under a single transaction.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	h, err := h.executeUnderTransaction(ctx)
	if err != nil {
		return err
	}

	var app AppNewTran
	if err := web.Decode(r, &app); err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	np, err := toCoreNewProduct(app.Product)
	if err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	nu, err := toCoreNewUser(app.User)
	if err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nu)
	if err != nil {
		return err
	}

	np.UserID = usr.ID

	prd, err := h.product.Create(ctx, np)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, toAppProduct(prd), http.StatusCreated)
}
