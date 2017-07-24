# whatsApp Analytics
Parse WhatsApp group conversations and extract simple information.

Example:
```
package main

import (
	"github.com/mmwa/whatsApp_Analytics"
	"fmt"
)

func main() {
	//import file into program
	data, err := whatsApp_Analytics.ReadFile("data/testSet.txt")
	whatsApp_Analytics.Check(err)

	//align
	dataSet, _ := whatsApp_Analytics.ParseData(data)

	//Make sure the names are names
	dataSet.CLeanAndVerifyNames()
	dataSet.FindByString("Happy", "Freinds Group Report report:")



	fmt.Println("First One to say Eid: ",dataSet.FindFirstOccurunce("Eid"))

	return
}
```
