package controllers

import (
	"assignment2/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Order struct {
	Customer_Name string
	Ordered_At    time.Time
}

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

func (o *OrderController) UpdateOrderAndItems(ctx *gin.Context) {
	orderid := ctx.Param("orderid")

	var (
		orderAndItems models.Order
		order         models.Order
		newOrder      Order
	)

	err := ctx.ShouldBindJSON(&orderAndItems)
	if err != nil {
		badRequestResponse(ctx, err.Error())
		return
	}

	// fmt.Printf("%+v <<<items/n", newOrder.Items)

	// find by id
	err = o.db.First(&order, orderid).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			notFoundResponse(ctx, err.Error())
			return
		}
		badRequestResponse(ctx, err.Error())
		return
	}

	// update order by id
	newOrder.Customer_Name = orderAndItems.Customer_Name
	newOrder.Ordered_At = orderAndItems.Ordered_At
	err = o.db.Model(&order).Updates(newOrder).Error
	if err != nil {
		writeJsonResponse(ctx, http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// update item by orderId
	for _, itm := range orderAndItems.Items {
		var item models.Item
		// find item by orderid
		err = o.db.Where("id = ? AND order_id = ?", itm.ID, orderid).First(&item).Error
		if err == nil {
			// update item by orderid
			err = o.db.Model(&item).Updates(itm).Error
			if err != nil {
				writeJsonResponse(ctx, http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}
		}

	}

	writeJsonResponse(ctx, http.StatusOK, gin.H{
		"success": true,
		"message": "order and items - update success",
	})
}

func (o *OrderController) DeleteOrderAndItems(ctx *gin.Context) {
	orderid := ctx.Param("orderid")

	var order models.Order

	// fmt.Printf("%+v <<<items/n", newOrder.Items)

	// find by id
	err := o.db.First(&order, orderid).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			notFoundResponse(ctx, err.Error())
			return
		}
		badRequestResponse(ctx, err.Error())
		return
	}

	// delete order by id
	err = o.db.Delete(&order).Error
	if err != nil {
		writeJsonResponse(ctx, http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// delete item by orderId

	var item models.Item
	// find item by orderid
	err = o.db.Where("order_id = ?", orderid).Delete(&item).Error
	if err != nil {
		writeJsonResponse(ctx, http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	writeJsonResponse(ctx, http.StatusOK, gin.H{
		"success": true,
		"message": "order and items - delete success",
	})
}
