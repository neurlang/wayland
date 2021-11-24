package main

var controlsDescriptor = "" +
	`Ctr0l	Ctr0r	Ctr1l	Ctr1r	Ctr2l	Ctr2r	 `

var ControlFont Font

func init() {
	(&ControlFont).Load("controls.png", controlsDescriptor, "")
	(&ControlFont).Alias("", " ")
}
