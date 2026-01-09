package src

import (
	"reflect"
	"strings"
	"testing"
)

func TestSlackListGenerator_Generate(t *testing.T) {
	tests := []struct {
		name          string
		htmlContent   string
		wantPlainText string
		wantOps       []map[string]interface{}
	}{
		{
			name: "Simple Bullet List",
			htmlContent: `
				<ul>
					<li>Item 1</li>
					<li>Item 2</li>
				</ul>
			`,
			wantPlainText: "- Item 1\n- Item 2",
			wantOps: []map[string]interface{}{
				{"insert": "Item 1"},
				{"attributes": map[string]interface{}{"list": "bullet"}, "insert": "\n"},
				{"insert": "Item 2"},
				{"attributes": map[string]interface{}{"list": "bullet"}, "insert": "\n"},
			},
		},
		{
			name: "Simple Ordered List",
			htmlContent: `
				<ol>
					<li>First</li>
					<li>Second</li>
				</ol>
			`,
			wantPlainText: "1. First\n2. Second",
			wantOps: []map[string]interface{}{
				{"insert": "First"},
				{"attributes": map[string]interface{}{"list": "ordered"}, "insert": "\n"},
				{"insert": "Second"},
				{"attributes": map[string]interface{}{"list": "ordered"}, "insert": "\n"},
			},
		},
		{
			name: "Nested List (Google Docs Style - Flat with aria-level)",
			htmlContent: `
				<ul>
					<li aria-level="1">Level 1</li>
					<li aria-level="2">Level 2</li>
					<li aria-level="1">Level 1 again</li>
				</ul>
			`,
			wantPlainText: "- Level 1\n    - Level 2\n- Level 1 again",
			wantOps: []map[string]interface{}{
				{"insert": "Level 1"},
				{"attributes": map[string]interface{}{"list": "bullet"}, "insert": "\n"},
				{"insert": "Level 2"},
				{"attributes": map[string]interface{}{"indent": 1, "list": "bullet"}, "insert": "\n"},
				{"insert": "Level 1 again"},
				{"attributes": map[string]interface{}{"list": "bullet"}, "insert": "\n"},
			},
		},
		{
			name: "Nested List (Standard HTML Structure)",
			htmlContent: `
				<ul>
					<li>Parent
						<ul>
							<li>Child</li>
						</ul>
					</li>
				</ul>
			`,
			wantPlainText: "- Parent\n    - Child",
			wantOps: []map[string]interface{}{
				{"insert": "Parent"},
				{"attributes": map[string]interface{}{"list": "bullet"}, "insert": "\n"},
				{"insert": "Child"},
				{"attributes": map[string]interface{}{"indent": 1, "list": "bullet"}, "insert": "\n"},
			},
		},
		{
			name:          "Plain Text Fallback",
			htmlContent:   `<p>Just some text</p>`,
			wantPlainText: "Just some text",
			wantOps: []map[string]interface{}{
				{"insert": "Just some text"},
				{"insert": "\n"},
			},
		},
		{
			name: "Multiple Separated Lists",
			htmlContent: `
				<ul>
					<li>List 1</li>
				</ul>
				<br>
				<ul>
					<li>List 2</li>
				</ul>
			`,
			wantPlainText: "- List 1\n- List 2",
			wantOps: []map[string]interface{}{
				{"insert": "List 1"},
				{"attributes": map[string]interface{}{"list": "bullet"}, "insert": "\n"},
				{"insert": "List 2"},
				{"attributes": map[string]interface{}{"list": "bullet"}, "insert": "\n"},
			},
		},
		{
			name: "Mixed Content (Text and Lists)",
			htmlContent: `
				<p>not-list-1</p>
				<ul>
					<li aria-level="1">level1-1</li>
				</ul>
				<ul>
					<li aria-level="2">level2</li>
				</ul>
				<ul>
					<li aria-level="1">level1-2</li>
				</ul>
				<p>not-list-2</p>
			`,
			wantPlainText: "not-list-1\n- level1-1\n    - level2\n- level1-2\nnot-list-2",
			wantOps: []map[string]interface{}{
				{"insert": "not-list-1"},
				{"insert": "\n"},
				{"insert": "level1-1"},
				{"attributes": map[string]interface{}{"list": "bullet"}, "insert": "\n"},
				{"insert": "level2"},
				{"attributes": map[string]interface{}{"indent": 1, "list": "bullet"}, "insert": "\n"},
				{"insert": "level1-2"},
				{"attributes": map[string]interface{}{"list": "bullet"}, "insert": "\n"},
				{"insert": "not-list-2"},
				{"insert": "\n"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewSlackListGenerator()
			got, err := g.Generate(tt.htmlContent)
			if err != nil {
				t.Errorf("Generate() error = %v", err)
				return
			}

			// Normalize newlines for comparison
			gotPlainText := strings.TrimSpace(got.PlainText)
			wantPlainText := strings.TrimSpace(tt.wantPlainText)

			if gotPlainText != wantPlainText {
				t.Errorf("Generate() PlainText = %v, want %v", gotPlainText, wantPlainText)
			}

			gotOps := got.TextyJSON["ops"].([]map[string]interface{})
			if !reflect.DeepEqual(gotOps, tt.wantOps) {
				t.Errorf("Generate() Ops = %v, want %v", gotOps, tt.wantOps)
			}
		})
	}
}
