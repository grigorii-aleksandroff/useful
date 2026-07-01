package main

import (
	"context"
	"fmt"
	_ "time/tzdata"

	appPkg "git.itechpsp.com/e46/box/platform.git/app"
	migratePkg "git.itechpsp.com/e46/box/platform.git/modules/migrate"
	serverPkg "git.itechpsp.com/e46/box/platform.git/pkg/micro/server"
	grpcServerPkg "git.itechpsp.com/e46/box/platform.git/servers/grpc"
	hcPkg "git.itechpsp.com/e46/box/platform.git/servers/healthcheck"

	contract "git.bububla.com/kilogramix/asia/contracts.git/proto/reporter"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/config"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/formatter"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/middleware"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/report"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/reportqueue"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/repository"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/scheduler"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/service"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/storage"
	grpcHandler "git.bububla.com/kilogramix/asia/reporter.git/internal/transport/grpc"
	migrations "git.bububla.com/kilogramix/asia/reporter.git/migrations"
)

var (
	tag    string
	commit string
)

const project = "asia"

func main() {
	app := appPkg.New(
		appPkg.Project(project),
		appPkg.Service(contract.ServiceName),
	)

	app.SetMigrations(migratePkg.New(
		migratePkg.Name(app.Service()),
		migratePkg.Logger(app.Logger()),
		migratePkg.FS(migrations.FS()),
		migratePkg.Db(app.Db()),
	))

	hc := hcPkg.New(hcPkg.Logger(app.Logger()), hcPkg.Config(app.ServiceConfig()))
	app.AddServer(hc)

	gs := grpcServerPkg.New(
		grpcServerPkg.Name(contract.ServiceName),
		grpcServerPkg.Logger(app.Logger()),
		grpcServerPkg.Config(app.ServiceConfig()),
		grpcServerPkg.Registry(app.Registry()),
	)

	serverPkg.NewOptions()
	_ = gs.Server().Init(serverPkg.WrapHandler(middleware.LoggingWrapper(app.Logger())))
	app.AddServer(gs)

	cfg := config.Load(app.ServiceConfig())

	repo := repository.NewReporterRepository(app.Db().Db())

	formatters, formatterCodes, err := formatter.Register(
		formatter.NewCSV(),
		formatter.NewXLSX(),
	)
	if err != nil {
		panic(err)
	}

	storages, storageCodes, err := storage.Register(cfg.Storages)
	if err != nil {
		panic(err)
	}

	registry, err := report.Register(cfg.Report, formatterCodes, storageCodes)
	if err != nil {
		panic(err)
	}

	publisher := reportqueue.NewJobsPublisher(repo.Jobs)
	exportService := service.NewExportService(registry, repo.Files, formatters, storages)
	reportService := service.NewReporterServices(repo, registry, formatters, storages)

	handler := grpcHandler.New(
		grpcHandler.Logger(app.Logger()),
		grpcHandler.Config(app.ServiceConfig()),
		grpcHandler.Service(reportService),
	)
	if err := contract.RegisterReporterServiceHandler(gs.Server(), handler); err != nil {
		panic(err)
	}

	app.Init()
	hc.SetLive()

	worker := reportqueue.NewWorker(app.Logger(), repo.Jobs, exportService, cfg.Worker.Interval, cfg.Worker.BatchSize)
	sched := scheduler.New(app.Logger(), repo.Schedules, publisher, cfg.Scheduler.Interval)

	go worker.Run(context.Background())
	go sched.Run(context.Background())

	app.Logger().Bg().Warn(fmt.Sprintf("service started tag %s commit %s", tag, commit))
	app.Run()
	app.Logger().Bg().Warn("service shutdown.")
}
