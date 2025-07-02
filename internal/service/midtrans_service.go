package service

import (
	"github.com/hiuncy/spp-payment-api/internal/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService interface {
	CreateTransaction(orderID string, amount int64) (string, error)
}

type midtransService struct {
	snapClient snap.Client
}

func NewMidtransService(cfg *config.Config) MidtransService {
	var client snap.Client
	env := midtrans.Sandbox
	if cfg.MidtransEnvironment == "production" {
		env = midtrans.Production
	}

	client.New(cfg.MidtransServerKey, env)
	return &midtransService{snapClient: client}
}

func (s *midtransService) CreateTransaction(orderID string, amount int64) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	token, err := s.snapClient.CreateTransactionToken(req)
	if err != nil {
		return "", err
	}

	return token, nil
}
