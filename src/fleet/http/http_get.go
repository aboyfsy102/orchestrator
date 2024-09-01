package main

import (
	"context"
	"encoding/json"
	"fmt"
	"j5v3/llib"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackc/pgx/v4"
)

func handleGet(ctx context.Context, queryParams map[string]string) (events.ALBTargetGroupResponse, error) {

	// if queryParams["id"] is empty, return all orders from RDS Postgres
	if queryParams["id"] == "" {
		orders, err := getOrders(ctx)
		if err != nil {
			return errorResponse(http.StatusInternalServerError, err.Error())
		}
		return orders, nil
	}

	// if queryParams["id"] is not empty, return the order with the id
	order, err := getOrder(ctx, queryParams["id"])
	if err != nil {
		return errorResponse(http.StatusInternalServerError, err.Error())
	}
	return order, nil
}

// New function to get a single order from the "orders" table
func getOrder(ctx context.Context, id string) (events.ALBTargetGroupResponse, error) {
	query := "SELECT * FROM orders WHERE id = $1"
	row := pgsqlClient.QueryRow(ctx, query, id)

	var order llib.FleetOrder                                         // Assuming you have an Order struct defined
	err := row.Scan(&order.ID, &order.CustomerName, &order.OrderDate) // Add more fields as needed

	if err != nil {
		if err == pgx.ErrNoRows {
			return events.ALBTargetGroupResponse{}, fmt.Errorf("order not found")
		}
		return events.ALBTargetGroupResponse{}, err
	}

	// Convert the order to JSON
	jsonResponse, err := json.Marshal(order)
	if err != nil {
		return events.ALBTargetGroupResponse{}, err
	}

	return events.ALBTargetGroupResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonResponse),
	}, nil
}
