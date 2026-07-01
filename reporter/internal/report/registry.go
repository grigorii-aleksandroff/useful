package report

import (
	"fmt"
	"sync"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
)

type Registry interface {
	Register(reportType ReportType)
	Get(code defs.ReportCode) (ReportType, bool)
	List() []ReportType
}

type registry struct {
	mu    sync.RWMutex
	types map[defs.ReportCode]ReportType
	order []defs.ReportCode
}

func newRegistry() *registry {
	return &registry{types: map[defs.ReportCode]ReportType{}}
}

func (r *registry) Register(reportType ReportType) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.types[reportType.Code()]; !ok {
		r.order = append(r.order, reportType.Code())
	}
	r.types[reportType.Code()] = reportType
}

func (r *registry) Get(code defs.ReportCode) (ReportType, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reportType, ok := r.types[code]
	return reportType, ok
}

func (r *registry) List() []ReportType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]ReportType, 0, len(r.order))
	for _, code := range r.order {
		result = append(result, r.types[code])
	}
	return result
}

var defaultRegistry Registry = newRegistry()

func Register(
	cfg Config,
	formatterCodes map[defs.FormatterCode]bool,
	storageCodes map[defs.StorageCode]bool,
) (Registry, error) {
	built := newRegistry()

	for code, reportConfig := range cfg.Reports {
		prototype, ok := prototypes[code]
		if !ok {
			return nil, fmt.Errorf("report type %q is not registered", code)
		}

		for _, formatterCode := range reportConfig.Formatters {
			if !formatterCodes[formatterCode] {
				return nil, fmt.Errorf("report type %q references unknown formatter %q", code, formatterCode)
			}
		}

		if !storageCodes[reportConfig.Storage] {
			return nil, fmt.Errorf("report type %q references unknown storage %q", code, reportConfig.Storage)
		}

		built.Register(prototype.Init(reportConfig.Formatters, reportConfig.Storage))
	}

	defaultRegistry = built
	return built, nil
}

var prototypes = map[defs.ReportCode]ReportType{}

func registerType(reportType ReportType) {
	prototypes[reportType.Code()] = reportType
}
