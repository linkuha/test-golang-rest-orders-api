package entity

type Follower struct {
	UserID     string `json:"user_id"`
	FollowerID string `json:"follower_id"`
	//Status     int
}
