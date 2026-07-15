package gateway

import (
	"context"
	"fmt"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type PaymentGateway interface {
	Name() string
	Authorize(ctx context.Context, amount valueobject.Money, method valueobject.PaymentMethodType, ref string) (*GatewayResult, error)
	Capture(ctx context.Context, gatewayRef string, amount valueobject.Money) (*GatewayResult, error)
	Refund(ctx context.Context, gatewayRef string, amount valueobject.Money) (*GatewayResult, error)
	Cancel(ctx context.Context, gatewayRef string) error
}

type GatewayResult struct {
	Success    bool
	GatewayRef string
	AuthCode   string
	Error      string
}

type GatewayRegistry struct {
	gateways map[valueobject.PaymentMethodType]PaymentGateway
}

func NewGatewayRegistry() *GatewayRegistry {
	return &GatewayRegistry{gateways: make(map[valueobject.PaymentMethodType]PaymentGateway)}
}

func (r *GatewayRegistry) Register(method valueobject.PaymentMethodType, gw PaymentGateway) {
	r.gateways[method] = gw
}

func (r *GatewayRegistry) Get(method valueobject.PaymentMethodType) (PaymentGateway, error) {
	gw, ok := r.gateways[method]
	if !ok {
		return nil, fmt.Errorf("no gateway for method %s", method)
	}
	return gw, nil
}
