package utils

const CSV_FILE_CODE = 1

const APSIMX_FILE_CODE = 2

const CONSTS_FILE_CODE = 0

//Helper message format struct for gouroutines who send the abs paths of newly created stage files
type Message struct {
	ID      int
	Payload string
}
