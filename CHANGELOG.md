# Changelog

## v0.9.0 (February 6, 2021)

Big thank you to [mattlqx](https://github.com/mattlqx) for the great enhancements.

ENHANCEMENTS:

* Proper arg parsing with help text for subcommands ([#73](https://github.com/fishi0x01/vsh/pull/73) - Thank you for implementation [mattlqx](https://github.com/mattlqx))
* Add replace command ([#69](https://github.com/fishi0x01/vsh/pull/69) - Thank you for implementation [mattlqx](https://github.com/mattlqx))
* Add key selector to replace command ([#72](https://github.com/fishi0x01/vsh/pull/72) - Thank you for implementation [mattlqx](https://github.com/mattlqx))
* Allow limiting scope of grep to keys or values ([#66](https://github.com/fishi0x01/vsh/pull/66) - Thank you for implementation [mattlqx](https://github.com/mattlqx))
* Do not show and operate on KV2 metadata ([#68](https://github.com/fishi0x01/vsh/pull/68))

## v0.8.0 (January 27, 2021)

ENHANCEMENTS:

* Allow regex on `grep` operation ([#61](https://github.com/fishi0x01/vsh/pull/61) - Thank you for implementation [mattlqx](https://github.com/mattlqx))
* Allow quotes and escapes in input ([#61](https://github.com/fishi0x01/vsh/pull/61) - Thank you for implementation [mattlqx](https://github.com/mattlqx))

BUG FIXES:

* Fix panic on `data` keys in KV1 ([#63](https://github.com/fishi0x01/vsh/pull/63) - Thank you for issue submission [tommartensen](https://github.com/tommartensen))

## v0.7.2 (October 4, 2020)

BUG FIXES:

* Fix copy of ambiguous sub-file ([#55](https://github.com/fishi0x01/vsh/pull/55))

## v0.7.1 (October 2, 2020)

BUG FIXES:

* Proper return codes ([#51](https://github.com/fishi0x01/vsh/pull/51))
* Proper logging ([#52](https://github.com/fishi0x01/vsh/pull/52))

## v0.7.0 (September 26, 2020)

ENHANCEMENTS:

* Add option to disable path auto-completion ([#48](https://github.com/fishi0x01/vsh/pull/48))

## v0.6.3 (September 22, 2020)

BUG FIXES:

* Properly handle ambiguous files without read permission [#46](https://github.com/fishi0x01/vsh/pull/46) - Thank you for detailed report [agaudreault-jive](https://github.com/agaudreault-jive)

## v0.6.2 (June 18, 2020)

ENHANCEMENTS:

* Allow separation of command params with multiple whitespaces ([#40](https://github.com/fishi0x01/vsh/pull/40)) - Thank you [vikin91](https://github.com/vikin91)
* Unify verbose behavior ([#41](https://github.com/fishi0x01/vsh/pull/41)) - Thank you [vikin91](https://github.com/vikin91)

BUG FIXES:

* Properly handle ambiguous files ([#37](https://github.com/fishi0x01/vsh/pull/37))

## v0.6.1 (June 15, 2020)

BUG FIXES:

* Properly handle suffix '/' in source path ([#39](https://github.com/fishi0x01/vsh/pull/39))

## v0.6.0 (June 13, 2020)

ENHANCEMENTS:

* add `append` command ([#30](https://github.com/fishi0x01/vsh/issues/30)) - Thank you [vikin91](https://github.com/vikin91)

BUG FIXES:

* Remove file/dir ambiguity for `rm` ([#29](https://github.com/fishi0x01/vsh/issues/29))

## v0.5.0 (April 5, 2020)

ENHANCEMENTS:

* add `grep` command ([#25](https://github.com/fishi0x01/vsh/issues/25))
* `ls` with new line instead of single line ([#27](https://github.com/fishi0x01/vsh/issues/27))

BUG FIXES:

* remove `//` from paths ([#26](https://github.com/fishi0x01/vsh/issues/26))
* fix broken tests

## v0.4.1 (March 21, 2020)

ENHANCEMENTS:

* performance: cache `List()` queries ([#23](https://github.com/fishi0x01/vsh/issues/23))

## v0.4.0 (March 10, 2020)

ENHANCEMENTS:

* use TokenHelper mechanism ([#20](https://github.com/fishi0x01/vsh/issues/20))

## v0.3.1 (October 31, 2019)

BUG FIXES:

* fix top-level path panic ([#17](https://github.com/fishi0x01/vsh/issues/17))

## v0.3.0 (October 31, 2019)

ENHANCEMENTS:

* token list permission on sys/mounts is not mandatory

## v0.2.0 (October 20, 2019)

ENHANCEMENTS:

* use `~/.vault-token` as fallback if `VAULT_TOKEN` is not set ([#12](https://github.com/fishi0x01/vsh/issues/12))

BUG FIXES:

* error handling to catch bad input ([#13](https://github.com/fishi0x01/vsh/issues/13))

## v0.1.1 (October 8, 2019)

BUG FIXES:

* more sanity checks on user input to avoid crashes ([#10](https://github.com/fishi0x01/vsh/issues/10))

## v0.1.0 (October 7, 2019)

Initial release
