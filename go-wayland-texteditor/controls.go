package main

var controls1Descriptor = "" +
	`Ctr0l	Ctr0r	Ctr1l	Ctr1r	Ctr2l	Ctr2r	 `

var controls2Descriptor = "" +
	`Ctr3l	Ctr3r	Ctr4l	Ctr4r	Ctr5l	Ctr5r	 `

var ControlFont Font

func init() {
	_ = (&ControlFont).Load("controls.png", controls1Descriptor, "")
	_ = (&ControlFont).Load("controls2.png", controls2Descriptor, "")
	_ = (&ControlFont).Alias("", " ")
}
