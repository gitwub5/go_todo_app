package clock

import (
	"time"
)

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{} // 실제 시간을 반환하는 타입

func (r RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct{} // 테스트용 고정 시간을 반환하는 타입

func (fc FixedClocker) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
}
