package velvet

import "html/template"

func unlessHelper(conditional bool, help HelperContext) (template.HTML, error) {
	return ifHelper(!conditional, help)
}
