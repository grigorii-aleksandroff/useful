package grpc

import (
	"context"

	contract "git.bububla.com/kilogramix/asia/contracts.git/proto/reporter"
)

func (h *Handler) Ping(ctx context.Context, req *contract.HealthcheckRequest, resp *contract.HealthcheckResponse) error {
	h.logger.Info("Healthcheck ping...")

	resp.Message = "Pong"

	return nil
}
