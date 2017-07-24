package whatsApp_Analytics

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"sort"
)

//Sort ----------------------
// A data structure to hold key/value pairs
type Pair struct {
	Key   string
	Value int
}

// A slice of pairs that implements sort.Interface to sort by values
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadFile(filename string) ([]string, error) {
	OpenedFile, err := os.Open(filename)
	Check(err)

	// close fi on exit and Check for its returned error
	defer func() {
		if err := OpenedFile.Close(); err != nil {
			fmt.Println("Closing: " + OpenedFile.Name())
			panic(err)
		}
	}()
	fmt.Println("Opened: " + OpenedFile.Name())

	var lines []string

	buffer := bufio.NewScanner(OpenedFile)
	for buffer.Scan() {
		lines = append(lines, buffer.Text())
	}
	return lines, buffer.Err()
}

func ParseData(data []string) (chatContainer, int) {
	tempbuffer := make([]messageData, 0)
	lastLine := int(0)

	//big container with users data and parameters
	tempChat := chatContainer{}

	tempChat.Init()

	//loop through data
	for i := range data {

		SplitMain := strings.SplitN(data[i], "- ", 2)
		if len(SplitMain) >= 2 {
			splitmeta := SplitMain[0]
			splitdata := SplitMain[1]

			if len(SplitMain) == 2 {
				splitMeta := strings.Split(splitmeta, ", ")

				splitData := strings.Split(splitdata, ": ")

				if len(splitData) != 1 || len(splitData) != 0 {

					if len(splitData) > 2 {
						for x := 2; x <= len(splitData); x++ {
							splitData[1] += splitData[x-1]
						}
					}

					if len(splitData) == 1 {
						splitData = append(splitData, "")
					}

					tempChat.AddSender(splitData[0])
					tempChat.AddString(splitData[0], strings.ToLower(splitData[1]))
					if len(splitMeta) == 1 {
						fmt.Println(splitMeta[0], lastLine, "  ", i)
						//tempbuffer[lastLine].Message += data[i]
					} else {
						lastLine = i
						tempbuffer = append(tempbuffer, messageData{splitMeta[0], splitMeta[1], splitData[0], strings.ToLower(splitData[1])})

					}
				}
			}

		} else {
			//some code to add text to last message
			//fmt.Println(len(SplitMain))
		}
	}
	tempChat.messages = tempbuffer
	return tempChat, len(tempbuffer)
}

type messageData struct {
	Date    string
	Time    string
	Sender  string
	Message string
}

type chatContainer struct {
	messages          []messageData
	senders           []string
	sendersParameters map[string][]string
	sendersCount      map[string]int
}

func (c *chatContainer) AddSender(name string) {
	for _, n := range c.senders {
		if n == name {
			return
		}
	}
	c.senders = append(c.senders, name)
}

func (c *chatContainer) Init() {
	c.messages = make([]messageData, 0)
	c.senders = make([]string, 0)
	c.sendersParameters = make(map[string][]string)
	c.sendersCount = make(map[string]int)
}

func (c *chatContainer) AddString(name string, data string) {
	c.sendersParameters[name] = append(c.sendersParameters[name], data)
}

func (c *chatContainer) IncrementCount(name string) {
	c.sendersCount[name]++
}

func (c *chatContainer) CLeanAndVerifyNames() {
	senders := make([]string, 0)

	//loop to run through the names and Check they dont have empty chat log
	for _, x := range c.senders {
		for _, u := range c.sendersParameters[x] {
			if u == "" {
				delete(c.sendersParameters, x)
				senders = append(senders, x)

			}
		}
		if len(c.sendersParameters[x]) == 1 {
			delete(c.sendersParameters, x)
		} else {
			senders = append(senders, x)
		}
	}

	c.senders = senders
}

func (c *chatContainer) FindFirstOccurunce(word string)(string){
	for _,name := range c.messages{
		if strings.Contains(name.Message,strings.ToLower(word)){
			return name.Sender
		}
	}
	return "No One"
}

func (c *chatContainer) FindByString(word string, reportTitle string) (map[string]int) {
	//analytics
	for _, name := range c.senders {
		for _, n := range c.sendersParameters[name] {
			if strings.Contains(n, strings.ToLower(word)) {
				c.IncrementCount(name)
			}
		}
	}

	p := make(PairList, len(c.sendersCount))

	i := 0
	for k, v := range c.sendersCount {
		p[i] = Pair{k, v}
		i++
	}

	sort.Sort(p)

	fmt.Println("")
	fmt.Println("")
	fmt.Println(reportTitle)
	for _, t := range p {
		fmt.Println(t.Key, " :", t.Value)
	}
	trannySum := int(0)

	for _, t := range p {
		trannySum+= t.Value
	}
	fmt.Println("Total number of occurunces = ",trannySum)
	temp:= c.sendersCount
	c.sendersCount = make(map[string]int)
	return temp

}
