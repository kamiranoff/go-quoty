package quoty

import (
	"fmt"
	"github.com/slack-go/slack"
	"go-quoty/pkg/quotes"
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
//				"text": "> Great"
//			}
//		},
//		{
//			"type": "section",
//			"text": {
//				"type": "mrkdwn",
//				"text": "*Benjamin Franklin*, _The Illiad_"
//			}
//		},
//		{
//					"type": "button",
//					"text": {
//						"type": "plain_text",
//						"text": "Ler Ros",
//						"emoji": true
//					},
//					"value": "click_me_123",
//					"url": "https://google.com"
//				}
//	]
//}

func buildCategoryMessage(categories []quotes.Category) string {
	var category quotes.Category
	if len(categories) > 0 {
        category = categories[0]
    }

	switch category {
	case quotes.Mythology:
		return ":zap:  " + strings.Title(string(category)) + "  :trident:"
	case quotes.Education:
		return ":mortar_board:  " + strings.Title(string(category)) + "  :books:"
	default:
		return ":books:  General  :book:"
	}
}

func buildHeaderSection(categories []quotes.Category) *slack.SectionBlock {
	categoryMsg := buildCategoryMessage(categories)
	headerBlock := slack.NewTextBlockObject(slack.MarkdownType, categoryMsg, false, false)
	return slack.NewSectionBlock(headerBlock, nil, nil)
}

func buildQuoteSection(quote string) *slack.SectionBlock {
	quoteBlock := slack.NewTextBlockObject(slack.MarkdownType, quote, false, false)
	return slack.NewSectionBlock(quoteBlock, nil, nil)
}

func buildFooterSection(author string, book *string) *slack.SectionBlock {
	authorText := "Unknown author"
	if author != "" {
		authorText = author
	}
	authorMsg := fmt.Sprintf("*%s*", authorText)

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

func buildButtonSection(id string) *slack.ActionBlock {
	actionId := "press_one_more"
	actionCancelId := "press_no_more"
	buttonId := "one_more"
	buttonText := "One more?"
	buttonTextBlock := slack.NewTextBlockObject(slack.PlainTextType, buttonText, false, false)

	buttonCancelId := id
	buttonCancelText := "No, thanks"
	buttonCancelTextBlock := slack.NewTextBlockObject(slack.PlainTextType, buttonCancelText, false, false)

	buttonBlockElement := slack.NewButtonBlockElement(actionId, buttonId, buttonTextBlock)
	buttonCancelBlockElement := slack.NewButtonBlockElement(actionCancelId, buttonCancelId, buttonCancelTextBlock)

	action := slack.NewActionBlock(buttonId, buttonBlockElement, buttonCancelBlockElement)

	return action
}

func BuildQuote(withButton bool, quoteId *string) slack.MsgOption {
	var quote quotes.Quote
	if quoteId != nil {
		quote = quotes.GetQuoteById(*quoteId)
    } else {
        quote = quotes.GetRandomQuote()
	}

	if withButton == true {
		return slack.MsgOptionBlocks(
			buildSeparatorSection(),
			buildHeaderSection(quote.Categories),
			buildSeparatorSection(),
			buildQuoteSection(quote.Quote),
			buildSeparatorSection(),
			buildFooterSection(quote.Author, quote.Book),
			buildSeparatorSection(),
			buildButtonSection(quote.Id),
		)
	}
	return slack.MsgOptionBlocks(
		buildSeparatorSection(),
		buildHeaderSection(quote.Categories),
		buildSeparatorSection(),
		buildQuoteSection(quote.Quote),
		buildSeparatorSection(),
		buildFooterSection(quote.Author, quote.Book),
		buildSeparatorSection(),
	)
}
