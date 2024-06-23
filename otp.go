package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key       string
	CreatedAt time.Time
}

type RetentionMap map[string]OTP

func NewRetantionMap(ctx context.Context, retentionPeriod time.Duration) RetentionMap {
	rm := make(RetentionMap)
	go rm.Retention(ctx, retentionPeriod)

	return rm
}

func (rm RetentionMap) NewOTP() OTP {
	o := OTP{
		Key:       uuid.NewString(),
		CreatedAt: time.Now(),
	}

	rm[o.Key] = o
	return o
}

func (rm RetentionMap) VerifyOTP(otp string) bool {
	if _, ok := rm[otp]; !ok {
		// otp not found...
		return false
	}

	delete(rm, otp)
	return true
}

func (rm RetentionMap) Retention(ctx context.Context, retentionPeriod time.Duration) {
	tiker := time.NewTicker(400 * time.Millisecond)

	for {
		select {
		case <-tiker.C:
			for _, otp := range rm {
				if time.Since(otp.CreatedAt) > retentionPeriod {
					delete(rm, otp.Key)
				}
			}

		case <-ctx.Done():
			return // exit the goroutine...
		}
	}
}
