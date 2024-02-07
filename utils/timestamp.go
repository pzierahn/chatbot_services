package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func TimeToProto(ts *time.Time) *timestamppb.Timestamp {
	if ts != nil {
		return timestamppb.New(*ts)
	}

	return nil
}

func ProtoToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts != nil {
		tmp := ts.AsTime()
		return &tmp
	}

	return nil
}
