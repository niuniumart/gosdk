package tools

import (
	"fmt"
	"testing"
	"time"
)

//Man man接口
type Man interface {
	Test(*Ppl)
}

//Person struct
type Person struct {
	Name string
	Age  int
}

//Ppl struct
type Ppl struct {
	Name string
	Age  int
}

func (p *Person) Test(ppl *Ppl) {
	fmt.Printf("in%s\n", ppl.Name)
}

func TestRoutineRecover(t *testing.T) {
	var man Man

	man = &Person{
		Name: "222",
		Age:  2,
	}

	GoRoutine(man.Test, &Ppl{
		Name: "ab",
	})

	ret, err := FuncProxy(func(person *Person) *Ppl {
		fmt.Printf("www%s\n", person.Name)
		panic("pa")
		return &Ppl{
			Name: person.Name,
		}
	}, man)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	if ppl, ok := ret.(*Ppl); ok {
		fmt.Printf("%s\n", ppl.Name)
	}

	/*GoRoutine(func(person Person) {
		fmt.Printf("%s,%d\n", person.Name, person.Age)
		panic("panic er")
	}, &Person{
		Name: "alex",
		Age:  10,
	})*/
	time.Sleep(100 * time.Second)
}
