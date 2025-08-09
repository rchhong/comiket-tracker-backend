package models

import (
	"fmt"
	"time"
)

type User struct {
	Discord_Id          int    `json:"discord_id"`
	Discord_Name        string `json:"discord_name"`
	Discord_Global_Name string `json:"discord_global_name"`
}

type UserWithMetadata struct {
	User
	Created_At time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func (user User) String() string {
	return fmt.Sprintf("{Discord_Id: %d Discord_Name: %s Discord_Global_Name: %s}\n", user.Discord_Id, user.Discord_Name, user.Discord_Global_Name)
}
