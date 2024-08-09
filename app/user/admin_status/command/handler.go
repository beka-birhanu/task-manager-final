// Package promotcmd provides the logic for promoting a user to admin status.
// It includes the necessary command structure and a handler to execute the promotion.
package promotcmd

import (
	"log"

	irepo "github.com/beka-birhanu/app/common/i_repo"
)

// Handler handles the promote command logic.
type Handler struct {
	userRepo irepo.User
}

// Promote promotes a user to admin status.
func (h *Handler) Promote(cmd *Command) error {
	user, err := h.userRepo.ByUsername(cmd.Username)
	if err != nil {
		return err
	}

	admin, err := h.userRepo.ById(cmd.PromoterID)
	if err != nil {
		return err
	}

	user.UpdateAdminStatus(true)
	if err := h.userRepo.Save(user); err != nil {
		return err
	}

	// TODO: Implement proper logging mechanism.
	log.Printf("Admin %v promoted user %v", admin.Username(), user.Username())
	return nil
}

