package config
type (
config struct {
File string `json:"-" name:"config",`
Gopher bool `json:"gopher,omitempty"`
Greta struct {
Name string `json:"name,omitempty"`
Unexported string `json:"-"`
CustomVar bool `json:"-" name:"custom",`
} `json:"greet,omitempty" name:"greet"`
}
)
