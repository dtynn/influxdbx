package internal

import "time"

// optional int64
func (oi *OptionalInt64) Int64Ptr() *int64 {
	if oi == nil {
		return nil
	}

	i := oi.GetVal()

	return &i
}

func (oi *OptionalInt64) IntPtr() *int {
	if oi == nil {
		return nil
	}

	i := int(oi.GetVal())

	return &i
}

func (oi *OptionalInt64) DurationPtr() *time.Duration {
	if oi == nil {
		return nil
	}

	i := time.Duration(oi.GetVal())

	return &i
}

// optional string
func (os *OptionalString) StringPtr() *string {
	if os == nil {
		return nil
	}

	s := os.GetVal()
	return &s
}
