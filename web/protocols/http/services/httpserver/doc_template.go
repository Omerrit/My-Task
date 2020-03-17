package httpserver

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/utils/maps"
	"gerrit-share.lan/go/web/protocols/http/services/httpserver/internal/metadata"
	"html/template"
	"io"
)

const (
	docTemplateName     = "documentation"
	cssTemplateName     = "styles"
	linkTemplateName    = "links"
	cssLinkTemplateName = "links_styles"
)

const cssTemplate = `
	<style>
		h2 {margin-top: 0px;}
		h3 {margin-bottom: 8px;}
		h4 {margin: 0;}
		li > h3 {cursor: pointer; }
		li > h3:hover {text-decoration: underline;}
		.endpoint-info {display:none;}
		body {background: #DEDEDE; padding: 10px 20px;}
		.parameter {padding-top: 12px;}
		.parameters {background: #FFF; padding: 14px; margin-right: 10px;}
		.red {color: #de495d;}
		li {min-width: 1200px;}
		li > .description {margin-bottom: 8px;}
		.description {font-style: italic;}
		.type {font-weight: 700; font-style: italic;}
		.type, .param-name, .value {width: 140px; display: inline-block;}
		.type {width: 80px;}
		.value {width: 240px;}
		.divider {border-bottom: 1px #DEDEDE solid; width: 480px;}
	</style>
`

const docTemplate = `
		{{ template "result-style" }}
		<div>
			<script>
				function toggleInfo(e) {
					let elem = e.parentNode.getElementsByClassName('endpoint-info')[0]
					let ellipsis = e.parentNode.getElementsByClassName('ellipsis')[0]
					if (ellipsis && e !== ellipsis.parentNode) {
						ellipsis = undefined
					}
					let style = window.getComputedStyle(elem)
					if (style.display !== 'none') {
						elem.style.display = 'none'
						if (ellipsis) {
							ellipsis.style.display = 'inline-block'
						}
					} else {
						elem.style.display = 'block'
						if (ellipsis) {
							ellipsis.style.display = 'none'
						}
					}
					window.event.stopPropagation()
				}
			</script>
			{{ range $resourceName, $resource := . }}
					<h2>Resource: {{ $resourceName }}</h2>
					<div>
						<h3>Methods: </h3>
						<ul>
							{{ range $endpointName, $endpoint := $resource }}
								<li>
									<h3 onclick="toggleInfo(this)">{{ $endpointName }}</h3>
									<div class="description">{{ or ($endpoint.CommandMetaData.Description) ($endpoint.CommandMetaData.TypeDescription) }}</div>
									<div class="endpoint-info">
										{{ if $endpoint.CommandMetaData.UnderlyingValues }}
											<div class="parameters">
												<h4>Params</h4>
												{{ range $parameterName, $info := $endpoint.CommandMetaData.UnderlyingValues }}
													<div class="parameter">
														<div>
															<span class="param-name">{{ $info.Name }}</span>
															<span class="type"> {{getTypeName $info.TypeId}}</span>
															{{ if $info.IsMandatory }}
																<span class="red value">[mandatory]</span>
															{{ end }}
															{{ if not $info.IsMandatory }}
																<span class="value">	
																	(default: 
																	{{ if eq $info.Default "" }}
																		" "
																	{{ else }}
																		{{$info.Default}}
																	{{ end }}
																	)
																</span>
															{{ end }}
														</div>
														<div class="description">{{ $info.Description }}</div>
													</div>
												{{ end }}
											</div>
											<div class="divider"></div>
										{{ end }}
										{{ if not (eq (getTypeName $endpoint.ResultInfo.TypeId) "") }}
											<div class="parameters">
												<h4>Response</h4>
												<div class="parameter">
													<div>
														{{ if $endpoint.ResultInfo}}
															{{ template "result-info" $endpoint.ResultInfo }}
														{{ end }}
														{{ if not $endpoint.ResultInfo}}
															<span class="type"> {{getTypeName $endpoint.ResultInfo.TypeId}}</span>
														{{ end }}
													</div>
												</div>
											</div>
										{{ end }}
									</div>
								</li>
							{{ end }}
						</ul>
					</div>
			{{ end }}
		</div>
	`

const linkTemplate = `
	<div class="links">
		<h3>Sub Resources:</h3>
		<ul>
			{{ range $name, $link := . }}
				<li><a href={{ $link }}>{{ $name }}</a></li>
			{{ end }}
		</ul>
	</div>
`

const cssLinkTemplate = `
	<style>
		.links {display: flex; flex-direction: column;}
		body {background: #DEDEDE; padding: 10px 20px;}
	</style>
`

const resultTemplate = `
	{{ define "result-info" }}
		<div class="parameter">
			<div>
				{{ if .Name }}
					<span class="param-name">{{ .Name }}</span>
				{{ end }}
				{{ if .TypeName }}
					<span title={{.TypeName}} class="type dotted">{{ getTypeName .TypeId }}</span>
				{{ else }}
					<span class="type">{{ getTypeName .TypeId }}</span>
				{{ end }}
			</div>
			<div class="description">{{ or (.Description) (.TypeDescription) }}</div>
			<div class="clicker-wrapper" onclick="toggleInfo(this)">
				{{ if .UnderlyingValues }}
					{{ if isObject .TypeId }}
						<span class="braces">{</span>
						<span class="ellipsis">...</span>
					{{ end }}
					{{ if isArray .TypeId }}
						<span class="braces">[</span>
						<span class="ellipsis">...</span>
					{{ end }}
					<div class="underlying endpoint-info">
						{{ range $i, $info := .UnderlyingValues }}
							{{ template "result-info" $info }}
						{{ end }}
					</div>
					{{ if isObject .TypeId }}
						<span class="braces">}</span>
					{{ end }}
					{{ if isArray .TypeId }}
						<span class="braces">]</span>
					{{ end }}
				{{ end }}	
			</div>
		</div>
	{{ end }}
`

const cssResultTemplate = `
	{{ define "result-style" }}
		<style>
			.underlying {margin-left: 80px;}
			.braces {margin: 0 40px; margin-top: 20px;}
			.clicker-wrapper {cursor: pointer;}
			.dotted {text-decoration:underline; color:#72add4;}
		</style>
	{{ end }}
`

func executeMethodsTemplate(writer io.Writer, data map[string]map[string]endpointInfo) error {
	cssT, err := template.New(cssTemplateName).Parse(cssTemplate)
	if err != nil {
		return err
	}
	err = cssT.Execute(writer, nil)
	if err != nil {
		return err
	}
	t := template.New(docTemplateName)

	t = t.Funcs(template.FuncMap{
		"getTypeName": metadata.GetTypeName,
		"isArray":     isArrayLike,
		"isObject":    isObjectLike,
	})
	t, err = t.Parse(docTemplate)
	if err != nil {
		return err
	}
	t, err = t.Parse(resultTemplate)
	if err != nil {
		return err
	}
	t, err = t.Parse(cssResultTemplate)
	if err != nil {
		return err
	}
	return t.Execute(writer, data)
}

func executeLinksTemplate(writer io.Writer, data maps.String) error {
	cssT, err := template.New(cssLinkTemplateName).Parse(cssLinkTemplate)
	if err != nil {
		return err
	}
	err = cssT.Execute(writer, nil)
	if err != nil {
		return err
	}
	t, err := template.New(linkTemplateName).Parse(linkTemplate)
	if err != nil {
		return err
	}
	return t.Execute(writer, data)
}

func isObjectLike(typeId inspect.TypeId) bool {
	if typeId == inspect.TypeObject || typeId == inspect.TypeMap || typeId == inspect.TypeValue {
		return true
	}
	return false
}

func isArrayLike(typeId inspect.TypeId) bool {
	if typeId == inspect.TypeArray {
		return true
	}
	return false
}
