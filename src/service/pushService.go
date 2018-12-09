package service

type Glasses struct {
	Price int64
	From  string
}

type Builder interface {
	BuildPrice(int64)
	BuildForm()
	Build() Glasses
}

type ImageBuilder struct {
	glasses Glasses
}

type TextBuilder struct {
	glasses Glasses
}

func (pS *ImageBuilder) BuildPrice(iP int64) {
	pS.glasses.Price = iP * 10
}

func (pS *ImageBuilder) BuildForm() {
	pS.glasses.From = "shenzhen"
}

func (pS *ImageBuilder) Build() Glasses {
	return pS.glasses
}
func main() {

}
