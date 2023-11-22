package types

type Polarity int

const (
	POSITIVE Polarity = iota
	NEGATIVE
	UNKNOWN
)

var PolarityMap = map[Polarity]string{
	POSITIVE: "+ve",
	NEGATIVE: "-ve",
	UNKNOWN:  "Unknown",
}

// Positive types: 1, *, +{...}, \/ (downshift)
// Negative types:   -*, &{...}, /\ (upshift)
