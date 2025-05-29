package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hsibAD/api-gateway/internal/proxy"
	pb "github.com/hsibAD/order-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandler struct {
	orderClient *proxy.OrderServiceClient
}

func NewOrderHandler(orderClient *proxy.OrderServiceClient) *OrderHandler {
	return &OrderHandler{
		orderClient: orderClient,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request struct {
		Items            []pb.OrderItem `json:"items"`
		DeliveryAddressID string       `json:"delivery_address_id"`
		DeliveryTime     int64         `json:"delivery_time"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	req := &pb.CreateOrderRequest{
		UserID:           userID.(string),
		Items:            request.Items,
		DeliveryAddressID: request.DeliveryAddressID,
		DeliveryTime:     timestamppb.New(timestamppb.Now().AsTime()),
	}

	order, err := h.orderClient.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	req := &pb.GetOrderRequest{
		OrderId: orderID,
	}

	order, err := h.orderClient.GetOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &pb.UpdateOrderStatusRequest{
		OrderId: orderID,
		Status:  pb.OrderStatus(pb.OrderStatus_value[request.Status]),
	}

	order, err := h.orderClient.UpdateOrderStatus(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) AddDeliveryAddress(c *gin.Context) {
	var address pb.DeliveryAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	address.UserId = userID.(string)

	result, err := h.orderClient.AddDeliveryAddress(c.Request.Context(), &address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *OrderHandler) ListDeliveryAddresses(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &pb.ListAddressesRequest{
		UserId: userID.(string),
		Page:   int32(page),
		Limit:  int32(limit),
	}

	result, err := h.orderClient.ListDeliveryAddresses(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *OrderHandler) GetAvailableDeliverySlots(c *gin.Context) {
	date := c.Query("date")
	req := &pb.GetDeliverySlotsRequest{
		Date: timestamppb.Now(),
	}

	result, err := h.orderClient.GetAvailableDeliverySlots(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
} 