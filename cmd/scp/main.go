package main

import "github.com/APwhitehat/goscp"

func main() {
	var scpOp goscp.ScpOptions

	scpOp.Hostname = "hello"
	goscp.Scp(scpOp, "src", "desc")
}
