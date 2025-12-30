package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type GenerateResult struct {
	PlainText string
	TextyJSON map[string]interface{}
}

type SlackListGenerator struct{}

func NewSlackListGenerator() *SlackListGenerator {
	return &SlackListGenerator{}
}

func (g *SlackListGenerator) Generate(htmlContent string) (*GenerateResult, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	ops := []map[string]interface{}{}
	plainTextLines := []string{}

	var processList func(*html.Node, int)
	processList = func(listNode *html.Node, level int) {
		listType := "bullet"
		if listNode.Data == "ol" {
			listType = "ordered"
		}
		index := 1

		for child := listNode.FirstChild; child != nil; child = child.NextSibling {
			if child.Type != html.ElementNode {
				continue
			}

			if child.Data == "li" {
				currentLevel := level
				// aria-level check
				for _, attr := range child.Attr {
					if attr.Key == "aria-level" {
						if val, err := strconv.Atoi(attr.Val); err == nil {
							currentLevel = val - 1
						}
						break
					}
				}

				// Extract text and nested lists
				itemText := ""
				var nestedLists []*html.Node

				for c := child.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						itemText += c.Data
					} else if c.Type == html.ElementNode {
						if c.Data == "ul" || c.Data == "ol" {
							nestedLists = append(nestedLists, c)
						} else {
							itemText += extractText(c)
						}
					}
				}

				itemText = strings.TrimSpace(itemText)
				if itemText != "" {
					// Add operation
					ops = append(ops, map[string]interface{}{
						"insert": itemText,
					})

					attrs := map[string]interface{}{
						"list": listType,
					}
					if currentLevel > 0 {
						attrs["indent"] = currentLevel
					}
					ops = append(ops, map[string]interface{}{
						"attributes": attrs,
						"insert":     "\n",
					})

					// Plain text
					indentStr := strings.Repeat("    ", currentLevel)
					prefix := "- "
					if listType == "ordered" {
						prefix = fmt.Sprintf("%d. ", index)
					}
					plainTextLines = append(plainTextLines, indentStr+prefix+itemText)

					if listType == "ordered" {
						index++
					}
				}

				// Process nested lists
				for _, nested := range nestedLists {
					processList(nested, currentLevel+1)
				}

			} else if child.Data == "ul" || child.Data == "ol" {
				// Handle nested lists as siblings (Google Docs structure)
				processList(child, level+1)
			}
		}
	}

	// Find first list
	var firstList *html.Node
	var findList func(*html.Node)
	findList = func(n *html.Node) {
		if firstList != nil {
			return
		}
		if n.Type == html.ElementNode && (n.Data == "ul" || n.Data == "ol") {
			firstList = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findList(c)
		}
	}
	findList(doc)

	if firstList != nil {
		processList(firstList, 0)
	} else {
		// Fallback: treat as plain text
		text := extractText(doc)
		text = strings.TrimSpace(text)
		if text != "" {
			ops = append(ops, map[string]interface{}{"insert": text})
			ops = append(ops, map[string]interface{}{"insert": "\n"})
			plainTextLines = append(plainTextLines, text)
		}
	}

	return &GenerateResult{
		PlainText: strings.Join(plainTextLines, "\n"),
		TextyJSON: map[string]interface{}{"ops": ops},
	}, nil
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	// Skip nested lists in text extraction if they appear inside other tags
	// (Though usually they are direct children of li)
	if n.Type == html.ElementNode && (n.Data == "ul" || n.Data == "ol") {
		return ""
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return text
}
