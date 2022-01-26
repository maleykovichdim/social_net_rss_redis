package common

import "time"

// Request structs
type CreateUserRequest struct {
	Name      string
	Surname   string
	Birthdate string
	Gender    string
	City      string
	Email     string
	Password  string
}

type LoginRequest struct {
	Email    string
	Password string
}

type InWrap struct {
	FriendId string `json:"friend_id"`
}

type InWrapTest struct {
	MyId     string `json:"my_id"`
	FriendId string `json:"friend_id"`
}

//Respons structs
type UserProfile struct {
	Id                int64  `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Surname           string `json:"surname,omitempty"`
	Birthdate         string `json:"birthdate,omitempty"`
	Gender            string `json:"gender,omitempty"`
	City              string `json:"city,omitempty"`
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	Avatar            string `json:"avatarURL,omitempty"`
	Has_personal_page bool   `json:"has_personal_page,omitempty"`
}

type IdResponse struct {
	Id int64 `json:"id"`
}

type PersonalPageResponse struct {
	Id        int64  `json:"id,omitempty"`
	UserId    int64  `json:"user_id,omitempty"`
	Interests string `json:"interests,omitempty"`
	About     string `json:"about,omitempty"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type FriendRequestResponse struct {
	UserId          int  `json:"user_id"`
	FriendUserId    int  `json:"friend_user_id"`
	RequestAccepted bool `json:"request_accepted"`
}

type Post struct {
	Id        int       `json:"id"`
	AuthorId  int       `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Follow struct {
	FollowerId int64 `json:"follower_id"`
	FolloweeId int64 `json:"followee_id"`
}
