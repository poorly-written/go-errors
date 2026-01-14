package errors

type Attribute struct {
	Min     interface{}   `json:"min"`
	Max     interface{}   `json:"max"`
	In      interface{}   `json:"in"`
	Pattern *string       `json:"pattern"`
	Format  *string       `json:"format"`
	Value   interface{}   `json:"value"`
	Values  []interface{} `json:"values"`
	Field   *string       `json:"field"`
	Fields  []any         `json:"fields"`
}

func (attr Attribute) toHashMap() map[string]interface{} {
	d := make(map[string]interface{})

	if attr.Min != nil {
		d["min"] = attr.Min
	}

	if attr.Max != nil {
		d["max"] = attr.Max
	}

	if attr.In != nil {
		d["in"] = attr.In
	}

	if attr.Pattern != nil {
		d["pattern"] = *attr.Pattern
	}

	if attr.Format != nil {
		d["format"] = *attr.Format
	}

	if attr.Value != nil {
		d["value"] = attr.Value
	}

	if attr.Values != nil {
		d["values"] = attr.Values
	}

	if attr.Field != nil {
		d["field"] = *attr.Field
	}

	if attr.Fields != nil {
		d["fields"] = attr.Fields
	}

	if len(d) == 0 {
		return nil
	}

	return d
}
