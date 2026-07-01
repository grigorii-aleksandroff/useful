package grpc

import (
	configPkg "git.itechpsp.com/e46/box/platform.git/modules/config"
	loggerPkg "git.itechpsp.com/e46/box/platform.git/modules/logger"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/service"
)

type Handler struct {
	logger  loggerPkg.Logger
	config  configPkg.Config
	service service.ReportService
}

type Option func(*Handler)

func Logger(logger loggerPkg.Logger) Option {
	return func(h *Handler) {
		h.logger = logger
	}
}

func Config(config configPkg.Config) Option {
	return func(h *Handler) {
		h.config = config
	}
}

func Service(reportService service.ReportService) Option {
	return func(h *Handler) {
		h.service = reportService
	}
}

func New(opts ...Option) *Handler {
	h := &Handler{}

	for _, opt := range opts {
		opt(h)
	}

	return h
}
