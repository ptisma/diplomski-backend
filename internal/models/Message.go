package models

const CSV_FILE_CODE = 1

const APSIMX_FILE_CODE = 2

const CONSTS_FILE_CODE = 0

type Message struct {
	ID      int
	Payload string
}
