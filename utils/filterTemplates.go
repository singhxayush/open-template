package utils

import "strings"

// FilterTemplates returns a list of templates whose names contain the query (case-insensitive).
func FilterTemplates(templates []string, query string) []string {
	if query == "" {
		return templates
	}
	var result []string
	lowerQuery := strings.ToLower(query)
	for _, tmpl := range templates {
		if strings.Contains(strings.ToLower(tmpl), lowerQuery) {
			result = append(result, tmpl)
		}
	}
	return result
}
