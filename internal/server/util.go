package server

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

// renderComponentToBuffer renders the provided component to the given bytes.Buffer
// using the context. It returns the rendered content as a string and any error.
func RenderComponentToBuffer(component templ.Component, ctx context.Context, buf *bytes.Buffer) (string, error) {
	err := component.Render(ctx, buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
