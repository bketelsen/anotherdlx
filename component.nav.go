package main

import "html/template"

type ComponentNav struct {
	Title string
	Items []ComponentNavbarHref
}

type ComponentNavbarHref struct {
	Title     template.HTML
	Image     template.HTML
	Href      string
	AriaLabel string
	Current   bool
}
