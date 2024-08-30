package main

import (
  "fmt"
  "os/exec"
)


  /*
     1. grep git log  and save to a file
        1.1 perform using golang or perform command line  script
        1.2 send output to array
     2. map through an array
2.1 hashmap ; key=":wq
"

*/

//2024-06-28 - finish app
// const (
// comd = "git log --pretty=format:"%ad - %s" --date=short"
//)

func main(){
      // Define the git command with arguments
    cmd := exec.Command("git", "log", "--pretty=format:%ad - %s", "--date=short")

    // Run the command and capture the output
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }

    // Print the output of the command
    fmt.Printf("Output:\n%s\n", output)
}
