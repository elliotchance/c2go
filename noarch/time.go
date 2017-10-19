package noarch

import "time"

type TimeT int32

func NullToTimeT(i int32) []TimeT {
	return []TimeT{}
}

func Time(tloc []TimeT) TimeT {
	var t = TimeT(int32(time.Now().Unix()))

	if len(tloc) > 0 {
		tloc[0] = t
	}

	return t
}

func IntToTimeT(t int32) TimeT {
	return TimeT(t)
}

func Ctime(tloc []TimeT) []byte {
	if len(tloc) > 0 {
		var t = time.Unix(int64(tloc[0]), 0)
		return []byte(t.Format(time.ANSIC) + "\n")
	}

	return nil
}
