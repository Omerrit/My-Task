package imports

import (
	"gerrit-share.lan/go/tools/sourcecheck/internal/filter"
	"gerrit-share.lan/go/utils/maps"
	"gerrit-share.lan/go/utils/sets"
)

type packageSettings struct {
	DirName        string
	Deny           filter.FilterSet
	AllowedImports sets.String
}

func (settings *packageSettings) Update(other packageSettings) {
	settings.Deny.Append(other.Deny)
	settings.AllowedImports = sets.Join(settings.AllowedImports, other.AllowedImports)
}

type dirSettings struct {
	RootInternalDir string
	ImportPath      string
	Variables       maps.String
}
