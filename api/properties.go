package api

type Properties struct {
	properties	map[string]interface{}
	layers  	[]*Source
	transformers	[]*Transformer
	initialized	bool
}

func NewProperties() *Properties {
	return &Properties{
		properties: make(map[string]interface{}),
		layers: make([]*Source, 0),
		transformers: make([]*Transformer, 0),
		initialized: false,
	}
}

func (p *Properties) Sources() {

}

func (p *Properties) Dynamic(source Source, namespace string) {
	p.layers = append(p.layers, source)
}

func (p *Properties) Static(properties map[string]interface{}, namespace string) {
	//source := make(Source)
	// Create source from property map

}

func (p *Properties) Initialize() {

}