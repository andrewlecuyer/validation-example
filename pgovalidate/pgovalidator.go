package pgovalidate

import (
	validator "gopkg.in/go-playground/validator.v9"
)

func ValidateBackupOpts(optsStruct *PgbasebackupOptions, setFlagFieldNames []string) error {
	validate := validator.New()
	err := validate.StructPartial(optsStruct, setFlagFieldNames...)
	return err
}
