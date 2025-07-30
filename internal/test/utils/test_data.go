package test_utils

import "flea-market/models"

// Todo: better to use func instead of global variables.
// global variables can be rewrite easily and unintentionally.
// Like the following codes
// db.Create(&test_utils.UserData)
// Passing a pointer, so gorm can rewrite data.
var ItemData = []models.Item{
	{Name: "test1", Price: 100, Description: "", SoldOut: false, UserID: 1},
	{Name: "test2", Price: 200, Description: "テスト2", SoldOut: true, UserID: 1},
	{Name: "test3", Price: 300, Description: "テスト3", SoldOut: false, UserID: 2},
}

var UserData = []models.User{
	{Email: "test1@test.com", Password: "testpass"},
	{Email: "test2@test.com", Password: "testpass"},
}
