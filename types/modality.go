package types

type Modality int

const (
	UNRESTRICTED Modality = iota
	LINEAR
)

var ModalityMap = map[Modality]string{
	UNRESTRICTED: "U",
	LINEAR:       "L",
}
