// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graphql

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

// GraphQLFields represent a list of GitHub PNFs (Project Next Field).
type GraphQLFields struct {
	TotalCount githubv4.Int
	Nodes      []GraphQLField
}

// GraphQLField represents a single GitHub PNF (Project Next Field).
type GraphQLField struct {
	ID       githubv4.ID
	Name     githubv4.String
	DataType githubv4.ProjectNextFieldType
	Settings githubv4.String
}

// GraphQLFieldValues represents list of GitHub PNFVs (Project Next Field Value).
type GraphQLFieldValues struct {
	TotalCount githubv4.Int
	Nodes      []GraphQLFieldValue
}

// GraphQLFieldValue represents a single GitHub PNFV (Project Next Field Value).
type GraphQLFieldValue struct {
	ID           githubv4.ID
	Value        githubv4.String
	ProjectField GraphQLField

	// ValueTitle is a special field to display the value in a more human-readable way.
	ValueTitle string `graphql:"value"`
}

// Items represents a list of GitHub project items.
// Data -> Node -> ProjectItems
type Items struct {
	TotalCount githubv4.Int
	Nodes      []ProjectV2FieldValues `graphql:"fieldValues(first: $fieldValuesMax)"`
}

// ProjectV2FieldValues represents an item within project.
// Data -> Node -> ProjectItems -> Nodes
type ProjectV2FieldValues struct {
	FieldValues any
}

// Data -> Node -> ProjectItems -> Nodes -> []
type ProjectItemsNodes struct {
	TotalCount  githubv4.Int
	FieldValues any
}

// ProjectV2ItemFieldValueConnection represents the connection type for ProjectV2ItemFieldValue.
type ProjectV2ItemFieldValueConnection struct {
	TypeName githubv4.String
	Query    string `graphql:"... on ProjectV2ItemFieldIterationValue"`
	//GenericField GenericField
	ID githubv4.ID
}

type GenericField struct {
	ID githubv4.ID
}

// PRItem represents Pull Request item values.
type PRItem struct {
	FieldName string
	Value     string
}

// Querier describes a GitHub GraphQL client that can make a query.
type Querier interface {
	// Query executes the given GraphQL query `q` with the given variables `vars` and stores the results in `q`.
	Query(ctx context.Context, q any, vars map[string]any) error
}

// GetPRItems returns the list of PNIs - Project Next Items (cards) associated with the given PR.
func GetPRItems(client Querier, nodeID string) ([]PRItem, error) {
	var q struct {
		Data struct {
			// data -> node
			Node struct {
				ID    githubv4.String
				Title githubv4.String
				State githubv4.String

				// data -> node -> projectItems
				ProjectItems Items `graphql:"projectItems(first: $itemsMax)"`
			} `graphql:"... on PullRequest"`
		} `graphql:"node(id: $nodeID)"`
	}

	variables := map[string]any{
		"nodeID":         githubv4.ID(nodeID),
		"itemsMax":       githubv4.Int(20),
		"fieldValuesMax": githubv4.Int(10),
	}

	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}

	fmt.Println(q.Data.Node.ProjectItems.Nodes[0].FieldValues)

	//if q.Node.PullRequest.ProjectItems.TotalCount == 0 {
	//	return nil, nil
	//}

	var result []PRItem

	//	for _, v := range q.Node.PullRequest.ProjectItems.Nodes {
	////		for _, vv := range v.Nodes {
	//			typename := vv.TypeName
	//			switch typename {
	/*case "ProjectV2ItemFieldIterationValue":
		title, ok := v["title"]
		if !ok {
			continue
		}

		result = append(result, PRItem{
			FieldName: getFieldName(v),
			Value:     title.(string),
		})
	case "ProjectV2ItemFieldMilestoneValue":
		milestone, ok := v["milestone"]
		if !ok {
			continue
		}
		title, ok := milestone.(map[string]any)["title"]
		if !ok {
			continue
		}

		result = append(result, PRItem{
			FieldName: getFieldName(v),
			Value:     title.(string),
		})
	case "ProjectV2ItemFieldSingleSelectValue":
		name, ok := v["name"]
		if !ok {
			continue
		}

		result = append(result, PRItem{
			FieldName: getFieldName(v),
			Value:     name.(string),
		})
	*/
	//			}
	//		}
	//	}

	return result, nil
}

// getFieldName returns  PR object item field name.
func getFieldName(v map[string]any) string {
	field, ok := v["field"]
	if !ok {
		return ""
	}
	fieldName, ok := field.(map[string]any)["name"]
	if !ok {
		return ""
	}

	return fieldName.(string)
}
