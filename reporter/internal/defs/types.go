package defs

type ReportCode string

func (c ReportCode) ToString() string {
	return string(c)
}

type FormatterCode string

func (c FormatterCode) ToString() string {
	return string(c)
}

type StorageCode string

func (c StorageCode) ToString() string {
	return string(c)
}

type DownloadType string

func (t DownloadType) ToString() string {
	return string(t)
}
