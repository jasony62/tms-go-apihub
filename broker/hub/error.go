package hub

type TmsError struct {
	Id       uint32
	ErrorMsg string
	Module   string
	Line     uint32
}
