package quoty

import (
	"fmt"
	"github.com/slack-go/slack"
	"go-quoty/pkg/quotes"
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


func buildCategoryMessage(category quotes.Category) string {
	switch category {
	case quotes.Mythology:
        return ":zap:  " + string(category) + "  :trident:"
	case quotes.Education:
		return ":mortar_board:  " + string(category) + "  :books:"
	default:
		return ":star:  General  :star:"
    }
}

func buildHeaderSection(category quotes.Category) *slack.SectionBlock {
	categoryMsg := buildCategoryMessage(category)
	headerBlock := slack.NewTextBlockObject(slack.MarkdownType, categoryMsg, false, false)
	return slack.NewSectionBlock(headerBlock, nil, nil)
}

func buildQuoteSection(quote string) *slack.SectionBlock {
	quoteBlock := slack.NewTextBlockObject(slack.MarkdownType, quote, false, false)
	return slack.NewSectionBlock(quoteBlock, nil, nil)
}

func buildFooterSection(author string, book *string) *slack.SectionBlock {
	authorMsg := fmt.Sprintf("*%s*", author)

	if book != nil {
		authorMsg = authorMsg + fmt.Sprintf(", _%s_", *book)
	}
	authorBlock := slack.NewTextBlockObject(slack.MarkdownType, authorMsg, false, false)
	return slack.NewSectionBlock(authorBlock, nil, nil)
}

func buildSeparatorSection() *slack.SectionBlock {
	separator := "\n\n"
	separatorBlock := slack.NewTextBlockObject(slack.MarkdownType, separator, false, false)
	return slack.NewSectionBlock(separatorBlock, nil, nil)
}

func BuildQuote(category quotes.Category, quote string, author string, book *string) slack.MsgOption {
	return slack.MsgOptionBlocks(
		buildSeparatorSection(),
		buildHeaderSection(category),
		buildSeparatorSection(),
		buildQuoteSection(quote),
		buildSeparatorSection(),
		buildFooterSection(author, book),
		buildSeparatorSection(),
	)
}
