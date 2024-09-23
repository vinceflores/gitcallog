package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetLogMap() (map[string][]string , error)  {
	// Define the git command with arguments
	cmd := exec.Command("git", "log", "--pretty=format:%ad/%s", "--date=short")
	// Run the command and capture the output
	output, err:= cmd.Output()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil , err
	}
  
	output_str_arr  := strings.Split(string(output), "\n")
  
	// make a hashmap
	m := make(map[string][]string)  
	for _, s := range output_str_arr {
	  temp := strings.Split(s, "/")
	  // iterate through temp and append to hashmap
	  for i , v := range temp {
		if i == 0 {
		  continue
		}else{
		  m[string(temp[0])] = append(m[string(temp[0])], string(v))
		}
	  }
	}

	return m , nil
  }