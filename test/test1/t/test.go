package main

import (
	"fmt"
	"strings"
)

func main(){
	fmt.Println(strings.NewReplacer("s","z").Replace("wwwwws"))
}


