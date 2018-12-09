package main

import (
	"fmt"
)

type Glasses struct {
	Price int64
	From  string
}

type Builder interface {
	BuildPrice(int64)
	BuildForm()
	GetGlasses() Glasses
}

type ShenZhenBuilder struct {
	glasses Glasses
}

type ShanWeiBuilder struct {
	glasses Glasses
}

func (pS *ShenZhenBuilder) BuildPrice(iP int64) {
	pS.glasses.Price = iP * 10
}

func (pS *ShenZhenBuilder) BuildForm() {
	pS.glasses.From = "shenzhen"
}

func (pS *ShenZhenBuilder) GetGlasses() Glasses {
	return pS.glasses
}

func (pS *ShanWeiBuilder) BuildPrice(iP int64) {
	pS.glasses.Price = iP * 2
}

func (pS *ShanWeiBuilder) BuildForm() {
	pS.glasses.From = "shanwei"
}

func (pS *ShanWeiBuilder) GetGlasses() Glasses {
	return pS.glasses
}

type LeshiGlasses struct {
	First_cost int64
}

func (L *LeshiGlasses) GetGlasses(builder Builder) Glasses {
	builder.BuildPrice(L.First_cost)
	builder.BuildForm()
	return builder.GetGlasses()
}

func main() {

	leshi := &LeshiGlasses{First_cost: 100}
	var glassesbuilder Builder
	glassesbuilder = &ShanWeiBuilder{}
	glasses := leshi.GetGlasses(glassesbuilder)
	fmt.Println("glasses's price is: ", glasses.Price, " glasses from :", glasses.From)
	glassesbuilder = &ShenZhenBuilder{}
	glasses = leshi.GetGlasses(glassesbuilder)
	fmt.Println("glasses's price is: ", glasses.Price, " glasses from :", glasses.From)

	return
}