package service

import (
	"log"

	"github.com/your-username/your-project/model"
)

// BillingSession represents a billing session
type BillingSession struct {
	UserId    int
	TokenId   int
	Model     string
	Quota     int
	PrePaid   bool
}

// PreConsume pre-consumes quota for a request
func PreConsume(userId, tokenId int, model string, estimatedQuota int) (*BillingSession, error) {
	// Get user
	user, err := model.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	// Check if user has enough quota
	if user.Quota < estimatedQuota {
		return nil, ErrInsufficientQuota
	}

	// Deduct quota
	user.Quota -= estimatedQuota
	if err := model.UpdateUser(user); err != nil {
		return nil, err
	}

	log.Printf("Pre-consumed %d quota from user %d for model %s", estimatedQuota, userId, model)

	return &BillingSession{
		UserId:  userId,
		TokenId: tokenId,
		Model:   model,
		Quota:   estimatedQuota,
		PrePaid: true,
	}, nil
}

// Settle settles the billing session with actual usage
func Settle(session *BillingSession, actualQuota int) error {
	// Get user
	user, err := model.GetUserById(session.UserId)
	if err != nil {
		return err
	}

	// Calculate difference
	diff := actualQuota - session.Quota

	if diff > 0 {
		// Under-paid, need to charge more
		if user.Quota < diff {
			// Not enough quota, refund original
			user.Quota += session.Quota
			log.Printf("Refunding %d quota to user %d due to insufficient balance", session.Quota, session.UserId)
		} else {
			// Charge additional
			user.Quota -= diff
			log.Printf("Charging additional %d quota from user %d", diff, session.UserId)
		}
	} else if diff < 0 {
		// Over-paid, need to refund
		user.Quota += -diff
		log.Printf("Refunding %d quota to user %d", -diff, session.UserId)
	}

	// Update user
	if err := model.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

// Refund refunds the entire pre-consumed amount
func Refund(session *BillingSession) error {
	if !session.PrePaid {
		return nil
	}

	// Get user
	user, err := model.GetUserById(session.UserId)
	if err != nil {
		return err
	}

	// Refund quota
	user.Quota += session.Quota
	log.Printf("Refunding %d quota to user %d", session.Quota, session.UserId)

	// Update user
	return model.UpdateUser(user)
}
