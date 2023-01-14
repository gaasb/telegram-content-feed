package utils

const DB_CONNECT_URI = "DB_CONNECT_URI"

type Repository interface {
	GetOne(id interface{}) interface{}
}
