package hub

type TmsError struct {
	Id       uint
	ErrorMsg string
	Module   string
	Line     int
	Error    error
}

const TmsErrorCoreId = 10000
const TmsErrorApisId = 20000
const TmsErrorUtilId = 30000
