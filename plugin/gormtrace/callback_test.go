package gormtrace

import (
	"testing"
)

// type user struct {
// 	ID   int
// 	Name string
// }

//PLACEHOLDER FOR TESTS
func TestBasicGormTrace(t *testing.T) {
	// //Do DB set up
	// var db *gorm.DB

	// _, _, err := sqlmock.NewWithDSN("sqlmock_db_0")
	// if err != nil {
	// 	panic("Got an unexpected error.")
	// }

	// db, err = gorm.Open("sqlmock", "sqlmock_db_0")
	// if err != nil {
	// 	panic("Got an unexpected error.")
	// }

	// //Set up mocked tracer
	// var id uint64
	// tracer := mocktrace.MockTracer{StartSpanID: &id}

	// //Set up the gormtrace components
	// RegisterCallbacks(db, WithTracer(&tracer))
	// ctx := context.Background()
	// orm := WithContext(ctx, db)

	// //Create a fake user record
	// testUser := user{
	// 	Name: "John Smith",
	// }

	// //Initialize a trace
	// err = orm.Create(&testUser).Error
	// if err != nil {
	// 	panic("Got an unexpected error writing to DB")
	// }

	// //Run assertions
	// if got, expected := id, uint64(1); got != expected {
	// 	t.Fatalf("got %d, expected %d", got, expected)
	// }
}
