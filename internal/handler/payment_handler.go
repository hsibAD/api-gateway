package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/api-gateway/internal/proxy"
	pb "github.com/yourusername/payment-service/proto"
)

type PaymentHandler struct {
	paymentService *proxy.PaymentServiceClient
}

func NewPaymentHandler(paymentService *proxy.PaymentServiceClient) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	var request struct {
		OrderID       string  `json:"order_id"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		PaymentMethod string  `json:"payment_method"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	req := &pb.InitiatePaymentRequest{
		OrderId:       request.OrderID,
		UserId:        userID.(string),
		Amount:        request.Amount,
		Currency:      request.Currency,
		PaymentMethod: pb.PaymentMethod(pb.PaymentMethod_value[request.PaymentMethod]),
	}

	payment, err := h.paymentService.InitiatePayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) ProcessCreditCardPayment(c *gin.Context) {
	var request struct {
		PaymentID      string `json:"payment_id"`
		CardNumber     string `json:"card_number"`
		ExpiryMonth    string `json:"expiry_month"`
		ExpiryYear     string `json:"expiry_year"`
		CVV            string `json:"cvv"`
		CardholderName string `json:"cardholder_name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &pb.CreditCardPaymentRequest{
		PaymentId: request.PaymentID,
		CardInfo: &pb.CreditCardInfo{
			CardNumber:     request.CardNumber,
			ExpiryMonth:    request.ExpiryMonth,
			ExpiryYear:     request.ExpiryYear,
			Cvv:            request.CVV,
			CardholderName: request.CardholderName,
		},
	}

	payment, err := h.paymentService.ProcessCreditCardPayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) InitiateMetaMaskPayment(c *gin.Context) {
	var request struct {
		PaymentID     string `json:"payment_id"`
		WalletAddress string `json:"wallet_address"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &pb.MetaMaskPaymentRequest{
		PaymentId:     request.PaymentID,
		WalletAddress: request.WalletAddress,
	}

	response, err := h.paymentService.InitiateMetaMaskPayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PaymentHandler) ConfirmMetaMaskPayment(c *gin.Context) {
	var request struct {
		PaymentID       string `json:"payment_id"`
		TransactionHash string `json:"transaction_hash"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &pb.ConfirmMetaMaskPaymentRequest{
		PaymentId:       request.PaymentID,
		TransactionHash: request.TransactionHash,
	}

	payment, err := h.paymentService.ConfirmMetaMaskPayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")
	req := &pb.GetPaymentRequest{
		PaymentId: paymentID,
	}

	payment, err := h.paymentService.GetPayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetPaymentsByOrder(c *gin.Context) {
	orderID := c.Param("order_id")
	req := &pb.GetPaymentsByOrderRequest{
		OrderId: orderID,
	}

	payments, err := h.paymentService.GetPaymentsByOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) GetPendingPayments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &pb.GetPendingPaymentsRequest{
		UserId: userID.(string),
		Page:   int32(page),
		Limit:  int32(limit),
	}

	payments, err := h.paymentService.GetPendingPayments(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
} 