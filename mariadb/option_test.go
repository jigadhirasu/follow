package mariadb_test

import (
	"fmt"
	"regexp"

	"gorm.io/gorm/schema"
)

func ExampleRegexp() {
	regexpField, _ := regexp.Compile(`^\w+!?$`)

	fmt.Println(regexpField.MatchString("Abc123"))
	fmt.Println(regexpField.MatchString("123bBc_"))
	fmt.Println(regexpField.MatchString("AA_AA"))
	fmt.Println(regexpField.MatchString("Abc123!"))
	fmt.Println(regexpField.MatchString("123bBc_!"))
	fmt.Println(regexpField.MatchString("AA_AA!"))

	// Output:
	// true
	// true
	// true
	// true
	// true
	// true
}

func ExampleLookup() {
	f := schema.Schema{}.LookUpField("abc_t")

	fmt.Println(f.Name)
	fmt.Println(f.DBName)

	// Output:
	// AbcT
	// abc_t
}
