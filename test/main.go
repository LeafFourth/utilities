package main

//import "unsafe"
//import "fmt"
import "utilities"

type A int
type B *A

var logger *utilities.Logger

func main() {
	logger = utilities.NewLogger()

	logger.SetLogPath("E:/log.txt");

	logger.Log("123");

	logger.SetLogPath("E:/log2.txt");
	logger.Log("3", "4", "5");

}