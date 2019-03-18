package pgovalidate

type PgbasebackupOptions struct {
	Gzip     bool   `flag:"gzip" flag-short:"z"`
	Label    string `flag:"label" flag-short:"l" validate:"required"`
	Progress bool   `flag:"progress" flag-short:"P"`
	Compress int    `flag:"compress-level" flag-short:"Z" validate:"numeric,min=0,max=9"`
}
