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

// NewReason allows to create new reason
//
// rType - reason type. e.g. "invalid", "exists", "required"
// info - additional information about the reason. e.g. "the age field is required", "the phone already exists"
// attribute - additional attributes for the reason. "min"/"max" values, "format"
func NewReason(rType string, info *string, attribute *Attribute) Reason {
	return reason{
		Info:      info,
		Type:      rType,
		Attribute: attribute,
	}
}

// SimpleReason creates a reason with a message
func SimpleReason(reason string) Reason {
	return NewReason(reason, nil, nil)
}
