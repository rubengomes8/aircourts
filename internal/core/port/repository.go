package port

type Repository interface {
	InsertMany(data []interface{}) string
}
