package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"
	"gorm.io/gorm"
)

type PaymentService interface {
	InitiatePayment(billID, userID uint) (string, error)
	GetPaymentHistory(userID uint) ([]model.Pembayaran, error)
	ProcessMidtransNotification(notificationPayload map[string]any) error
}

type paymentService struct {
	billRepo    repository.BillRepository
	studentRepo repository.StudentRepository
	paymentRepo repository.PaymentRepository
	midtransSvc MidtransService
	db          *gorm.DB
}

func NewPaymentService(billRepo repository.BillRepository, studentRepo repository.StudentRepository, paymentRepo repository.PaymentRepository, midtransSvc MidtransService, db *gorm.DB) PaymentService {
	return &paymentService{billRepo, studentRepo, paymentRepo, midtransSvc, db}
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

func (s *paymentService) ProcessMidtransNotification(notificationPayload map[string]interface{}) error {
	orderID, exists := notificationPayload["order_id"].(string)
	if !exists {
		return errors.New("invalid notification payload: missing order_id")
	}

	transactionStatus, _ := notificationPayload["transaction_status"].(string)
	transactionID, _ := notificationPayload["transaction_id"].(string)
	paymentType, _ := notificationPayload["payment_type"].(string)
	settlementTime, _ := notificationPayload["settlement_time"].(string)

	responseBytes, err := json.Marshal(notificationPayload)
	if err != nil {
		return err
	}
	midtransResponse := string(responseBytes)

	err = s.db.Exec("CALL UpdateStatusPembayaran(?, ?, ?, ?, ?, ?)",
		orderID,
		transactionStatus,
		transactionID,
		paymentType,
		settlementTime,
		midtransResponse,
	).Error

	if err != nil {
		return fmt.Errorf("gagal menjalankan stored procedure: %w", err)
	}

	return nil
}
