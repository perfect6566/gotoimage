# gotoimage    

A client for remote executing golang for Go users 
If you don't have a go environment and you want to execute your go code, Feel free to use this tool to execute your code from remote go environment     

```go
import "github.com/perfect6566/gotoimage"
```



## Basic usage notes:    

1. go build -o gotoimage example/main.go ; Build the execute file named gotoimage

2. Run gotoimage.exe (windows) or gotoimage(for linux),Input the source golang code filename, then it will open "https://play.golang.org/" with your golang by chrome, and paste your code to the webpage execution area,then it will click the format and run button to run your code, after Execution Completed, it will screenshot for the chrome page result  and save to  filename.png


