package tui

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type responseTreeRow struct {
	Path       string
	Key        string
	Value      string
	Depth      int
	Expandable bool
	Expanded   bool
	CanSave    bool
}

func buildResponseTreeRows(body string, expanded map[string]bool) []responseTreeRow {
	var data interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		body = strings.TrimSpace(body)
		if body == "" {
			return nil
		}
		return []responseTreeRow{{
			Path:    "$",
			Key:     "body",
			Value:   body,
			CanSave: true,
		}}
	}

	if expanded == nil {
		expanded = map[string]bool{"$": true}
	}

	rows := make([]responseTreeRow, 0)
	appendResponseTreeRows(&rows, "body", "$", data, 0, expanded)
	return rows
}

func appendResponseTreeRows(rows *[]responseTreeRow, key, path string, value interface{}, depth int, expanded map[string]bool) {
	switch v := value.(type) {
	case map[string]interface{}:
		isExpanded := expanded[path]
		*rows = append(*rows, responseTreeRow{
			Path:       path,
			Key:        key,
			Value:      fmt.Sprintf("Object {%d}", len(v)),
			Depth:      depth,
			Expandable: len(v) > 0,
			Expanded:   isExpanded,
		})
		if !isExpanded {
			return
		}

		keys := make([]string, 0, len(v))
		for childKey := range v {
			keys = append(keys, childKey)
		}
		sort.Strings(keys)

		for _, childKey := range keys {
			appendResponseTreeRows(rows, childKey, path+"."+childKey, v[childKey], depth+1, expanded)
		}
	case []interface{}:
		isExpanded := expanded[path]
		*rows = append(*rows, responseTreeRow{
			Path:       path,
			Key:        key,
			Value:      fmt.Sprintf("Array(%d)", len(v)),
			Depth:      depth,
			Expandable: len(v) > 0,
			Expanded:   isExpanded,
		})
		if !isExpanded {
			return
		}

		for i, child := range v {
			childKey := fmt.Sprintf("[%d]", i)
			appendResponseTreeRows(rows, childKey, fmt.Sprintf("%s[%d]", path, i), child, depth+1, expanded)
		}
	default:
		*rows = append(*rows, responseTreeRow{
			Path:    path,
			Key:     key,
			Value:   formatResponsePrimitive(v),
			Depth:   depth,
			CanSave: true,
		})
	}
}

func formatResponsePrimitive(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func selectedResponseRow(m Model) (responseTreeRow, bool) {
	rows := buildResponseTreeRows(m.ApiResponse.Body, m.ResponseExpanded)
	if m.pointer < 0 || m.pointer >= len(rows) {
		return responseTreeRow{}, false
	}
	return rows[m.pointer], true
}

func responseVariableKey(row responseTreeRow) string {
	key := strings.TrimPrefix(row.Path, "$.")
	key = strings.TrimPrefix(key, "$")
	if key == "" {
		return row.Key
	}
	return key
}

func visibleResponseWindow(m Model, rows []responseTreeRow) (int, int) {
	maxVisible := responseListHeight(m)
	start := m.responseScrollOffset
	if start < 0 {
		start = 0
	}
	if start > len(rows) {
		start = len(rows)
	}
	end := start + maxVisible
	if end > len(rows) {
		end = len(rows)
	}
	return start, end
}

func responseListHeight(m Model) int {
	return 5
}

func truncateText(text string, max int) string {
	if max <= 0 {
		return ""
	}

	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	if max <= 1 {
		return string(runes[:max])
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}
