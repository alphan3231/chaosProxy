package graphql

import (
	"encoding/json"
	"strings"
)

type Request struct {
	Query         string          `json:"query"`
	OperationName string          `json:"operationName"`
	Variables     json.RawMessage `json:"variables"`
}

type Operation struct {
	Name  string
	Query string
	Type  string // query, mutation, or subscription
}

// ParseRequest parses a JSON body to extract GraphQL operation details.
func ParseRequest(body []byte) (*Operation, error) {
	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	op := &Operation{
		Name:  req.OperationName,
		Query: req.Query,
	}

	// Simple heuristic to determine type if operation name is missing or just to be sure
	queryLower := strings.ToLower(req.Query)
	if strings.HasPrefix(strings.TrimSpace(queryLower), "mutation") {
		op.Type = "mutation"
	} else if strings.HasPrefix(strings.TrimSpace(queryLower), "subscription") {
		op.Type = "subscription"
	} else {
		op.Type = "query"
	}

	// If Name is empty, try to extract it from the query "query MyOp { ... }"
	if op.Name == "" {
		op.Name = extractOperationName(req.Query)
	}

	return op, nil
}

func extractOperationName(query string) string {
	// A very basic parser to find the second word after query/mutation/...
	parts := strings.Fields(query)
	if len(parts) > 1 {
		// e.g. "query GetUser {" -> GetUser
		// But verify it's not a brace
		if parts[1] != "{" && parts[1] != "(" {
			return parts[1]
		}
	}
	return "Anonymous"
}
