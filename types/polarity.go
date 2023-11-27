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

func (q *LabelType) Polarity() Polarity {
	// todo change to pass labelled environments
	panic("unfold type before checking for polarity")
	// return UNKNOWN
}

func (q *WIPType) Polarity() Polarity {
	return UNKNOWN
}

func (q *UnitType) Polarity() Polarity {
	return POSITIVE
}

func (q *SendType) Polarity() Polarity {
	return POSITIVE
}

func (q *ReceiveType) Polarity() Polarity {
	return NEGATIVE
}

func (q *SelectLabelType) Polarity() Polarity {
	return POSITIVE
}

func (q *BranchCaseType) Polarity() Polarity {
	return NEGATIVE
}

func (q *UpType) Polarity() Polarity {
	return NEGATIVE
}

func (q *DownType) Polarity() Polarity {
	return POSITIVE
}
