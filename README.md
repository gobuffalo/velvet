# Velvet [![GoDoc](https://godoc.org/github.com/gobuffalo/velvet?status.svg)](https://godoc.org/github.com/gobuffalo/velvet) [![Build Status](https://travis-ci.org/gobuffalo/velvet.svg?branch=master)](https://travis-ci.org/gobuffalo/velvet) [![Code Climate](https://codeclimate.com/github/gobuffalo/velvet/badges/gpa.svg)](https://codeclimate.com/github/gobuffalo/velvet)

Velvet is a templating package for Go. It bears a striking resemblance to "handlebars" based templates, there are a few small changes/tweaks, that make it slightly different.

## General Usage

```handlebars
<!-- some input -->
<h1>{{ name }}</h1>
<ul>
  {{#each names}}
    <li>{{ this }}</li>
  {{/each}}
</ul>
```

```go
ctx := velvet.NewContext()
ctx.Set("name", "Mark")
ctx.Set("names", []string{"John", "Paul", "George", "Ringo"})
s, _ := velvet.Render(input, ctx)
```

```html
<h1>Mark</h1>
<ul>
    <li>John</li>
    <li>Paul</li>
    <li>George</li>
    <li>Ringo</li>
</ul>
```

## If statements

```handlebars
{{#if true }}
  render this
{{/if}}
```

### Else statements

```handlebars
{{#if false }}
  won't render this
{{ else }}
  render this
{{/if}}
```
