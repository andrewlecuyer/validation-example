package pgovalidate

import (
	"reflect"
	"strings"

	"github.com/spf13/pflag"
)

var setFlagFieldNames []string

func ParseBackupOpts(backupOpts string) []string {

	newFields := []string{}
	var newField string
	for i, c := range backupOpts {
		// if another option is found, add current option to newFields array
		if !(c == ' ' && backupOpts[i+1] == '-') {
			newField = newField + string(c)
		}

		// append if at the end of the flag (i.e. if another new flag was found) or if at the end of the string
		if i == len(backupOpts)-1 || c == ' ' && backupOpts[i+1] == '-' {
			if len(strings.Split(newField, " ")) > 1 && !strings.Contains(strings.Split(newField, " ")[0], "=") {
				splitFlagNoEqualsSign := strings.SplitN(newField, " ", 2)
				if (len(splitFlagNoEqualsSign)) > 1 {
					newFields = append(newFields, strings.TrimSpace(splitFlagNoEqualsSign[0]))
					newFields = append(newFields, strings.TrimSpace(splitFlagNoEqualsSign[1]))
				}
			} else {
				newFields = append(newFields, strings.TrimSpace(newField))
			}
			newField = ""
		}
	}

	return newFields
}

func ConvertBackupOptsToStruct(backupOpts []string) (*PgbasebackupOptions, []string) {

	optsStruct := &PgbasebackupOptions{}

	structValue := reflect.Indirect(reflect.ValueOf(optsStruct))
	structType := structValue.Type()

	commandLine := pflag.NewFlagSet("backup-opts", pflag.ExitOnError)

	for i := 0; i < structValue.NumField(); i++ {
		fieldVal := structValue.Field(i)

		flag, _ := structType.Field(i).Tag.Lookup("flag")
		flagShort, _ := structType.Field(i).Tag.Lookup("flag-short")

		switch fieldVal.Kind() {
		case reflect.String:
			commandLine.StringVarP(fieldVal.Addr().Interface().(*string), flag, flagShort, "", "")
		case reflect.Int:
			commandLine.IntVarP(fieldVal.Addr().Interface().(*int), flag, flagShort, 0, "")
		case reflect.Bool:
			commandLine.BoolVarP(fieldVal.Addr().Interface().(*bool), flag, flagShort, false, "")
		}
	}

	commandLine.Parse(backupOpts)

	commandLine.Visit(visitBackupOptFlags)
	return optsStruct, setFlagFieldNames
}

func visitBackupOptFlags(flag *pflag.Flag) {
	optsType := reflect.TypeOf(PgbasebackupOptions{})
	for i := 0; i < optsType.NumField(); i++ {
		field := optsType.Field(i)
		flagName, _ := field.Tag.Lookup("flag")
		flagNameShort, _ := field.Tag.Lookup("flag-short")
		if flag.Name == flagName || flag.Name == flagNameShort {
			setFlagFieldNames = append(setFlagFieldNames, field.Name)
		}
	}
}
