package darkplus

func MergeStyles() map[string]string {
	result := make(map[string]string)
	styleSets := []map[string]string{
		SharedStyles,
		JsStyles,
		BashStyles,
		PythonStyles,
		YamlStyles,
		TypeScriptStyles,
		CssStyles,
		SqlStyles,
		MarkdownStyles,
		CStyles,
		GoStyles,
		RustStyles,
		CppStyles,
		HtmlStyles,
		CsharpStyles,
		JavaStyles,
		JsonStyles,
		PhpStyles,
	}
	for _, set := range styleSets {
		for k, v := range set {
			result[k] = v
		}
	}
	return result
}
