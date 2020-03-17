package httpserver

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/metadata"
)

func MakeResultInfo(sampleValue interface{}) (*metadata.Metadata, error) {
	switch v := sampleValue.(type) {
	case nil:
		return &metadata.Metadata{}, nil
	case inspect.Inspectable:
		inspectorImpl := metadata.NewMetadataCreator()
		inspector := inspect.NewGenericInspector(inspectorImpl)
		v.Inspect(inspector)
		return inspectorImpl.Metadata, inspector.GetError()
	case string:
		return &metadata.Metadata{TypeId: inspect.TypeString}, nil
	case int:
		return &metadata.Metadata{TypeId: inspect.TypeInt}, nil
	case int32:
		return &metadata.Metadata{TypeId: inspect.TypeInt32}, nil
	case int64:
		return &metadata.Metadata{TypeId: inspect.TypeInt64}, nil
	case float32:
		return &metadata.Metadata{TypeId: inspect.TypeFloat32}, nil
	case float64:
		return &metadata.Metadata{TypeId: inspect.TypeFloat64}, nil
	}
	return nil, ErrUnsupportedResultType
}
