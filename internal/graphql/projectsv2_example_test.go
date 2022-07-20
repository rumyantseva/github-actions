package graphql

import (
	"context"

	"golang.org/x/oauth2"
)

func Example() {
	token := "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"
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), ts)
	qlClient := NewClient(httpClient)

	err := qlClient.GetPRProjectItems(context.Background(), "PR_kwDOHbB198459Yt9")
	if err != nil {
		panic(err)
	}

	// Output:
	// [Hello]
}
