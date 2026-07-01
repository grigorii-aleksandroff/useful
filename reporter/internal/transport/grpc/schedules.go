package grpc

import (
	"context"
	"fmt"
	"strconv"

	microErrors "git.itechpsp.com/e46/box/platform.git/pkg/micro/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	contract "git.bububla.com/kilogramix/asia/contracts.git/proto/reporter"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/helpers"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/service"
)

func (h *Handler) CreateSchedule(ctx context.Context, req *contract.CreateScheduleRequest, resp *contract.CreateScheduleResponse) error {
	in := scheduleInput(req.Kind, req.ReportCode, req.TimeOfDay, req.DayOfWeek, req.DayOfMonth, req.RunAt, req.Enabled, req.Formats, req.Timezone, req.RunnerType, req.RunnerId)

	if err := validateScheduleInput(in); err != nil {
		h.logger.For(ctx).Error(fmt.Sprintf("create schedule validation: %v", err))
		return microErrors.BadRequest("validation", "invalid request: %v", err)
	}

	schedule, err := h.service.CreateSchedule(ctx, in)
	if err != nil {
		h.logger.For(ctx).Error(fmt.Sprintf("create schedule: %v", err))
		return microErrors.InternalServerError("create_schedule", "cannot create schedule: %v", err)
	}

	resp.Schedule = scheduleToProto(schedule)
	return nil
}

func (h *Handler) UpdateSchedule(ctx context.Context, req *contract.UpdateScheduleRequest, resp *contract.UpdateScheduleResponse) error {
	id, err := helpers.ParseID(req.Id)
	if err != nil {
		return microErrors.BadRequest("validation", "invalid id: %v", err)
	}

	in := scheduleInput(req.Kind, req.ReportCode, req.TimeOfDay, req.DayOfWeek, req.DayOfMonth, req.RunAt, req.Enabled, req.Formats, req.Timezone, "", "")

	if err := validateScheduleInput(in); err != nil {
		h.logger.For(ctx).Error(fmt.Sprintf("update schedule validation: %v", err))
		return microErrors.BadRequest("validation", "invalid request: %v", err)
	}

	schedule, err := h.service.UpdateSchedule(ctx, id, in)
	if err != nil {
		h.logger.For(ctx).Error(fmt.Sprintf("update schedule: %v", err))
		return microErrors.InternalServerError("update_schedule", "cannot update schedule: %v", err)
	}

	resp.Schedule = scheduleToProto(schedule)
	return nil
}

func (h *Handler) DeleteSchedule(ctx context.Context, req *contract.DeleteScheduleRequest, resp *contract.DeleteScheduleResponse) error {
	id, err := helpers.ParseID(req.Id)
	if err != nil {
		return microErrors.BadRequest("validation", "invalid id: %v", err)
	}

	if err := h.service.DeleteSchedule(ctx, id); err != nil {
		h.logger.For(ctx).Error(fmt.Sprintf("delete schedule: %v", err))
		return microErrors.InternalServerError("delete_schedule", "cannot delete schedule: %v", err)
	}

	resp.Success = true
	return nil
}

func (h *Handler) ListSchedules(ctx context.Context, req *contract.ListSchedulesRequest, resp *contract.ListSchedulesResponse) error {
	limit, offset := helpers.Pagination(req.Page, req.PerPage)

	schedules, count, err := h.service.ListSchedules(ctx, limit, offset, runnerFilters(req.Runners))
	if err != nil {
		h.logger.For(ctx).Error(fmt.Sprintf("list schedules: %v", err))
		return microErrors.InternalServerError("list_schedules", "cannot list schedules: %v", err)
	}

	resp.Data = make([]*contract.Schedule, 0, len(schedules))
	for i := range schedules {
		resp.Data = append(resp.Data, scheduleToProto(&schedules[i]))
	}
	resp.Count = count
	return nil
}

func (h *Handler) ListReportTypes(ctx context.Context, req *contract.ListReportTypesRequest, resp *contract.ListReportTypesResponse) error {
	types := h.service.ListReportTypes(ctx)

	resp.Data = make([]*contract.ReportType, 0, len(types))
	for _, reportType := range types {
		resp.Data = append(resp.Data, &contract.ReportType{
			Code: reportType.Code().ToString(),
		})
	}
	resp.Count = int64(len(types))
	return nil
}

func (h *Handler) ListReportFormats(ctx context.Context, req *contract.ListReportFormatsRequest, resp *contract.ListReportFormatsResponse) error {
	formats := h.service.Formatters(ctx, req.ReportCode)

	resp.Data = make([]*contract.ReportFormat, 0, len(formats))
	for _, format := range formats {
		resp.Data = append(resp.Data, &contract.ReportFormat{
			Code:      format.Code().ToString(),
			Extension: format.Extension(),
		})
	}
	resp.Count = int64(len(formats))
	return nil
}

func (h *Handler) ListFiles(ctx context.Context, req *contract.ListFilesRequest, resp *contract.ListFilesResponse) error {
	limit, offset := helpers.Pagination(req.Page, req.PerPage)

	files, count, err := h.service.ListFiles(ctx, limit, offset, runnerFilters(req.Runners))
	if err != nil {
		h.logger.For(ctx).Error(fmt.Sprintf("list files: %v", err))
		return microErrors.InternalServerError("list_files", "cannot list files: %v", err)
	}

	resp.Data = make([]*contract.ReportFile, 0, len(files))
	for i := range files {
		resp.Data = append(resp.Data, fileToProto(&files[i]))
	}
	resp.Count = count
	return nil
}

func (h *Handler) DownloadFile(ctx context.Context, req *contract.DownloadFileRequest, resp *contract.DownloadFileResponse) error {
	id, err := helpers.ParseID(req.Id)
	if err != nil {
		return microErrors.BadRequest("download_file", "invalid id: %v", err)
	}

	file, download, err := h.service.DownloadFile(ctx, id, runnerFilters(req.Runners))
	if err != nil {
		h.logger.For(ctx).Error(fmt.Sprintf("download file: %v", err))
		return microErrors.InternalServerError("download_file", "cannot download file: %v", err)
	}

	resp.Name = file.Name
	resp.Format = file.Format
	resp.Type = download.Type.ToString()
	resp.Content = download.Content
	resp.Link = download.Link
	return nil
}

func scheduleInput(
	kind contract.ScheduleKind,
	reportCode, timeOfDay string,
	dayOfWeek, dayOfMonth int32,
	runAt *timestamppb.Timestamp,
	enabled bool,
	formats []string,
	timezone string,
	runnerType, runnerID string,
) service.ScheduleInput {
	return service.ScheduleInput{
		Kind:       contract.ScheduleKind_name[int32(kind)],
		ReportCode: reportCode,
		TimeOfDay:  timeOfDay,
		DayOfWeek:  int(dayOfWeek),
		DayOfMonth: int(dayOfMonth),
		RunAt:      helpers.TimestampToTime(runAt),
		Enabled:    enabled,
		Formats:    formats,
		Timezone:   timezone,
		RunnerType: runnerType,
		RunnerID:   runnerID,
	}
}

func runnerFilters(runners []*contract.RunnerFilter) []helpers.RunnerFilter {
	result := make([]helpers.RunnerFilter, 0, len(runners))
	for _, runner := range runners {
		result = append(result, helpers.RunnerFilter{Type: runner.GetType(), ID: runner.GetId()})
	}
	return result
}

func scheduleToProto(schedule *model.Schedule) *contract.Schedule {
	result := &contract.Schedule{
		Id:         strconv.FormatUint(schedule.ID, 10),
		Kind:       contract.ScheduleKind(contract.ScheduleKind_value[schedule.Kind]),
		ReportCode: schedule.ReportCode,
		TimeOfDay:  schedule.TimeOfDay,
		DayOfWeek:  int32(schedule.DayOfWeek),
		DayOfMonth: int32(schedule.DayOfMonth),
		Status:     schedule.Status,
		Formats:    []string(schedule.Formats),
		Timezone:   schedule.Timezone,
		RunnerType: schedule.RunnerType,
		RunnerId:   schedule.RunnerID,
		CreatedAt:  timestamppb.New(schedule.CreatedAt),
		UpdatedAt:  timestamppb.New(schedule.UpdatedAt),
	}
	if schedule.NextRun != nil {
		result.NextRun = timestamppb.New(*schedule.NextRun)
	}
	return result
}

func fileToProto(file *model.File) *contract.ReportFile {
	return &contract.ReportFile{
		Id:         strconv.FormatUint(file.ID, 10),
		JobId:      strconv.FormatUint(file.JobID, 10),
		ReportCode: file.ReportCode,
		Name:       file.Name,
		Format:     file.Format,
		CreatedAt:  timestamppb.New(file.CreatedAt),
	}
}
