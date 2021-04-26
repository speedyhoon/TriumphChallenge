package main

import (
	"fmt"
	"os/exec"
)

func openBrowser() {
	if err := exec.Command("xdg-open", natSoftURL).Start(); err != nil {
		fmt.Println("Unable to open a web browser for", natSoftURL, err)
	}
}
