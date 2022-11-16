package sqlkit

import "time"

func TimeNowUTCStripNano() time.Time {
	t := time.Now().UTC()
	ns := t.Nanosecond() % 1000
	if ns > 0 {
		t = t.Add(-1 * time.Duration(ns) * time.Nanosecond)
		return t
	}
	return t
}

func TimeAddStripNano(t time.Time, d time.Duration) time.Time {
	u := t.Add(d)
	ns := u.Nanosecond()
	if ns > 9999 {
		ns %= 1000
		if ns > 0 {
			u = u.Add(-1 * time.Duration(ns) * time.Nanosecond)
			return u
		}
	}
	return u
}

func TimeSubStripNano(t time.Time, u time.Time) time.Duration {
	d := t.Sub(u)
	if d > 9999 {
		ns := d % 1000
		if ns > 0 {
			d = d - ns
			return d
		}
	}
	return d
}
