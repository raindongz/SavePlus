package utils

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numbers = "0123456789"

func init(){
	rand.NewSource(time.Now().UnixNano())
}

// randomInt generate random integer between min and max
func RandomInt(min, max int) int{
	return min + rand.Intn(max - min + 1)
}

func RandomStringWithSpecifiedLenth(n int) string{
	var sb strings.Builder
	k := len(alphabet)

	for i:=0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}


func RandomEmail()string{
	return fmt.Sprintf("%s@email.com", RandomStringWithSpecifiedLenth(6))
}

func RandomAvatar()string{
	return fmt.Sprintf("www.%s.com", RandomStringWithSpecifiedLenth(6))
}

func RandomPhoneNumber()string{
	// return fmt.Sprintf("+1%s", strconv.Itoa(int(RandomInt(10,10))))
	var sb strings.Builder
	lenth := len(numbers)

	for i:=0; i<10; i++{
		ranNum := numbers[rand.Intn(lenth)]
		sb.WriteByte(byte(ranNum))
	}
	log.Println(fmt.Sprintf("+1%s", sb.String()))
	return fmt.Sprintf("+1%s", sb.String()) 
}