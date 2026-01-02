package handler

import (
	"time"
)


type CreateUserRequest struct{
	Email string `json:"email"`
}

type UserResponse struct{
	ID	int64 			 `json:"id"`
	Email	string		 `json:"email"`
	CreatedAt time.Time	 `json:"created_at"`
}

