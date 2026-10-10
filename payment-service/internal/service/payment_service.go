package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"

	"com.ecommerce/payment-service/internal/model"
	"com.ecommerce/payment-service/internal/repository"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrOrderAlreadyPaid = errors.New("Order has already been paid.")
)

const (
	TopicPaymentSuccessful = "SUCCESSFUL"
	TopicPaymentCompleted  = "paymentCompleted"
)

type PaymentService struct {
	repo     *repository.PaymentRepository
	producer *EventProducer
}

func NewPaymentService(repo *repository.PaymentRepository, producer *EventProducer) *PaymentService {
	return &PaymentService{
		repo:     repo,
		producer: producer,
	}
}

func (s *PaymentService) FindAll(ctx context.Context) ([]model.PaymentDto, error) {
	return s.repo.FindAll(ctx)
}

func (s *PaymentService) FindAllPaged(ctx context.Context, page, size int, sortBy, sortOrder string) (*model.PageResponse[model.PaymentDto], error) {
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}

	payments, total, err := s.repo.FindAllPaged(ctx, page, size, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))
	if totalPages == 0 && total == 0 {
		totalPages = 0
	}

	return &model.PageResponse[model.PaymentDto]{
		Content:          payments,
		TotalElements:    total,
		TotalPages:       totalPages,
		Size:             size,
		Number:           page,
		NumberOfElements: len(payments),
		First:            page == 0,
		Last:             page >= totalPages-1,
		Empty:            len(payments) == 0,
	}, nil
}

func (s *PaymentService) FindById(ctx context.Context, id int) (*model.PaymentDto, error) {
	payment, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

func (s *PaymentService) GetOrderDto(ctx context.Context, orderId int) *model.OrderDto {
	return &model.OrderDto{
		OrderId: &orderId,
	}
}

func (s *PaymentService) Save(ctx context.Context, payment model.PaymentDto) (*model.PaymentDto, error) {
	if payment.OrderId != nil {
		alreadyPaid, err := s.repo.ExistsByOrderIdAndIsPayed(ctx, payment.OrderId)
		if err != nil {
			return nil, err
		}
		if alreadyPaid {
			return nil, ErrOrderAlreadyPaid
		}
	}

	saved, err := s.repo.Save(ctx, payment)
	if err != nil {
		return nil, err
	}

	kafkaDto := model.KafkaPaymentDto{
		PaymentId:     saved.PaymentId,
		IsPayed:       saved.IsPayed,
		PaymentStatus: saved.PaymentStatus,
		OrderId:       saved.OrderId,
		UserId:        saved.UserId,
	}
	if msgBytes, err := json.Marshal(kafkaDto); err == nil {
		_ = s.producer.Send(ctx, TopicPaymentSuccessful, string(msgBytes))
		_ = s.producer.Send(ctx, TopicPaymentCompleted, string(msgBytes))
	}

	return saved, nil
}

func (s *PaymentService) Update(ctx context.Context, payment model.PaymentDto) (*model.PaymentDto, error) {
	if payment.PaymentId == nil {
		return nil, errors.New("paymentId is required for update")
	}
	return s.UpdateById(ctx, *payment.PaymentId, payment)
}

func (s *PaymentService) UpdateById(ctx context.Context, id int, payment model.PaymentDto) (*model.PaymentDto, error) {
	updated, err := s.repo.Update(ctx, id, payment)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrPaymentNotFound
	}
	return updated, nil
}

func (s *PaymentService) DeleteById(ctx context.Context, id int) error {
	return s.repo.DeleteById(ctx, id)
}
