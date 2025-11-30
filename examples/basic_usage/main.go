package main

import (
	"fmt"
	"log"

	"github.com/mikeschinkel/go-dt"
)

func main() {
	// Example 1: Parse and use DirPath
	dirPath, err := dt.ParseDirPath("/usr/local/bin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Directory path: %s\n", dirPath)

	// Example 2: Use Filename type
	filename := dt.Filename("config.json")
	fmt.Printf("Filename: %s\n", filename)

	// Example 3: Use Filepath type
	filepath := dt.Filepath("/etc/config/settings.conf")
	fmt.Printf("Filepath: %s\n", filepath)

	// Example 4: Join path components using method
	joined := dirPath.Join(filename)
	fmt.Printf("Joined path (method): %s\n", joined)

	// Example 5: Join path components using function
	joined2 := dt.FilepathJoin(string(dirPath), string(filename))
	fmt.Printf("Joined path (function): %s\n", joined2)

	fmt.Println("\nBasic go-dt types demonstration complete!")
}
