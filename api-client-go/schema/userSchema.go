package schema

type Username struct {
	Username string `json:"username"`
}

type NewPassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword 	string `json:"new_password"`
}

type SessionUser struct {
	Username string `json:"username"`
	IsAdmin  bool	`json:"is_admin"`
}

type User struct {
	UserId 			int 	`json:"user_id"`
	Username 		string  `json:"username"`
	Email 			string  `json:"email"`
	HashedPassword 	string  `json:"hashed_password"`
	IsAdmin 		bool 	`json:"is_admin"`
}

type NewUser struct {
	Token 	 string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
*---------------------------------------------------------- 
* Only for testing
*----------------------------------------------------------
*/
type TestNewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email 	 string `json:"email"`
}