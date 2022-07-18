package controllers

import (
	"assignment2/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type OrderController struct {
	db *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{
		db: db,
	}
}

func (o *OrderController) CreateOrder(ctx *gin.Context) {
	// create order and items
	var order models.Order

	err := ctx.ShouldBindJSON(&order)
	if err != nil {
		badRequestResponse(ctx, err.Error())
		return
	}

	// fmt.Printf("%+v <<< order\n", order)

	err = o.db.Create(&order).Error
	if err != nil {
		badRequestResponse(ctx, err.Error())
		return
	}

	writeJsonResponse(ctx, http.StatusCreated, gin.H{
		"success": true,
		"message": "created success",
	})
}

func (o *OrderController) GetOrders(ctx *gin.Context) {
	limit := ctx.Query("limit")
	limitInt := 10

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			limitInt = l
		}
	}

	var orders []models.Order
	var total int

	err := o.db.Preload("Items").Limit(limitInt).Find(&orders).Count(&total).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			notFoundResponse(ctx, err.Error())
			return
		}
		badRequestResponse(ctx, err.Error())
		return
	}

	writeJsonResponse(ctx, http.StatusOK, gin.H{
		"success": true,
		"payload": orders,
		"query": map[string]interface{}{
			"limit": limitInt,
			"total": total,
		},
	})
}

func (o *OrderController) GetOrderByID(ctx *gin.Context) {
	orderid := ctx.Param("orderid")

	var order models.Order

	err := o.db.Preload("Items").First(&order, orderid).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			notFoundResponse(ctx, err.Error())
			return
		}
		badRequestResponse(ctx, err.Error())
		return
	}

	writeJsonResponse(ctx, http.StatusOK, gin.H{
		"success": true,
		"payload": order,
	})
}

// func (o *OrderController) UpdateOrder(ctx *gin.Context) {
// 	orderid := ctx.Param("orderid")

// 	var order models.Order

// }
