package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/eaok/khlashe/config"
	"github.com/lonelyevil/khl"
	"github.com/phuslu/log"
)

func AddRole(roleID int64, roleName string, ctx *khl.ReactionAddContext) error {
	iFlagRoleExist := false

	log.Info().Str("MsgID", ctx.Extra.MsgID).Str("UserID", ctx.Extra.UserID).Str("TargetID", ctx.Common.TargetID).Str("EmojiName", ctx.Extra.Emoji.Name).Msg("AddRole")
	if ctx.Extra.MsgID == config.IDMsgRS {
		if ctx.Extra.Emoji.Name == config.RSEmoji[roleID] ||
			ctx.Extra.Emoji.Name == EmojiHexToDec(config.RSEmoji[roleID]) {
			user, _ := ctx.Session.UserView(ctx.Extra.UserID, ctx.Common.TargetID)
			fmt.Println(user.Roles)
			for _, role := range user.Roles {
				if int64(role) == roleID {
					iFlagRoleExist = true
				}
			}
			if !iFlagRoleExist {
				_, err := ctx.Session.GuildRoleGrant(ctx.Common.TargetID, ctx.Extra.UserID, roleID)
				if err != nil {
					log.Error().Err(errors.New("an error")).Msg("Role grant failed!")
					return err
				}

				resp, err := ctx.Session.MessageCreate(&khl.MessageCreate{
					MessageCreateBase: khl.MessageCreateBase{
						TargetID: ctx.Extra.ChannelID,
						Content:  fmt.Sprintf("You successfully joined %s!", roleName),
					},
				})
				if err != nil {
					return err
				}
				go func() {
					time.Sleep(2 * time.Second)
					ctx.Session.MessageDelete(resp.MsgID)
				}()
				// log.Info().Msgf("%s successfully joined %s!", user.Username, roleName)
			} else {
				resp, err := ctx.Session.MessageCreate(&khl.MessageCreate{
					MessageCreateBase: khl.MessageCreateBase{
						TargetID: ctx.Extra.ChannelID,
						Content:  fmt.Sprintf("You're already group %s!", roleName),
					},
				})
				if err != nil {
					return err
				}
				go func() {
					time.Sleep(2 * time.Second)
					ctx.Session.MessageDelete(resp.MsgID)
				}()
			}
		}
	}

	return nil
}

func DeleteRole(roleID int64, roleName string, ctx *khl.ReactionDeleteContext) error {
	iFlagRoleExist := false

	if ctx.Extra.MsgID == config.IDMsgRS {
		if ctx.Extra.Emoji.Name == config.RSEmoji[roleID] ||
			ctx.Extra.Emoji.Name == EmojiHexToDec(config.RSEmoji[roleID]) {
			user, _ := ctx.Session.UserView(ctx.Extra.UserID, ctx.Common.TargetID)
			fmt.Println(user.Roles)
			for _, role := range user.Roles {
				if int64(role) == roleID {
					iFlagRoleExist = true
				}
			}

			if iFlagRoleExist {
				_, err := ctx.Session.GuildRoleRevoke(ctx.Common.TargetID, ctx.Extra.UserID, roleID)
				if err != nil {
					return err
				}

				resp, err := ctx.Session.MessageCreate(&khl.MessageCreate{
					MessageCreateBase: khl.MessageCreateBase{
						TargetID: ctx.Extra.ChannelID,
						Content:  fmt.Sprintf("You successfully left group %s!", roleName),
					},
				})
				if err != nil {
					return err
				}
				go func() {
					time.Sleep(2 * time.Second)
					ctx.Session.MessageDelete(resp.MsgID)
				}()

				log.Info().Msgf("%s successfully left group %s!", user.Username, roleName)
			} else {
				resp, err := ctx.Session.MessageCreate(&khl.MessageCreate{
					MessageCreateBase: khl.MessageCreateBase{
						TargetID: ctx.Extra.ChannelID,
						Content:  fmt.Sprintf("You not in group %s!", roleName),
					},
				})
				if err != nil {
					return err
				}
				go func() {
					time.Sleep(2 * time.Second)
					ctx.Session.MessageDelete(resp.MsgID)
				}()
			}
		}
	}

	return nil
}

func AddRoles(ctx *khl.ReactionAddContext) error {
	if ctx.Extra.MsgID == config.IDMsgRS {
		log.Info().Str("EmojiName", ctx.Extra.Emoji.Name).Str("DecEmojiTem", EmojiHexToDec(config.EmojiTen)).Msg("AddRoles")
		switch ctx.Extra.Emoji.Name {
		case EmojiHexToDec(config.EmojiEleven):
			if err := AddRole(config.RoleRS11, "RS11", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("add role failed!")
				return err
			}
		case EmojiHexToDec(config.EmojiTen):
			if err := AddRole(config.RoleRS10, "RS10", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("add role failed!")
				return err
			}
		case config.EmojiNine:
			if err := AddRole(config.RoleRS9, "RS9", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("add role failed!")
				return err
			}
		case config.EmojiNeight:
			if err := AddRole(config.RoleRS8, "RS8", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("add role failed!")
				return err
			}
		case config.EmojiSeven:
			if err := AddRole(config.RoleRS7, "RS7", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("add role failed!")
				return err
			}
		case config.EmojiSix:
			if err := AddRole(config.RoleRS6, "RS6", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("add role failed!")
				return err
			}
		case config.EmojiFive:
			if err := AddRole(config.RoleRS5, "RS5", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("add role failed!")
				return err
			}
		case config.EmojiFour:
			if err := AddRole(config.RoleRS4, "RS4", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("add role failed!")
				return err
			}
		default:
		}
	}

	return nil
}

func DeleteRoles(ctx *khl.ReactionDeleteContext) error {
	if ctx.Extra.MsgID == config.IDMsgRS {
		log.Info().Str("EmojiName", ctx.Extra.Emoji.Name).Msg("AddRoles")

		switch ctx.Extra.Emoji.Name {
		case EmojiHexToDec(config.EmojiEleven):
			if err := DeleteRole(config.RoleRS11, "RS11", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("delete role failed!")
				return err
			}
		case EmojiHexToDec(config.EmojiTen):
			if err := DeleteRole(config.RoleRS10, "RS10", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("delete role failed!")
				return err
			}
		case config.EmojiNine:
			if err := DeleteRole(config.RoleRS9, "RS9", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("delete role failed!")
				return err
			}
		case config.EmojiNeight:
			if err := DeleteRole(config.RoleRS8, "RS8", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("delete role failed!")
				return err
			}
		case config.EmojiSeven:
			if err := DeleteRole(config.RoleRS7, "RS7", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("delete role failed!")
				return err
			}
		case config.EmojiSix:
			if err := DeleteRole(config.RoleRS6, "RS6", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("delete role failed!")
				return err
			}
		case config.EmojiFive:
			if err := DeleteRole(config.RoleRS5, "RS5", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("delete role failed!")
				return err
			}
		case config.EmojiFour:
			if err := DeleteRole(config.RoleRS4, "RS4", ctx); err != nil {
				log.Error().Err(errors.New("an error")).Msg("delete role failed!")
				return err
			}
		default:
		}
	}

	return nil
}
