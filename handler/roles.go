package handler

import (
	"fmt"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
)

func AddRole(roleID int64, roleName string, ctx *khl.ReactionAddContext) {
	iFlagRoleExist := false

	if ctx.Extra.Emoji.Name == config.RSEmoji[roleID] ||
		ctx.Extra.Emoji.Name == EmojiHexToDec(config.RSEmoji[roleID]) {
		user, err := ctx.Session.UserView(ctx.Extra.UserID, khl.UserViewWithGuildID(ctx.Common.TargetID))
		if err != nil {
			ctx.Session.Logger.Error().Err("", err).Msg("AddRole UserView")
		}
		for _, role := range user.Roles {
			if int64(role) == roleID {
				iFlagRoleExist = true
			}
		}
		ctx.Session.Logger.Warn().Interface("user.Roles", user.Roles).Bool("iFlagRoleExist", iFlagRoleExist).Msg("AddRole")

		if !iFlagRoleExist {
			_, err := ctx.Session.GuildRoleGrant(ctx.Common.TargetID, ctx.Extra.UserID, roleID)
			if err != nil {
				ctx.Session.Logger.Error().Err("", err).Msg("AddRole GuildRoleGrant")
				return
			}

			SendDirectMessage(ctx.Session, ctx.Extra.UserID, fmt.Sprintf("You successfully joined %s!", roleName))
		} else {
			SendDirectMessage(ctx.Session, ctx.Extra.UserID, fmt.Sprintf("You're already group %s!", roleName))
		}
	}
}

func DeleteRole(roleID int64, roleName string, ctx *khl.ReactionDeleteContext) {
	iFlagRoleExist := false

	if ctx.Extra.Emoji.Name == config.RSEmoji[roleID] ||
		ctx.Extra.Emoji.Name == EmojiHexToDec(config.RSEmoji[roleID]) {
		user, err := ctx.Session.UserView(ctx.Extra.UserID, khl.UserViewWithGuildID(ctx.Common.TargetID))
		if err != nil {
			ctx.Session.Logger.Error().Err("", err).Msg("DeleteRole UserView")
		}
		for _, role := range user.Roles {
			if int64(role) == roleID {
				iFlagRoleExist = true
			}
		}
		ctx.Session.Logger.Warn().Interface("user.Roles", user.Roles).Bool("iFlagRoleExist", iFlagRoleExist).Msg("DeleteRole")

		if iFlagRoleExist {
			_, err := ctx.Session.GuildRoleRevoke(ctx.Common.TargetID, ctx.Extra.UserID, roleID)
			if err != nil {
				ctx.Session.Logger.Error().Err("", err).Msg("DeleteRole GuildRoleRevoke")
				return
			}

			SendDirectMessage(ctx.Session, ctx.Extra.UserID, fmt.Sprintf("You successfully left group %s!", roleName))
		} else {
			SendDirectMessage(ctx.Session, ctx.Extra.UserID, fmt.Sprintf("You not in group %s!", roleName))
		}
	}
}

func AddRoles(ctx *khl.ReactionAddContext) error {
	if ctx.Extra.MsgID == config.Data.IDMsgRS {
		ctx.Session.Logger.Warn().Str("EmojiName", ctx.Extra.Emoji.Name).Str("DecEmojiTem", EmojiHexToDec(config.Data.EmojiTen)).Msg("AddRoles")

		switch ctx.Extra.Emoji.Name {
		case EmojiHexToDec(config.Data.EmojiEleven):
			AddRole(config.Data.RoleRS11, "RS11", ctx)
		case EmojiHexToDec(config.Data.EmojiTen):
			AddRole(config.Data.RoleRS10, "RS10", ctx)
		case config.Data.EmojiNine:
			AddRole(config.Data.RoleRS9, "RS9", ctx)
		case config.Data.EmojiNeight:
			AddRole(config.Data.RoleRS8, "RS8", ctx)
		case config.Data.EmojiSeven:
			AddRole(config.Data.RoleRS7, "RS7", ctx)
		case config.Data.EmojiSix:
			AddRole(config.Data.RoleRS6, "RS6", ctx)
		case config.Data.EmojiFive:
			AddRole(config.Data.RoleRS5, "RS5", ctx)
		case config.Data.EmojiFour:
			AddRole(config.Data.RoleRS4, "RS4", ctx)
		default:
		}
	} else if ctx.Extra.MsgID == config.Data.IDMsgBS {
		switch ctx.Extra.Emoji.Name {
		case config.Data.EmojiNeight:
			AddRole(config.Data.RoleBS8, "BS8", ctx)
		case config.Data.EmojiSeven:
			AddRole(config.Data.RoleBS7, "BS7", ctx)
		case config.Data.EmojiSix:
			AddRole(config.Data.RoleBS6, "BS6", ctx)
		case config.Data.EmojiFive:
			AddRole(config.Data.RoleBS5, "BS5", ctx)
		case config.Data.EmojiFour:
			AddRole(config.Data.RoleBS4, "BS4", ctx)
		case config.Data.EmojiThree:
			AddRole(config.Data.RoleBS3, "BS3", ctx)
		case config.Data.EmojiTwo:
			AddRole(config.Data.RoleBS2, "BS2", ctx)
		case config.Data.EmojiOne:
			AddRole(config.Data.RoleBS1, "BS1", ctx)
		default:
		}
	}

	return nil
}

func DeleteRoles(ctx *khl.ReactionDeleteContext) error {
	if ctx.Extra.MsgID == config.Data.IDMsgRS {
		ctx.Session.Logger.Warn().Str("EmojiName", ctx.Extra.Emoji.Name).Str("DecEmojiTem", EmojiHexToDec(config.Data.EmojiTen)).Msg("DeleteRoles")

		switch ctx.Extra.Emoji.Name {
		case EmojiHexToDec(config.Data.EmojiEleven):
			DeleteRole(config.Data.RoleRS11, "RS11", ctx)
		case EmojiHexToDec(config.Data.EmojiTen):
			DeleteRole(config.Data.RoleRS10, "RS10", ctx)
		case config.Data.EmojiNine:
			DeleteRole(config.Data.RoleRS9, "RS9", ctx)
		case config.Data.EmojiNeight:
			DeleteRole(config.Data.RoleRS8, "RS8", ctx)
		case config.Data.EmojiSeven:
			DeleteRole(config.Data.RoleRS7, "RS7", ctx)
		case config.Data.EmojiSix:
			DeleteRole(config.Data.RoleRS6, "RS6", ctx)
		case config.Data.EmojiFive:
			DeleteRole(config.Data.RoleRS5, "RS5", ctx)
		case config.Data.EmojiFour:
			DeleteRole(config.Data.RoleRS4, "RS4", ctx)
		default:
		}
	} else if ctx.Extra.MsgID == config.Data.IDMsgBS {
		switch ctx.Extra.Emoji.Name {
		case config.Data.EmojiNeight:
			DeleteRole(config.Data.RoleBS8, "BS8", ctx)
		case config.Data.EmojiSeven:
			DeleteRole(config.Data.RoleBS7, "BS7", ctx)
		case config.Data.EmojiSix:
			DeleteRole(config.Data.RoleBS6, "BS6", ctx)
		case config.Data.EmojiFive:
			DeleteRole(config.Data.RoleBS5, "BS5", ctx)
		case config.Data.EmojiFour:
			DeleteRole(config.Data.RoleBS4, "BS4", ctx)
		case config.Data.EmojiThree:
			DeleteRole(config.Data.RoleBS3, "BS3", ctx)
		case config.Data.EmojiTwo:
			DeleteRole(config.Data.RoleBS2, "BS2", ctx)
		case config.Data.EmojiOne:
			DeleteRole(config.Data.RoleBS1, "BS1", ctx)
		default:
		}
	}

	return nil
}
