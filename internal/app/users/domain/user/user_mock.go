package user

import "time"

var user1Now = time.Now()
var User1 = User{
	id:        "1",
	firstName: "John",
	lastName:  "Doe",
	nickname:  "john-123",
	password:  "password",
	email:     "me@john.com",
	country:   "US",
	createdAt: user1Now,
	updatedAt: user1Now,
}
