package toutiao

import (
	"fmt"
	"testing"
)

func TestAESDecrypt(t *testing.T) {
	content, err := AESDecrypt(
		"wKbL7MvrmciitnkMCXvdePo7pbli//nr3nUW5NBj9UPBNsLWi+aWpGgEC7lrLN7HCt5eSqvO1jS0H2xkCEi5lMJHzPuO3IMzAMb8H0zL/XkMtviXUlviLNnB3tHGgmKOqSsuuH9rs7rSknlha68nO6kkKK+12ew+kn+FXPNoQd9L0eDOm1wdbUpNASWkrgS8Yw8Gp9c31OHzcJjv8OVhmw==",
		"Zg0x7DDrN56lG0fncoXQ9A==",
		"L2ZxbwZl5xwve0heY9l2Aw==",
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content)
	fmt.Println(string(content))
}
