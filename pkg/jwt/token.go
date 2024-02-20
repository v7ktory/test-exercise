package jwt

type JWT struct {
}

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) Generate() (string, error) {

}
