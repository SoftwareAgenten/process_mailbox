// process_mailbox
//
// Part of the Software-Agenten im Internet project.
// <https://github.com/SoftwareAgenten>
// <https://github.com/SoftwareAgenten/process_mailbox>
//
// The MIT License (MIT)
//
// Copyright (c) 2016 Florian Pircher
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Comment struct {
	Id       string   `json:"id"`
	Datetime Datetime `json:"datetime"`
	Name     string   `json:"name"`
	Message  string   `json:"message"`
}

type Datetime struct {
	Day   int    `json:"day"`
	Month string `json:"month"`
	Year  int    `json:"year"`
	Time  Time   `json:"time"`
}

type Time struct {
	Hours    int    `json:"hours"`
	Minutes  int    `json:"minutes"`
	Seconds  int    `json:"seconds"`
	Timezone string `json:"timezone"`
}

// regex, send date: day, month, year, time (H:i:s), timezone
const DATE string = `Date: .*(\d+) (\w{3}) (\d+) (\d+):(\d+):(\d+) (.\d{4})`

// regex, message id: id
const ID string = `Message-id: <(.*)>`

// regex, commentator name: name
const NAME string = `Name<\/strong>: ([^<]+)<\/p>`

// regex, html message: message
const MESSAGE string = `Name<\/strong>:.*?<\/p><p>([\s\S]+?)<\/p>`

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: process_mailbox source [target=source.json]")
		os.Exit(0)
	}

	source := os.Args[1]
	target := source + ".json"

	if len(os.Args) > 2 {
		target = os.Args[2]
	}

	// read argument file
	data, err := ioutil.ReadFile(source)
	checkError(err)
	contents := string(data)
	mails := strings.Split(contents, "From kommentar@blogfill.de")[1:]

	reDate := regexp.MustCompile(DATE)
	reId := regexp.MustCompile(ID)
	reName := regexp.MustCompile(NAME)
	reMessage := regexp.MustCompile(MESSAGE)

	f, err := os.Create(target)
	checkError(err)

	defer f.Close()

	f.WriteString("[")

	for i, mail := range mails {
		rawDT := reDate.FindAllStringSubmatch(mail, -1)[0]
		hours, _ := strconv.Atoi(rawDT[4])
		minutes, _ := strconv.Atoi(rawDT[5])
		seconds, _ := strconv.Atoi(rawDT[6])
		timezone := rawDT[7]
		time := Time{hours, minutes, seconds, timezone}
		day, _ := strconv.Atoi(rawDT[1])
		year, _ := strconv.Atoi(rawDT[3])
		datetime := Datetime{day, rawDT[2], year, time}
		id := reId.FindAllStringSubmatch(mail, -1)[0][1]
		name := reName.FindAllStringSubmatch(mail, -1)[0][1]
		message := reMessage.FindAllStringSubmatch(mail, -1)[0][1]
		comment := Comment{id, datetime, name, message}

		foo, _ := json.Marshal(comment)
		f.Write(foo)

		if i < len(mails)-1 {
			f.WriteString(",")
		}
	}

	f.Sync()
	f.WriteString("]")
}
