package image

import (
	"fmt"
	"os"

	"github.com/TheRaizer/GolangGame/util"
)

func DecodeImage(name string) {
	file, err := os.Open(name)

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	util.CheckErr(err)

	buffer := make([]byte, 5)
	numBytesRead, err := file.Read(buffer)

	util.CheckErr(err)

	fmt.Println(buffer)
	fmt.Println(numBytesRead)
}
