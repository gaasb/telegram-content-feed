package bot

var (
	categories map[string]string
)

type Category struct {
	Period interface{}
}

func (c *Category) Append(name string) {

}
func (c *Category) WithPeriod() {
	c.Period = "week"
}
