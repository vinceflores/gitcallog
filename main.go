package main

import (
  "fmt"
  "os/exec"
  "strings"
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
    output, err:= cmd.Output()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }

    output_str_arr  := strings.Split(string(output), "\n")

    // make a hashmap
    m := make(map[string]string)  
    for _, s := range output_str_arr {
      temp := strings.Split(s, " ")
      m[string(temp[0])] = temp[2]
    }
    fmt.Println(m) 
}
