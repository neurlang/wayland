package main

var extendedaDescriptor = `Ā	Ġ	Ŀ	ş	ā	ġ	ŀ	Š
Ă	Ģ	Ł	š	ă	ģ	ł	Ţ
Ą	Ĥ	Ń	ţ	ą	ĥ	ń	Ť
Ć	Ħ	Ņ	ť	ć	ħ	ņ	Ŧ
Ĉ	Ĩ	Ň	ŧ	ĉ	ĩ	ň	Ũ
Ċ	Ī	ŉ	ũ	ċ	ī	Ŋ	Ū
Č	Ĭ	ŋ	ū	č	ĭ	Ō	Ŭ
Ď	Į	ō	ŭ	ď	į	Ŏ	Ů
Đ	İ	ŏ	ů	đ	ı	Ő	Ű
Ē	Ĳ	ő	ű	ē	ĳ	Œ	Ų
Ĕ	Ĵ	œ	ų	ĕ	ĵ	Ŕ	Ŵ
Ė	Ķ	ŕ	ŵ	ė	ķ	Ŗ	Ŷ
Ę	ſ	ŗ	ŷ	ę	ĸ	Ř	Ÿ
Ě	Ĺ	ř	Ź	ě	ĺ	Ś	ź
Ĝ	Ļ	ś	Ż	ĝ	ļ	Ŝ	ż
Ğ	Ľ	ŝ	Ž	ğ	ľ	Ş	ž`

var supplementDescriptor = `¡	±	À	Ï	à	ï
¢	²	Á	Ñ	á	ñ
£	³	Â	Ò	â	ò
¤	´	Ã	Ó	ã	ó
¥	µ	Ä	Ô	ä	ô
¦	¶	Å	Õ	å	õ
§	·	Æ	Ö	æ	ö
¨	¸	Ç	Ø	ç	ø
©	¹	Ð	Ù	ð	ù
ª	º	È	Ú	è	ú
«	»	É	Û	é	û
¬	¼	Ê	Ü	ê	ü
½	¾	Ë	Ý	ë	ý
­®	¿	Ì	Ÿ	ì	ÿ
¯	÷	Í	Þ	í	þ
°	×	Î	ẞ	î	ß`

var UnicodeFont Font

func init() {
	(&UnicodeFont).Load("ascii.png", asciiDescriptor)
	(&UnicodeFont).Load("extendeda.png", extendedaDescriptor)
	(&UnicodeFont).Load("supplement.png", supplementDescriptor)
}
