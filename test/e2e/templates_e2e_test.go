package e2e

type User struct {
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
}

var User0 = User{
	FirstName: "User",
	LastName:  "One",
	Nickname:  "user1",
	Password:  "password1",
	Email:     "one@users.com",
	Country:   "US",
}

var User1 = User{
	FirstName: "User",
	LastName:  "Two",
	Nickname:  "user2",
	Password:  "password2",
	Email:     "two@user.com",
	Country:   "UK",
}

var User2 = User{
	FirstName: "User",
	LastName:  "Three",
	Nickname:  "user3",
	Password:  "password3",
	Email:     "three@user.com",
	Country:   "US",
}

var InvalidUser = User{
	FirstName: "User",
	LastName:  "",
	Nickname:  "userInvalid",
	Password:  "passwordInvalid",
	Email:     "invalid@user.com",
	Country:   "UK",
}

var UpdatedUser0 = User{
	FirstName: "User",
	LastName:  "One Updated",
	Nickname:  "user1",
	Password:  "password1",
	Email:     "one@users.com",
	Country:   "US",
}

var InvalidUpdatedUser1 = User{
	FirstName: "User",
	LastName:  "One",
	Nickname:  "user1",
	Password:  "",
	Email:     "one@users.com",
	Country:   "US",
}
