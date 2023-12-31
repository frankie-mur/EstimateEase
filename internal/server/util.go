package server

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

// RenderComponent renders the provided component to the given bytes.Buffer
// using the context. It returns the rendered content as a string and any error.
func RenderComponent(component templ.Component, ctx context.Context, buf *bytes.Buffer) (string, error) {
	err := component.Render(ctx, buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderComponentToString renders a templ component and returns as string
func RenderComponentToString(component templ.Component, ctx context.Context) (string, error) {
	var buf bytes.Buffer
	data, err := RenderComponent(
		component,
		ctx,
		&buf,
	)
	if err != nil {
		return "", err
	}

	return data, nil
}
