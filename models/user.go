package models

import (
	"fmt"
	"time"
)

type User struct {
	DiscordId         int    `json:"discord_id"`
	DiscordName       string `json:"discord_name"`
	DiscordGlobalName string `json:"discord_global_name"`
}

type UserWithMetadata struct {
	User
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (user User) String() string {
	return fmt.Sprintf("{DiscordId: %d DiscordName: %s DiscordGlobalName: %s}\n", user.DiscordId, user.DiscordName, user.DiscordGlobalName)
}
