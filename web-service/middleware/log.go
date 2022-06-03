package middleware

import (
	"fmt"
	"log"
)

func LogRun(source string) {
	log.Println(fmt.Sprintf("[REQUEST] %v", source))
}

func LogError(source string, err error) {
	log.Println(fmt.Sprintf("[ERROR] %v | %v", source, err.Error()))
}
