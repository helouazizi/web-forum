package main

/*
#include <stdio.h>
void hello() {
    printf("Hassan using C!\n");
}
*/
import "C"

func main() {
	C.hello() // Calls the C function from Go
}
