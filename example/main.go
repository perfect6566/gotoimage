package main

import (
	"fmt"
	"github.com/perfect6566/gotoimage"
	"log"
	"time"
)

//you can create the source code file with golang, then input the filename
//and it will go to https://play.golang.org/ to excute your code and give back the
//screenshot as png file for your code execution result screenshot
//please input the correct and existing filename ,or it will fall in loop until it get
//correct source file
func main()  {

getsourcefile:

	log.Println("Please input the  source code file name ... ")

	var filename string
	_,err:=fmt.Scan(&filename)
	if err!=nil{
		log.Println(err)
	}
	_,err=gotoimage.Render(filename)
	if err!=nil{
		log.Println(err)
		time.Sleep(5*time.Second)
		if err.Error()=="source file is not exist"{
			goto getsourcefile
		}
	}
}
