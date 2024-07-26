package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key     string
	Created time.Time
}

type RetentionMap map[string]OTP

func NewRetentionMap(ctx context.Context, retentionPeriond time.Duration) RetentionMap {
	rm := make(RetentionMap)

	// this will run in the background and delete the otp after the retention period has expired...
	go rm.Retention(ctx, retentionPeriond)

	return rm
}

func (rm RetentionMap) NewOTP() OTP {
	o := OTP{
		Key:     uuid.NewString(),
		Created: time.Now(),
	}

	rm[o.Key] = o
	return o
}

func (rm RetentionMap) VerifyOTP(otp string) bool {
	if _, ok := rm[otp]; !ok {
		return false
	}
	// delete the otp after it has been used...
	delete(rm, otp)
	return true
}

func (rm RetentionMap) Retention(ctx context.Context, retentionPeriod time.Duration) {
	// for every 400 milliseconds, check if the otp has expired or not...
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, otp := range rm {
				if time.Since(otp.Created) > retentionPeriod {
					delete(rm, otp.Key)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
