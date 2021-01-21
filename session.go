package jsessions

type ISessions interface {
	Get(id string) ISession
}

type ISession interface {
	ID() string
	Set(name string, newValue interface{}) (savedValue interface{}, err error)
	SetString(name string, newValue string) (savedValue string, err error)
	SetInt(name string, newValue int) (savedValue int, err error)
	SetBool(name string, newValue bool) (savedValue bool, err error)

	Get(name string) (value interface{}, err error)
	GetString(name string) (string, bool)
	GetInt(name string) (int, bool)
	GetBool(name string) (bool, bool)

	Data() map[string]interface{}

	Save() error
	Close()
}
