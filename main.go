package main

import (
	"fmt"
	"hackathon-basar-backend/internal"
)

func main() {

	fmt.Println(` 
██████╗  █████╗ ███████╗███████╗ █████╗ ██████╗ 
██╔══██╗██╔══██╗╚══███╔╝╚══███╔╝██╔══██╗██╔══██╗
██████╔╝███████║  ███╔╝   ███╔╝ ███████║██████╔╝
██╔══██╗██╔══██║ ███╔╝   ███╔╝  ██╔══██║██╔══██╗
██████╔╝██║  ██║███████╗███████╗██║  ██║██║  ██║
╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝
	`)

	internal.LoadApplicationConfig()

	chi := internal.ChiConfig()

	internal.Serve(chi)
}
