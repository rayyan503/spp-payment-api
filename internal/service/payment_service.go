package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"
)

type PaymentService interface {
	InitiatePayment(billID, userID uint) (string, error)
	GetPaymentHistory(userID uint) ([]model.Pembayaran, error)
}

type paymentService struct {
	billRepo    repository.BillRepository
	studentRepo repository.StudentRepository
	paymentRepo repository.PaymentRepository
	midtransSvc MidtransService
}

func NewPaymentService(billRepo repository.BillRepository, studentRepo repository.StudentRepository, paymentRepo repository.PaymentRepository, midtransSvc MidtransService) PaymentService {
	return &paymentService{billRepo, studentRepo, paymentRepo, midtransSvc}
}

func (s *paymentService) InitiatePayment(billID, userID uint) (string, error) {
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return "", errors.New("tagihan tidak ditemukan")
	}

	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return "", errors.New("profil siswa tidak ditemukan")
	}

	if bill.SiswaID != student.ID {
		return "", errors.New("tagihan ini bukan milik Anda")
	}

	if bill.StatusPembayaran != "belum_bayar" {
		return "", errors.New("tagihan ini sudah dibayar atau sedang diproses")
	}

	orderID := fmt.Sprintf("SPP-%d-%d", bill.ID, time.Now().Unix())
	newPayment := &model.Pembayaran{
		TagihanID:        bill.ID,
		SiswaID:          student.ID,
		OrderID:          orderID,
		JumlahBayar:      bill.JumlahTagihan,
		StatusPembayaran: "pending",
	}

	if err := s.paymentRepo.Create(newPayment); err != nil {
		return "", err
	}

	snapToken, err := s.midtransSvc.CreateTransaction(orderID, int64(bill.JumlahTagihan))
	if err != nil {
		_ = s.paymentRepo.Delete(newPayment.ID)
		return "", err
	}

	return snapToken, nil
}

func (s *paymentService) GetPaymentHistory(userID uint) ([]model.Pembayaran, error) {
	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("profil siswa tidak ditemukan")
	}

	return s.paymentRepo.FindAllBySiswaID(student.ID)
}
