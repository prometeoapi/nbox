package models

type CommandName string

var (
	UpsertTemplate CommandName = "upsert.template"
	UpsertVariable CommandName = "upsert.variables"
)

type Command[T any] struct {
	Command CommandName `json:"command"`
	Payload T           `json:"payload"`
	//Stage   string      `json:"stage"`
	//Payload string      `json:"payload"`
}
