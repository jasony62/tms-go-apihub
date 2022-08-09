package main

type ApiHubHttpConf struct {
	ID                 string `json:"id"`
	Description        string `json:"description"`
	URL                string `json:"url"`
	Method             string `json:"method"`
	Private            string `json:"private"`
	Requestcontenttype string `json:"requestContentType"`
	Args               []Args `json:"args"`
}

type Value struct {
	From    string `json:"from"`
	Content string `json:"content"`
}
type Args struct {
	In    string `json:"in"`
	Name  string `json:"name"`
	Value Value  `json:"value"`
}
