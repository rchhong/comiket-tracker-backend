package models

import (
	"fmt"
	"time"
)

type User struct {
	DiscordId         int    `json:"discord_id"`
	DiscordGlobalName string `json:"discord_global_name"`
}

type UserWithMetadata struct {
	User
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (user User) String() string {
	return fmt.Sprintf("{DiscordId: %d DiscordGlobalName: %s}\n", user.DiscordId, user.DiscordGlobalName)
}
