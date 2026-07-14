package clock

import "time"

// CalculateAngles returns the angles (in degrees) for hour, minute, and second hands
// based on the given time. Angles follow standard clock convention:
// 0° = 12 o'clock, rotating clockwise.
//
//	secondAngle = second * 6            (360° / 60)
//	minuteAngle = minute * 6 + second * 0.1
//	hourAngle   = (hour % 12) * 30 + minute * 0.5
func CalculateAngles(t time.Time) (hourAngle, minuteAngle, secondAngle float64) {
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()

	secondAngle = float64(second) * 6.0
	minuteAngle = float64(minute)*6.0 + float64(second)*0.1
	hourAngle = float64(hour%12)*30.0 + float64(minute)*0.5

	return hourAngle, minuteAngle, secondAngle
}
