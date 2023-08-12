package fun

// RealNumber is a generic number interface that covers all Go real number types.
type RealNumber interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

// Number is a generic number interface that covers all Go number types.
type Number interface {
	RealNumber | complex64 | complex128
}
