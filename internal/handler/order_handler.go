package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/api-gateway/internal/proxy"
	pb "github.com/yourusername/order-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandler struct {
	orderService *proxy.OrderServiceClient
}

func NewOrderHandler(orderService *proxy.OrderServiceClient) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
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

	order, err := h.orderService.CreateOrder(c.Request.Context(), req)
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

	order, err := h.orderService.GetOrder(c.Request.Context(), req)
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

	order, err := h.orderService.UpdateOrderStatus(c.Request.Context(), req)
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

	result, err := h.orderService.AddDeliveryAddress(c.Request.Context(), &address)
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

	result, err := h.orderService.ListDeliveryAddresses(c.Request.Context(), req)
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

	result, err := h.orderService.GetAvailableDeliverySlots(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
} 