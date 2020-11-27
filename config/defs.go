package config
const configFlag = "config"
type (
config struct {
File string`json:"-"`
Gopher string `json:"gopher,omitempty"`
Greet struct {
Name string `json:"name,omitempty"`
} `json:"greet,omitempty"`
}
)
