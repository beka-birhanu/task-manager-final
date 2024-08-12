// Package promotcmd provides the logic for promoting a user to admin status.
// It includes the necessary command structure and a handler to execute the promotion.
package promotcmd

import (
	"log"

	icmd "github.com/beka-birhanu/app/common/cqrs/command"
	irepo "github.com/beka-birhanu/app/common/i_repo"
)

// Handler handles the promote command logic.
type Handler struct {
	userRepo irepo.User
}

// Ensure Handler implement icmd.Handler
var _ icmd.IHandler[*Command, bool] = &Handler{}

// New creates a new instance of the Handler with the provided user repository.
func New(userRepo irepo.User) *Handler {
	return &Handler{userRepo: userRepo}
}

// Promote promotes a user to admin status based on the provided command.
func (h *Handler) Handle(cmd *Command) (bool, error) {
	user, err := h.userRepo.ByUsername(cmd.Username)
	if err != nil {
		return false, err
	}

	admin, err := h.userRepo.ById(cmd.PromoterID)
	if err != nil {
		return false, err
	}

	user.UpdateAdminStatus(true)
	if err := h.userRepo.Save(user); err != nil {
		return false, err
	}

	// TODO: Implement a proper logging mechanism.
	log.Printf("Admin %v promoted user %v", admin.Username(), user.Username())
	return true, nil
}
