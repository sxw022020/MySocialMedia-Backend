package model

/*
*
`struct`:
  - a composite data type ("blueprint") that groups together zero or more named values of different types under a single name
  - this is similar to an `object` in Java

`xxx`: raw string literals:
  - can include any character that can appear in a text string, including control characters and special characters.
*/
type Post struct {
	Id      string `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	Url     string `json:"url"`
	Type    string `json:"type"`
}
