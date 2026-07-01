package helpers

import (
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ParseID(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

func TimestampToTime(value *timestamppb.Timestamp) *time.Time {
	if value == nil {
		return nil
	}

	t := value.AsTime()
	return &t
}
