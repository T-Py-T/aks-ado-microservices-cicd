// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
)

type placeOrderMockBackend struct {
	emptyCartErr error
}

type mockCartService struct {
	backend *placeOrderMockBackend
}

func (m *mockCartService) AddItem(context.Context, *pb.AddItemRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (m *mockCartService) GetCart(_ context.Context, req *pb.GetCartRequest) (*pb.Cart, error) {
	return &pb.Cart{
		UserId: req.UserId,
		Items:  []*pb.CartItem{{ProductId: "OLJCESPC7Z", Quantity: 1}},
	}, nil
}

func (m *mockCartService) EmptyCart(context.Context, *pb.EmptyCartRequest) (*pb.Empty, error) {
	if m.backend.emptyCartErr != nil {
		return nil, m.backend.emptyCartErr
	}
	return &pb.Empty{}, nil
}

type mockProductCatalogService struct{}

func (m *mockProductCatalogService) ListProducts(context.Context, *pb.Empty) (*pb.ListProductsResponse, error) {
	return &pb.ListProductsResponse{}, nil
}

func (m *mockProductCatalogService) GetProduct(_ context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	return &pb.Product{
		Id:   req.Id,
		Name: "test product",
		PriceUsd: &pb.Money{
			CurrencyCode: "USD",
			Units:        10,
			Nanos:        0,
		},
	}, nil
}

func (m *mockProductCatalogService) SearchProducts(context.Context, *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	return &pb.SearchProductsResponse{}, nil
}

type mockCurrencyService struct{}

func (m *mockCurrencyService) GetSupportedCurrencies(context.Context, *pb.Empty) (*pb.GetSupportedCurrenciesResponse, error) {
	return &pb.GetSupportedCurrenciesResponse{CurrencyCodes: []string{"USD"}}, nil
}

func (m *mockCurrencyService) Convert(_ context.Context, req *pb.CurrencyConversionRequest) (*pb.Money, error) {
	return &pb.Money{
		CurrencyCode: req.ToCode,
		Units:        req.From.Units,
		Nanos:        req.From.Nanos,
	}, nil
}

type mockPaymentService struct{}

func (m *mockPaymentService) Charge(context.Context, *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	return &pb.ChargeResponse{TransactionId: "test-tx"}, nil
}

type mockShippingService struct{}

func (m *mockShippingService) GetQuote(context.Context, *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
	return &pb.GetQuoteResponse{
		CostUsd: &pb.Money{CurrencyCode: "USD", Units: 5, Nanos: 0},
	}, nil
}

func (m *mockShippingService) ShipOrder(context.Context, *pb.ShipOrderRequest) (*pb.ShipOrderResponse, error) {
	return &pb.ShipOrderResponse{TrackingId: "TRACK-123"}, nil
}

type mockEmailService struct{}

func (m *mockEmailService) SendOrderConfirmation(context.Context, *pb.SendOrderConfirmationRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func startPlaceOrderMockBackend(t *testing.T, backend *placeOrderMockBackend) (*grpc.ClientConn, func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	gs := grpc.NewServer()
	pb.RegisterCartServiceServer(gs, &mockCartService{backend: backend})
	pb.RegisterProductCatalogServiceServer(gs, &mockProductCatalogService{})
	pb.RegisterCurrencyServiceServer(gs, &mockCurrencyService{})
	pb.RegisterPaymentServiceServer(gs, &mockPaymentService{})
	pb.RegisterShippingServiceServer(gs, &mockShippingService{})
	pb.RegisterEmailServiceServer(gs, &mockEmailService{})

	go gs.Serve(lis)

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		gs.Stop()
		t.Fatalf("grpc client: %v", err)
	}

	cleanup := func() {
		conn.Close()
		gs.Stop()
	}
	return conn, cleanup
}

func TestPlaceOrderFailsWhenEmptyCartFails(t *testing.T) {
	backend := &placeOrderMockBackend{
		emptyCartErr: status.Error(codes.Unavailable, "cart store unavailable"),
	}
	conn, cleanup := startPlaceOrderMockBackend(t, backend)
	defer cleanup()

	svc := &checkoutService{
		cartSvcConn:          conn,
		productCatalogSvcConn: conn,
		currencySvcConn:      conn,
		paymentSvcConn:       conn,
		shippingSvcConn:      conn,
		emailSvcConn:         conn,
	}

	_, err := svc.PlaceOrder(context.Background(), &pb.PlaceOrderRequest{
		UserId:       "user-1",
		UserCurrency: "USD",
		Email:        "buyer@example.com",
		Address: &pb.Address{
			StreetAddress: "1600 Amphitheatre Parkway",
			City:          "Mountain View",
			State:         "CA",
			Country:       "United States",
			ZipCode:       94043,
		},
		CreditCard: &pb.CreditCardInfo{
			CreditCardNumber:          "4111111111111111",
			CreditCardCvv:             123,
			CreditCardExpirationYear:  2030,
			CreditCardExpirationMonth: 12,
		},
	})
	if err == nil {
		t.Fatal("expected PlaceOrder to fail when EmptyCart fails after payment")
	}
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected Internal status, got %v", status.Code(err))
	}
}
