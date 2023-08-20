package pizzahut

import (
	"fmt"

	"github.com/maxthom/artifact-store/store"
)

var ()

func init() {

}

type PizzaHutStorePlan struct {
	Address   string    `yaml:"address"`
	City      string    `yaml:"city"`
	Employees Employees `yaml:"employees"`
	Menu      []Menu    `yaml:"menu"`
}
type Employees struct {
	Managers int `yaml:"managers"`
	Drivers  int `yaml:"drivers"`
	Cooks    int `yaml:"cooks"`
	Counter  int `yaml:"counter"`
}
type Menu struct {
	Name  string  `yaml:"name"`
	Price float64 `yaml:"price"`
}

func OpenStore() {
	fmt.Println("OPEN!")

	files := []string{"pizzahut.yaml"}
	fmt.Println(store.ListFileContentToType[PizzaHutStorePlan]("pizzahut", "default", "v1.0.0", files...))

}

func CloseStore() {
	fmt.Println("CLOSE!")
}
