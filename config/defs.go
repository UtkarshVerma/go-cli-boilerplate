package config
type (
config struct {
Gopher bool `json:"gopher,omitempty"`
File string `json:"-" name:"config",`
Greta struct {
Name string `json:"name,omitempty"`
Unexported string `json:"-"`
CustomVar bool `json:"-" name:"custom",`
} `json:"greet,omitempty" name:"greet"`
}
)
