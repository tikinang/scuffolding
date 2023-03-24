package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Hello struct {
	app.Compo

	title string
}

func (r *Hello) OnNav(ctx app.Context) {
	fmt.Println("component navigated", ctx.Page().URL().String())
}

func (r *Hello) OnMount(ctx app.Context) {
	fmt.Println("component mounted")
}

func (r *Hello) OnPreRender(ctx app.Context) {
	fmt.Println("component pre-rendered")
}

func (r *Hello) Render() app.UI {
	return app.
		Div().
		Body(
			app.H1().Text("Hello World!"),
			app.Code().
				Style("color", "deepskyblue").
				Text("func() app.UI { well then }"),
		)
}

func main() {
	app.Route("/", &Hello{
		title: "That's my Jam!",
	})

	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "That's my Jam!",
		Lang:        "en",
		Description: "A simple web app to help you discover new music.",
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
