package velvet

import "html/template"

func ifHelper(conditional bool, help HelperContext) (template.HTML, error) {
	if conditional {
		s, err := help.Block()
		return template.HTML(s), err
	}
	s, err := help.ElseBlock()
	return template.HTML(s), err
}
