package metadata

import "gerrit-share.lan/go/inspect"

const NoPath = -1

type Metadata struct {
	TypeId           inspect.TypeId
	TypeName         string
	TypeDescription  string
	Name             string
	Description      string
	IsMandatory      bool
	Default          string
	UnderlyingValues []*Metadata
	ValueTypeName    string
	PathIndex        int
}

func (m *Metadata) Add(data *Metadata) {
	m.UnderlyingValues = append(m.UnderlyingValues, data)
}

func (m *Metadata) SetTypeInfo(typeId inspect.TypeId, name string, description string, defaultValue string) {
	m.TypeId = typeId
	m.TypeName = name
	m.TypeDescription = description
	m.Default = defaultValue
}

func NewMetadata(typeId inspect.TypeId, value string, name string,
	mandatory bool, description string) *Metadata {
	return &Metadata{
		TypeId:           typeId,
		Name:             name,
		IsMandatory:      mandatory,
		Description:      description,
		Default:          value,
		UnderlyingValues: nil,
		PathIndex:        NoPath,
	}
}
