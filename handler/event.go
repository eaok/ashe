package handler

import (
	"strconv"
	"strings"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
)

func CardButton(ctx *khl.MessageButtonClickContext) {
	uv, err := ctx.Session.UserView(ctx.Extra.UserID, khl.UserViewWithGuildID(ctx.Common.TargetID))
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("AddReaction UserView")
		return
	}
	ctx.Session.Logger.Warn().Str("value", ctx.Extra.Value).Str("MsgID", ctx.Extra.MsgID).Str("Username", uv.Username).Msg("CardButton")

	// team.ReactionAdd <- ctx
}

func AddReaction(ctx *khl.ReactionAddContext) {
	uv, err := ctx.Session.UserView(ctx.Extra.UserID, khl.UserViewWithGuildID(ctx.Common.TargetID))
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("AddReaction UserView")
		return
	}
	ctx.Session.Logger.Warn().Str("EmojiName", ctx.Extra.Emoji.Name).Str("MsgID", ctx.Extra.MsgID).Str("Username", uv.Username).Msg("AddReaction")
	// bot event ignore
	if uv.Bot {
		return
	}

	// 角色选择频道
	if ctx.Extra.ChannelID == config.Data.IDChannelSelectRole {
		if err := AddRoles(ctx); err != nil {
			ctx.Session.Logger.Error().Err("", err).Msg("AddReaction AddRoles")
			return
		}
	} else if ctx.Extra.ChannelID == config.Data.IDChannelBS {
		if ctx.Extra.Emoji.ID == config.Data.EmojiCheckMark {
			if _, ok := BSteam.MapBSAddGoroutine[ctx.Extra.MsgID]; ok {
				BSteam.MapBSAddGoroutine[ctx.Extra.MsgID] <- ctx
			}
		}
	} else if ctx.Extra.ChannelID == config.Data.IDChannelTradePublish {
		orderNumContent := GetMessageCoutent(ctx.Session, ctx.Extra.ChannelID, ctx.Extra.MsgID)
		if orderNumContent != "" {
			orderNumStr := strings.TrimPrefix(orderNumContent, "#")
			orderNum, _ := strconv.Atoi(orderNumStr)

			switch ctx.Extra.Emoji.ID {
			case config.Data.EmojiCrossMark:
				if _, ok := Trade.TradeAddCross[orderNum]; ok {
					Trade.TradeAddCross[orderNum] <- ctx
				}
			case config.EmojiNum[0], config.EmojiNum[1], config.EmojiNum[2], config.EmojiNum[3]:
				if _, ok := Trade.TradeAddNum[orderNum]; ok {
					Trade.TradeAddNum[orderNum] <- ctx
				}
			}
		}
	} else if ctx.Extra.ChannelID == config.Data.IDChannelTradeWait {
		if ctx.Extra.Emoji.ID == config.Data.EmojiCheckMark {
			orderNumContent := GetMessageCoutent(ctx.Session, ctx.Extra.ChannelID, ctx.Extra.MsgID)
			if orderNumContent != "" {
				orderNumStr := strings.TrimPrefix(orderNumContent, "#")
				orderNum, _ := strconv.Atoi(orderNumStr)
				if _, ok := Trade.TradeAddCheck[orderNum]; ok {
					Trade.TradeAddCheck[orderNum] <- ctx
				}
			}
		}
	}
}

func DeleteReaction(ctx *khl.ReactionDeleteContext) {
	uv, err := ctx.Session.UserView(ctx.Extra.UserID, khl.UserViewWithGuildID(ctx.Common.TargetID))
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("DeleteReaction UserView")
		return
	}
	ctx.Session.Logger.Warn().Str("EmojiName", ctx.Extra.Emoji.Name).Str("MsgID", ctx.Extra.MsgID).Str("Username", uv.Username).Msg("DeleteReaction")
	// bot event ignore
	if uv.Bot {
		return
	}

	// 角色选择频道
	if ctx.Extra.ChannelID == config.Data.IDChannelSelectRole {
		if err := DeleteRoles(ctx); err != nil {
			ctx.Session.Logger.Error().Err("", err).Msg("DeleteReaction DeleteRoles")
			return
		}
	} else if ctx.Extra.ChannelID == config.Data.IDChannelBS {
		if ctx.Extra.Emoji.ID == config.Data.EmojiCheckMark {
			if _, ok := BSteam.MapBSAddGoroutine[ctx.Extra.MsgID]; ok {
				BSteam.MapBSDeleteGoroutine[ctx.Extra.MsgID] <- ctx
			}
		}
	} else if ctx.Extra.ChannelID == config.Data.IDChannelTradeWait {
		if ctx.Extra.Emoji.ID == config.Data.EmojiCheckMark {
			orderNumContent := GetMessageCoutent(ctx.Session, ctx.Extra.ChannelID, ctx.Extra.MsgID)
			if orderNumContent != "" {
				orderNumStr := strings.TrimPrefix(orderNumContent, "#")
				orderNum, _ := strconv.Atoi(orderNumStr)
				if _, ok := Trade.TradeDeleteCheck[orderNum]; ok {
					Trade.TradeDeleteCheck[orderNum] <- ctx
				}
			}
		}
	}
}
