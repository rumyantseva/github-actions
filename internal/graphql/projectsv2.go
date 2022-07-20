package graphql

import (
	"context"
	"encoding/json"
	"fmt"
)

// query is used to get PR's project items.3
var query = `{
	node(id: "PR_kwDOHbB198459Yt9") {
		... on PullRequest {
			title
			state
			id
			projectItems(first: 20) {
				totalCount
				nodes {
					fieldValues(first: 20) {
						totalCount
						nodes {
							... on ProjectV2ItemFieldRepositoryValue {
								__typename
								field {
									... on ProjectV2Field {
										name
									}
								}

								repository {
									__typename
									name
								}
							}

							... on ProjectV2ItemFieldIterationValue {
								__typename
								title
								field {
									... on ProjectV2IterationField {
										name
									}
								}
							}

							... on ProjectV2ItemFieldMilestoneValue {
								__typename
								field {
									... on ProjectV2Field {
										name
									}
								}
								milestone {
									title
								}
							}

							... on ProjectV2ItemFieldTextValue {
								__typename
								field {
									... on ProjectV2Field {
										name
									}
								}
								text
							}

							... on ProjectV2ItemFieldSingleSelectValue {
								__typename
								field {
									... on ProjectV2SingleSelectField {
										name
									}
								}
								name
							}
						}
					}
				}
			}
		}
	}
}
`

func (c *Client) GetPRProjectItems(ctx context.Context, nodeID string) error {
	//out, err := c.do(ctx, query, map[string]any{"nodeID": nodeID})
	out, err := c.do(ctx, query, map[string]any{})

	if err != nil {
		return err
	}

	var a any
	json.Unmarshal(out.Data, &a)

	var e any
	json.Unmarshal(out.Errors, &e)

	fmt.Println(a)

	return nil
}
