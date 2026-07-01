package config

import (
	"time"

	configPkg "git.itechpsp.com/e46/box/platform.git/modules/config"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/report"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/storage"
)

type Config struct {
	Report    report.Config
	Storages  storage.Config
	Scheduler Scheduler
	Worker    Worker
}

type Scheduler struct {
	Interval time.Duration
}

type Worker struct {
	Interval  time.Duration
	BatchSize int
}

func Load(cfg configPkg.Config) Config {
	return Config{
		Report:    report.Config{Reports: reports(cfg)},
		Storages:  storages(cfg),
		Scheduler: scheduler(cfg),
		Worker:    worker(cfg),
	}
}

func scheduler(cfg configPkg.Config) Scheduler {
	return Scheduler{
		Interval: cfg.Get("scheduler", "interval").Duration(30 * time.Second),
	}
}

func worker(cfg configPkg.Config) Worker {
	return Worker{
		Interval:  cfg.Get("worker", "interval").Duration(5 * time.Second),
		BatchSize: cfg.Get("worker", "batch_size").Int(5),
	}
}

func reports(cfg configPkg.Config) map[defs.ReportCode]report.ReportConfig {
	result := map[defs.ReportCode]report.ReportConfig{}

	raw, ok := cfg.Get("reports").Interface(nil).(map[string]interface{})
	if !ok {
		return result
	}

	for code := range raw {
		result[defs.ReportCode(code)] = report.ReportConfig{
			Formatters: formatterCodes(cfg.Get("reports", code, "formatters").StringSlice(nil)),
			Storage:    defs.StorageCode(cfg.Get("reports", code, "storage").String(defs.StorageFile.ToString())),
		}
	}

	return result
}

func storages(cfg configPkg.Config) storage.Config {
	return storage.Config{
		File: storage.FileConfig{
			Dir:          cfg.Get("storages", "file", "dir").String("/tmp/reporter"),
			DownloadType: defs.DownloadType(cfg.Get("storages", "file", "download_type").String(defs.DownloadTypeFile.ToString())),
		},
	}
}

func formatterCodes(list []string) []defs.FormatterCode {
	result := make([]defs.FormatterCode, 0, len(list))
	for _, code := range list {
		result = append(result, defs.FormatterCode(code))
	}

	return result
}
