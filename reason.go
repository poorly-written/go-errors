package errors

type Reason interface {
	ToHashMap() map[string]interface{}
}

type reason struct {
	Type      string     `json:"type"`
	Info      *string    `json:"info"`
	Attribute *Attribute `json:"attribute"`
}

func (r reason) ToHashMap() map[string]interface{} {
	d := make(map[string]interface{})

	d["type"] = r.Type

	if r.Info != nil {
		d["info"] = *r.Info
	}

	if r.Attribute != nil {
		d["attribute"] = r.Attribute.toHashMap()
	}

	return d
}
