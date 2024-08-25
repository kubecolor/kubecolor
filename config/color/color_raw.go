package color

// Raw is a color of raw ANSI codes
type Raw string

// Code implements [ColorCode]
func (r Raw) Code() string {
	return string(r)
}
