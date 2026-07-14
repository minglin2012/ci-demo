package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║   CI Multi-Platform Demo App v1.0   ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("  Platform   : %s\n", runtime.GOOS)
	fmt.Printf("  Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("  Go Version : %s\n", runtime.Version())
	fmt.Println()

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("CI Demo Tool - Version 1.0.0")
		fmt.Println("Built with GitHub Actions CI/CD")
		os.Exit(0)
	}

	fmt.Println("Hello from the CI Demo application!")
	fmt.Println("This binary was built automatically via GitHub Actions.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ci-demo           Run the application")
	fmt.Println("  ci-demo --version  Show version information")
}
