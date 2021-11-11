package quoty

import (
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

//  Template
// {
//	"blocks": [
//		{
//			"type": "section",
//			"text": {
//				"type": "mrkdwn",
//				"text": ":zap: Mythology :trident:"
//			}
//		},
//		{
//			"type": "section",
//			"text": {
//				"type": "mrkdwn",
//				"text": "> Great to see you here! App helps you to stay up-to-date with your meetings and events right here within Slack. These are just a few things which you will be able to do"
//			}
//		},
//		{
//			"type": "section",
//			"text": {
//				"type": "mrkdwn",
//				"text": "*Benjamin Franklin*, _The Illiad_"
//			}
//		}
//	]
//}

func BuildQuote(category string, quote string, author string, book string) slack.MsgOption {
	separator := "\n\n"
	separatorBlock := slack.NewTextBlockObject(slack.MarkdownType, separator, false, false)
    separatorSection := slack.NewSectionBlock(separatorBlock, nil, nil)

	categoryMsg :=  ":zap:  " + strings.Title(category) + "  :trident:"
	headerBlock := slack.NewTextBlockObject(slack.MarkdownType, categoryMsg + separator, false, false)
	headerSection := slack.NewSectionBlock(headerBlock, nil, nil)

	quoteMsg := fmt.Sprintf("> %s", quote) + separator
	body := quoteMsg + separator
	bodyBlock := slack.NewTextBlockObject(slack.MarkdownType, body, false, false)
	quoteSection := slack.NewSectionBlock(bodyBlock, nil, nil)

	authorMsg := fmt.Sprintf("*%s*, _%s_", author, book) + separator
	footerBlock := slack.NewTextBlockObject(slack.MarkdownType, authorMsg, false, false)
	footerSection := slack.NewSectionBlock(footerBlock, nil, nil)

	msg := slack.MsgOptionBlocks(
		separatorSection,
		headerSection,
		separatorSection,
		quoteSection,
		separatorSection,
		footerSection,
		separatorSection,
	)

	return msg
}
