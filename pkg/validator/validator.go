package validator

type ValidationErrors map[string][]string

func (v ValidationErrors) Error() string {
	return "validation error"
}

type ValidatorHandler func(fieldname string) error
type Validator struct {
	vErrors ValidationErrors
}

func New() *Validator {
	return &Validator{
		vErrors: make(ValidationErrors),
	}
}

type Entity interface {
	Validations() map[string][]ValidatorHandler
}

func (v *Validator) Validate(entity Entity) ValidationErrors {
	vList := entity.Validations()
	for name, vals := range vList {
		for _, validator := range vals {
			err := validator(name)
			if err != nil {
				if _, ok := v.vErrors[name]; ok {
					v.vErrors[name] = append(v.vErrors[name], err.Error())
				} else {
					v.vErrors[name] = []string{err.Error()}
				}
			}
		}

	}
	if len(v.vErrors) == 0 {
		return nil
	}
	return v.vErrors
}
