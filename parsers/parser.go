package parsers

type Parser interface {
	Parse(*map[string]string) map[string]interface{}
}