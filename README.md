## Instructions for Running the Example

1.) First get the required dependencies.  Please note that **pflags** is used by Cobra, and is therefore an existing dependency in PGO, _NOT_ a new one.  [**The Validator package**](https://godoc.org/gopkg.in/go-playground/validator.v9) is the only new dependency:

```bash
go get github.com/spf13/pflag
go get gopkg.in/go-playground/validator.v9
```

2.) Clone the repo for the example validation application and install it:

```bash
git clone https://github.com/andrewlecuyer/validation-example.git
go install validate.go
```

3.) For the first example, we will assume we want to support three flags for **pg_basebackup**, which we will pass to the example application as follows:

```bash
validate --backup-opts="--gzip --label=mylabel --progress"
```

This will result in the validation of all flags, in essentially the same manner that they would be validated in PGO.  However, for this first run, since all flags are formatted correctly, don't contain typos, invalid values, etc., you will see a response indicating that all flags were valid:

```bash
flags are valid!
```

4.) Now lets say the user mistypes **gzip** by accidentally adding an extra 'p', such as follows:

```bash
validate --backup-opts="--gzipp --label=mylabel --progress"
```

This will result in pflags, the library used to parse the command line options provided (which is also the same flag package utilized by Cobra), reporting that `gizpp` is an unrecognized flag, as can be seen in the following error message:

```bash
unknown flag: --gzipp
Usage of backup-opts:
  -z, --gzip
  -l, --label string
  -P, --progress
unknown flag: --gzipp
```

So as you can see, the error shows that the flag provided was not recognized, while also clearly displaying to the user the values (i.e. backup options) that we do support.  In this sense the use of pflags allows us to establish a whitelist of flags that we support, displaying a consistent error message for any that we don't, similar to what the user would see with any command line application.  And as shown above, this approach also catches typos in any flags provided by the user while also parsing backup options provided in a POSIX-compliant manner, i.e. it should handle quotes, escapes, values provided with or without an equals sign, etc.  And all of this is easily done by leveraging the power of pflags.

5.) For this next example, we will correct the `gzip` flag to ensure all flags provided are one again valid (as determined by pflags).  However, this time we will include an error in the value provided for one of the options, which will be detected by the new validator package.  Specifically, we will specify the label flag, without specifying an actual label value:

```bash
validate --backup-opts="--gzip --label= --progress"
```

As you will see, an error is thrown indicating that an actual value must be specified for the label flag, i.e. a value is required when the label flag is included.

```bash
Key: 'PgbasebackupOptions.Label' Error:Field validation for 'Label' failed on the 'required' tag
```

And from a source code perspective, this is done by simply updating a struct containing all supported pg_basebackup options.  Specifically, the **Label** field specifies via a tag that it is `required`, with `required` being a built-in validation function provided by the validator package (of which there are many) that can be used to ensure a value has been provided for a specific field.  This can be seen in the following code snippet from `pgovalidate\backuptypes.go`:

```go
type PgbasebackupOptions struct {
	Gzip     bool   `flag:"gzip" flag-short:"z"`
	Label    string `flag:"label" flag-short:"l" validate:"required"`
	Progress bool   `flag:"progress" flag-short:"P"`
}
```

And thanks to the validation library, it is easy to only validate those options in the struct that were actually provided by the user when running a command.  This will, for instance, ensure that above error isn't thrown for the `label` option if the user does not actually provide that flag on the command line.

6.) And lastly, to again show the power of this solution when it comes to maintenance, lets say our users come back and say that while `gzip` support is great, they would now also like to be able to specify a compression level (i.e. pg_basebackup option `--compress=level`).

With this solution, only a single line of code would be required to support and validate this new option.  To demonstrate, `pgovalidate\backuptypes.go` can be updated as follows, adding the compression level option to type `PgbasebackupOptions`:

```go
type PgbasebackupOptions struct {
	Gzip     bool   `flag:"gzip" flag-short:"z"`
	Label    string `flag:"label" flag-short:"l" validate:"required"`
	Progress bool   `flag:"progress" flag-short:"P"`
	Compress int    `flag:"compress-level" flag-short:"Z" validate:"numeric,min=0,max=9"`
}
```

You can now re-install the validate application again to pick up this change, while then running another `validate` command, this time specifying a compression level via the `compress-level` flag:

```bash
go install validate.go
validate --backup-opts="--gzip --label=mylabel --progress --compress-level=10"
```

However, even though `--compress-level` will now be accepted as a valid/support backup option, by running the above command you will see an error, since 10 was provided for the compression level, and the max specified for this value via the `validate` tag in the struct was 9.

```bash
Key: 'PgbasebackupOptions.Compress' Error:Field validation for 'Compress' failed on the 'max' tag
```

So if you update the command to specify a value between 0 and 9, you should now once again see a message indicating that all flags were valid:

```bash
$ validate --backup-opts="--gzip --label=mylabel --progress --compress-level=5"
flags are valid!
```

As you can see above, the validate tag for the compression level specifies a few built-in validation functions provided by the validation package, specifically `numeric`, `min` and `max` functions, which verify that the value provided is numeric, and is a value between 0 and 10.  This again prevents us from having code commonly-used validation logic ourselves, while also showing how easy it is to add/remove/update validation logic for any specific option we support.  If you look a file `pgovalidate\pgovalidator.go`, you will see that the body of the current method to validate the backup opts only contains three lines of code thanks to the validator package and these built-in functions, which should cover the majority of our use cases.

So again, most changes in the future should require little more than a single line of code.  This means when future options are added to the backup/restore utilities across various PG and pgBackRest releases (and I have verified that new options have been added across that last 4 or so PG releases), and users are looking to utilize those new options, the changes required to not only allow their use, but also effectively validate them, should be trivial.
