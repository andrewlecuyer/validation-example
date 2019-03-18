package main

import (
	"fmt"
	"os"

	"github.com/andrewlecuyer/validation-example/pgovalidate"
	"github.com/spf13/pflag"
	validator "gopkg.in/go-playground/validator.v9"
)

func main() {

	var backupOpts string

	myCommandLine := pflag.NewFlagSet("mycommand", pflag.ExitOnError)
	myCommandLine.StringVar(&backupOpts, "backup-opts", "", "")
	myCommandLine.Parse(os.Args[1:])

	parsedOpts := pgovalidate.ParseBackupOpts(backupOpts)

	backupOptions, setFlagFieldNames := pgovalidate.ConvertBackupOptsToStruct(parsedOpts)

	err := pgovalidate.ValidateBackupOpts(backupOptions, setFlagFieldNames)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}
	} else {
		fmt.Println("flags are valid!")
	}
}
