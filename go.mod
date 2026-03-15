module uploadpass

go 1.23.0

toolchain go1.23.3

require (
	golang.org/x/crypto v0.31.0
	modernc.org/sqlite v1.38.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	modernc.org/libc v1.65.10 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

replace (
	golang.org/x/exp v0.0.0-20230315142452-642cacee5cc0 => golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842
	golang.org/x/mod v0.19.0 => golang.org/x/mod v0.21.0
	golang.org/x/sync v0.7.0 => golang.org/x/sync v0.17.0
	golang.org/x/tools v0.23.0 => golang.org/x/tools v0.24.0
	modernc.org/fileutil v1.3.3 => modernc.org/fileutil v1.3.0
	modernc.org/libc v1.65.10 => modernc.org/libc v1.61.9
)
