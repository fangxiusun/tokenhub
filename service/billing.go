package service

import (
	"log"

	"github.com/fangxiusun/tokenhub/model"
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
func PreConsume(userId, tokenId int, modelName string, estimatedQuota int) (*BillingSession, error) {
	user, err := model.GetUserById(userId, true)
	if err != nil {
		return nil, err
	}

	if user.Quota < estimatedQuota {
		return nil, ErrInsufficientQuota
	}

	user.Quota -= estimatedQuota
	if err := user.Update(false); err != nil {
		return nil, err
	}

	log.Printf("Pre-consumed %d quota from user %d for model %s", estimatedQuota, userId, modelName)

	return &BillingSession{
		UserId:  userId,
		TokenId: tokenId,
		Model:   modelName,
		Quota:   estimatedQuota,
		PrePaid: true,
	}, nil
}

// Settle settles the billing session with actual usage
func Settle(session *BillingSession, actualQuota int) error {
	user, err := model.GetUserById(session.UserId, true)
	if err != nil {
		return err
	}

	diff := actualQuota - session.Quota

	if diff > 0 {
		if user.Quota < diff {
			user.Quota += session.Quota
			log.Printf("Refunding %d quota to user %d due to insufficient balance", session.Quota, session.UserId)
		} else {
			user.Quota -= diff
			log.Printf("Charging additional %d quota from user %d", diff, session.UserId)
		}
	} else if diff < 0 {
		user.Quota += -diff
		log.Printf("Refunding %d quota to user %d", -diff, session.UserId)
	}

	if err := user.Update(false); err != nil {
		return err
	}

	return nil
}

// Refund refunds the entire pre-consumed amount
func Refund(session *BillingSession) error {
	if !session.PrePaid {
		return nil
	}

	user, err := model.GetUserById(session.UserId, true)
	if err != nil {
		return err
	}

	user.Quota += session.Quota
	log.Printf("Refunding %d quota to user %d", session.Quota, session.UserId)

	return user.Update(false)
}

