module example.com/main

go 1.27

require (
	example.com/dep v0.0.0
	example.com/helper v0.0.0
)

replace (
	example.com/dep => ../dep
	example.com/helper => ../helper
)
