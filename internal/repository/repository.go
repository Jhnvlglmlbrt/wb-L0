package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/Jhnvlglmlbrt/wb-order/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateTable() error
	SaveOrder(order models.Order) error
	GetAll() ([]models.Order, error)
}

type repo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repo {
	return &repo{
		pool: pool,
	}
}

func (r *repo) CreateTable() error {
	_, err := r.pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS main (
			order_uid VARCHAR(255) PRIMARY KEY,
			track_number VARCHAR(255),
			entry VARCHAR(255),
			delivery_info JSONB,
			payment_info JSONB,
			items JSONB,
			locale VARCHAR(255),
			internal_signature VARCHAR(255),
			customer_id VARCHAR(255),
		  	delivery_service VARCHAR(255),
		  	shardkey VARCHAR(255),
		  	sm_id INTEGER,
			date_created VARCHAR(255),
			oof_shard VARCHAR(255)
		)
		
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	return err
}

func (r *repo) SaveOrder(order models.Order) error {
	_, err := r.pool.Exec(context.Background(), `
		INSERT INTO main (order_uid, track_number, entry, delivery_info, payment_info, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, order.OrderUid, order.TrackNumber, order.Entry, order.Delivery, order.Payment, order.Items, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard)
	return err
}

func (r *repo) GetAll() ([]models.Order, error) {
	rows, err := r.pool.Query(context.Background(), `SELECT * FROM main`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}

	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var queryOrder models.Order
		if err := rows.Scan(
			&queryOrder.OrderUid,
			&queryOrder.TrackNumber,
			&queryOrder.Entry,
			&queryOrder.Delivery,
			&queryOrder.Payment,
			&queryOrder.Items,
			&queryOrder.Locale,
			&queryOrder.InternalSignature,
			&queryOrder.CustomerId,
			&queryOrder.DeliveryService,
			&queryOrder.ShardKey,
			&queryOrder.SmId,
			&queryOrder.DateCreated,
			&queryOrder.OofShard,
		); err != nil {
			log.Printf("Error at parsing preload data: %v\n", err)
			continue
		}
		orders = append(orders, queryOrder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error at scanning orders: %w", err)
	}
	return orders, nil
}
