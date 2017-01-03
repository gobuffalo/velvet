package velvet

import "html/template"

func equalHelper(a, b interface{}, help HelperContext) (template.HTML, error) {
	if a == b {
		s, err := help.Block()
		if err != nil {
			return "", err
		}
		return template.HTML(s), nil
	}
	s, err := help.ElseBlock()
	if err != nil {
		return "", err
	}
	return template.HTML(s), nil
}
