package container

type ThingDescription struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	AtContext   string `json:"@context"`
	AtType      string `json:"@type"`
	Description string `json:"description"`
	Base        string `json:"base"`
	BaseHref    string `json:"baseHref"`
	Href        string `json:"href"`

	Properties map[string]*Property `json:"properties"`
	Actions    map[string]*Action   `json:"actions"`
	Events     map[string]*Event    `json:"events"`

	CredentialsRequired bool   `json:"credentialsRequired"`
	FloorplanVisibility bool   `json:"floorplanVisibility"`
	FloorplanX          uint   `json:"floorplanX"`
	FloorplanY          uint   `json:"floorplanY"`
	LayoutIndex         uint   `json:"layoutIndex"`
	SelectedCapability  string `json:"selectedCapability"`
	IconHref            string `json:"iconHref"`

	Security string `json:"security"`

	GroupId string `json:"group_id"`
}

type Property struct {
	Name        string        `json:"name,omitempty"`
	AtType      string        `json:"@type,omitempty"`
	Title       string        `json:"title,omitempty"`
	Type        string        `json:"type"`
	Unit        string        `json:"unit,omitempty"`
	Description string        `json:"description,omitempty"`
	Minimum     interface{}   `json:"minimum,omitempty"`
	Maximum     interface{}   `json:"maximum,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
	ReadOnly    bool          `json:"readOnly,omitempty"`
	MultipleOf  float64       `json:"multipleOf,omitempty"`

	Value interface{} `json:"value,omitempty"`
}
type Action struct {
	Name        string `json:"name"`
	AtType      string `json:"@type,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	Input interface{} `json:"input"`
}

func (a Action) GetName() string {
	return a.Name
}

func (a Action) GetDescription() interface{} {
	return nil
}

type Event struct {
	AtType      string `json:"@type,omitempty"`
	Name        string `json:"name,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Type       string        `json:"type"`
	Unit       string        `json:"unit"`
	Minimum    interface{}   `json:"minimum"`
	Maximum    interface{}   `json:"maximum"`
	MultipleOf float64       `json:"multipleOf"`
	Enum       []interface{} `json:"enum"`
}

func (e Event) GetName() string {
	return e.Name
}

type PIN struct {
	Required bool   `json:"required,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
}

func (p Property) GetName() string {
	return p.Name
}

func (p Property) GetValue() interface{} {
	return p.Value
}
