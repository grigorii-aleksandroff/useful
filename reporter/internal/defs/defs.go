package defs

const (
	ScheduleActive   = "active"
	ScheduleDisabled = "disabled"
	ScheduleDeleted  = "deleted"
)

const (
	JobNew        = "new"
	JobProcessing = "processing"
	JobComplete   = "complete"
	JobFailed     = "failed"
)

const (
	KindOnce    = "ONCE"
	KindDaily   = "DAILY"
	KindWeekly  = "WEEKLY"
	KindMonthly = "MONTHLY"
)

const (
	ReportOperators ReportCode = "operators"
)

const (
	FormatCSV  FormatterCode = "csv"
	FormatXLSX FormatterCode = "xlsx"
)

const (
	StorageFile StorageCode = "file"
)

const (
	DownloadTypeFile DownloadType = "file"
	DownloadTypeLink DownloadType = "link"
)
