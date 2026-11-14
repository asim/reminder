package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/asim/reminder/daily"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

// NewMCP exposes Reminder's content and search tools through the MCP protocol.
func NewMCP(q *quran.Quran, n *names.Names, b *hadith.Collection, reminders *daily.Store, idx search.Searcher, indexed <-chan bool) *Server {
	mcpServer := NewServer("reminder", "1.0.0")

	mcpServer.AddTool("get_latest", "Get the latest reminder (updated hourly) with a verse, hadith and name of Allah", InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		resp := reminders.Latest()
		b, _ := json.Marshal(resp)
		return string(b), nil
	})

	mcpServer.AddTool("get_daily", "Get today's daily reminder with a verse, hadith and name of Allah", InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		resp := reminders.Today()
		b, _ := json.Marshal(resp)
		return string(b), nil
	})

	mcpServer.AddTool("get_daily_by_date", "Get a daily reminder for a specific date", InputSchema{
		Type: "object",
		Properties: map[string]Property{
			"date": {Type: "string", Description: "Date in YYYY-MM-DD format"},
		},
		Required: []string{"date"},
	}, func(args map[string]interface{}) (string, error) {
		date, _ := args["date"].(string)
		if date == "" {
			return "", fmt.Errorf("date is required")
		}
		if entry, ok := reminders.Get(date); ok {
			b, _ := json.Marshal(entry)
			return string(b), nil
		}
		return "", fmt.Errorf("not found")
	})

	mcpServer.AddTool("get_quran_chapters", "Get a list of all Quran chapters", InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		b, _ := json.Marshal(q.Index().Chapters)
		return string(b), nil
	})

	mcpServer.AddTool("get_quran_chapter", "Get a chapter of the Quran with all its verses", InputSchema{
		Type: "object",
		Properties: map[string]Property{
			"chapter": {Type: "number", Description: "Chapter number (1-114)"},
		},
		Required: []string{"chapter"},
	}, func(args map[string]interface{}) (string, error) {
		chNum, ok := args["chapter"].(float64)
		if !ok {
			return "", fmt.Errorf("chapter is required")
		}
		chapter := int(chNum)
		if chapter < 1 || chapter > 114 {
			return "", fmt.Errorf("chapter must be between 1 and 114")
		}
		return string(q.Get(chapter).JSON()), nil
	})

	mcpServer.AddTool("get_quran_verse", "Get a specific verse of the Quran with word-by-word translation", InputSchema{
		Type: "object",
		Properties: map[string]Property{
			"chapter": {Type: "number", Description: "Chapter number (1-114)"},
			"verse":   {Type: "number", Description: "Verse number"},
		},
		Required: []string{"chapter", "verse"},
	}, func(args map[string]interface{}) (string, error) {
		chNum, ok := args["chapter"].(float64)
		if !ok {
			return "", fmt.Errorf("chapter is required")
		}
		veNum, ok := args["verse"].(float64)
		if !ok {
			return "", fmt.Errorf("verse is required")
		}
		chapter := int(chNum)
		verse := int(veNum)
		if chapter < 1 || chapter > 114 {
			return "", fmt.Errorf("chapter must be between 1 and 114")
		}
		cc := q.Get(chapter)
		if cc.Verse(verse) == nil {
			return "", fmt.Errorf("verse out of range")
		}
		return string(cc.Verse(verse).JSON()), nil
	})

	mcpServer.AddTool("get_hadith_books", "Get a list of all hadith books in Sahih Bukhari", InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		bks, _ := json.Marshal(b.Index().Books)
		return string(bks), nil
	})

	mcpServer.AddTool("get_hadith_book", "Get a specific book from Sahih Bukhari with all its hadiths", InputSchema{
		Type: "object",
		Properties: map[string]Property{
			"book": {Type: "number", Description: "Book number"},
		},
		Required: []string{"book"},
	}, func(args map[string]interface{}) (string, error) {
		bkNum, ok := args["book"].(float64)
		if !ok {
			return "", fmt.Errorf("book is required")
		}
		book := int(bkNum)
		if book < 1 || book > len(b.Books) {
			return "", fmt.Errorf("book out of range")
		}
		return string(b.Get(book).JSON()), nil
	})

	mcpServer.AddTool("get_names", "Get all 99 Names of Allah", InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		return string(n.JSON()), nil
	})

	mcpServer.AddTool("get_name", "Get a specific Name of Allah with description and Quran references", InputSchema{
		Type: "object",
		Properties: map[string]Property{
			"id": {Type: "number", Description: "Name number (1-99)"},
		},
		Required: []string{"id"},
	}, func(args map[string]interface{}) (string, error) {
		idNum, ok := args["id"].(float64)
		if !ok {
			return "", fmt.Errorf("id is required")
		}
		id := int(idNum)
		if id < 1 || id > len(*n) {
			return "", fmt.Errorf("id must be between 1 and 99")
		}
		byt, _ := json.Marshal(n.Get(id))
		return string(byt), nil
	})

	mcpServer.AddContextTool("search", "Search Islamic content and get AI-summarised answers from the Quran, Hadith and Names of Allah", InputSchema{
		Type: "object",
		Properties: map[string]Property{
			"summarise": {Type: "boolean", Description: "Generate an AI summary; defaults to false"},
			"q":         {Type: "string", Description: "The question to ask"},
		},
		Required: []string{"q"},
	}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		select {
		case <-indexed:
		default:
			return "", fmt.Errorf("indexing content, please try again later")
		}

		question, ok := args["q"].(string)
		if !ok || question == "" {
			return "", fmt.Errorf("q is required")
		}

		res, err := search.References(idx, n, question)
		if err != nil {
			return "", err
		}
		res = search.CurrentTafsirResults(res, q)
		answer := ""
		if value, present := args["summarise"]; present {
			enabled, ok := value.(bool)
			if !ok {
				return "", fmt.Errorf("summarise must be true or false")
			}
			if enabled {
				result, err := model.RespondPublic(ctx, model.Request{Prompt: question}, func(query string) ([]*model.Reference, error) { return search.ModelReferences(idx, n, q, b, query) })
				if err != nil {
					return "", err
				}
				answer, res = result.Answer, result.References
			}
		}

		output, _ := json.Marshal(map[string]interface{}{
			"q":          question,
			"answer":     answer,
			"references": res,
		})
		return string(output), nil
	})

	return mcpServer
}
