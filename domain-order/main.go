package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/zackyyhw/latihan-datacontract/contracts/order/v1"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	db *sql.DB
}

func (s *server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, customer_id, total, status, created_at FROM orders WHERE id = $1",
		req.OrderId,
	)
	var o pb.Order
	var createdAt time.Time
	err := row.Scan(&o.OrderId, &o.CustomerId, &o.Total, &o.Status, &createdAt)
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "order %s tidak ditemukan", req.OrderId)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}
	o.CreatedAt = createdAt.Unix()
	return &o, nil
}

func (s *server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	var rows *sql.Rows
	var err error
	if req.CustomerId == "" {
		rows, err = s.db.QueryContext(ctx,
			"SELECT id, customer_id, total, status, created_at FROM orders")
	} else {
		rows, err = s.db.QueryContext(ctx,
			"SELECT id, customer_id, total, status, created_at FROM orders WHERE customer_id = $1",
			req.CustomerId)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}
	defer rows.Close()

	var orders []*pb.Order
	for rows.Next() {
		var o pb.Order
		var createdAt time.Time
		rows.Scan(&o.OrderId, &o.CustomerId, &o.Total, &o.Status, &createdAt)
		o.CreatedAt = createdAt.Unix()
		orders = append(orders, &o)
	}
	return &pb.ListOrdersResponse{Orders: orders}, nil
}

func main() {
	dsn := os.Getenv("ORDER_DB_URL")
	if dsn == "" {
		log.Fatal("ORDER_DB_URL belum di-set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("gagal konek DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("DB tidak bisa diping: %v", err)
	}
	log.Println("Berhasil konek ke order-db")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("gagal listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{db: db})
	log.Println("Order domain running on :50051")
	s.Serve(lis)
}
