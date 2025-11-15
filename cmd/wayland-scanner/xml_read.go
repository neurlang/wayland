package main

import (
	"log"
	"os"
)

func xmlProcess(inputFile, outputFile string) {
	var data, err = os.ReadFile(inputFile)
	if err != nil {
		log.Fatalln("Cannot read input file:", err)
	}

	xml, err := UnmarshalXML(data)
	if err != nil {
		log.Fatalln("Cannot parse input file:", err)
	}
	//spew.Dump(xml)

	var f GoFile

	f.PkgName = before_wl(xml.Name)
	f.Name = sanitizeSingleLineComment(xml.Name)
	f.FileName = sanitizeSingleLineComment(inputFile)

	for _, iface := range xml.Interface {
		var s GoStruct
		s.Name = removePrefixAndCamelCase(iface.Name, f.PkgName)

		if iface.Description != nil {
			s.Comment = sanitizeSingleLineComment(iface.Description.Summary)
		}

		f.AddStruct(&s)

		for _, req := range iface.Request {
			var m GoMethod
			m.Name = removePrefixAndCamelCase(req.Name, f.PkgName)

			if req.Description != nil {
				m.Comment = sanitizeSingleLineComment(req.Description.Summary)
			}

			s.AddMethod(&m)

			for _, arg := range req.Arg {

				var a GoArg

				var t GoType

				a.Name = removePrefixAndCamelCase(arg.Name, f.PkgName)

				a.Comment = sanitizeSingleLineComment(arg.Summary)

				t.Type = arg.Type

				if arg.Interface != nil {
					t.Interface = removePrefixAndCamelCase(*(arg.Interface), f.PkgName)
				}
				var isMagicalXdgGetSurface = iface.Name == "xdg_wm_base" &&
					req.Name == "get_xdg_surface" &&
					arg.Name == "surface" &&
					arg.Type == "object"

				if isMagicalXdgGetSurface {
					t.Interface = "Wl" + t.Interface
				}
				var isMagicalRegistryBind = iface.Name == "wl_registry" &&
					req.Name == "bind" &&
					arg.Name == "id" &&
					arg.Type == "new_id"

				if isMagicalRegistryBind {
					// add two special parameters to registry bind

					m.AddEventArg(&GoArg{
						Name:    "Iface",
						Comment: "magical parameter iface",
					}, &GoType{
						Type: "string",
					})

					m.AddEventArg(&GoArg{
						Name:    "Version",
						Comment: "magical parameter version",
					}, &GoType{
						Type: "uint",
					})
				}

				m.AddEventArg(&a, &t)
			}
		}

		for _, event := range iface.Event {
			var e GoEvent
			e.Name = removePrefixAndCamelCase(iface.Name+"_"+event.Name, f.PkgName)

			if event.Description != nil {
				e.Comment = sanitizeSingleLineComment(event.Description.Summary)
			}

			var isMagicalCalbackDoneEvent = iface.Name == "wl_callback" && event.Name == "done"

			e.IsCb = isMagicalCalbackDoneEvent

			var isMagicalPointerEvent = iface.Name == "wl_pointer" && (event.Name == "button" || event.Name == "motion")

			e.IsPt = isMagicalPointerEvent

			var isMagicalBufferReleaseEvent = iface.Name == "wl_buffer" && event.Name == "release"

			e.IsBf = isMagicalBufferReleaseEvent

			s.AddEvent(&e)

			for _, arg := range event.Arg {

				var a GoArg

				var t GoType

				a.Name = removePrefixAndCamelCase(arg.Name, f.PkgName)

				a.Comment = sanitizeSingleLineComment(arg.Summary)

				t.Type = arg.Type

				if arg.Interface != nil {
					t.Interface = removePrefixAndCamelCase(*(arg.Interface), f.PkgName)
				}

				e.AddEventArg(&a, &t)
			}

			var h GoHandler

			h.Name = removePrefixAndCamelCase(event.Name, f.PkgName)

			s.AddHandler(&h)
		}

		for _, enum := range iface.Enum {

			for _, entry := range enum.Entry {

				var c GoConstant
				c.Name = removePrefixAndCamelCase(iface.Name+"_"+enum.Name+"_"+entry.Name, f.PkgName)

				c.Value = entry.Value

				if entry.Summary != nil {
					c.Comment = sanitizeSingleLineComment(*entry.Summary)
				}

				f.AddConstant(&c)

			}
		}
	}
	var data2 = f.Serialize()
	err = os.WriteFile(outputFile, []byte(data2), 0666)
	if err != nil {
		log.Fatalln("Cannot write output file:", err)
	}
}
