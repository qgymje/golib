package provider

import "io"

type IView interface {
	Render(w io.Writer, temps ...string)
}
