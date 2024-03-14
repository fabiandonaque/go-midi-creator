package main

import(
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func hexToBytes(data ...string) ([]byte,error){
	str := strings.Join(data,"")
	raw,err := hex.DecodeString(str)
	return raw,err
}

func concatBytes(data ...[]byte) []byte {
	var raw []byte
	for _,d := range(data){
		raw = append(raw,d...)
	}
	return raw
}

func headerChunk() ([]byte,error){
	// Midi Header
	headerType,err := hexToBytes("4D","54","68","64")
	if err != nil { return nil,err }
	headerLength,err := hexToBytes("00","00","00","06")
	if err != nil { return nil,err }
	headerFormat,err := hexToBytes("00","00")
	if err != nil { return nil,err }
	headerNTracks,err := hexToBytes("00","01")
	if err != nil { return nil,err }
	headerDivision,err := hexToBytes("00","60")
	if err != nil { return nil,err }
	return concatBytes(headerType,headerLength,headerFormat,headerNTracks,headerDivision),nil
}

func trackChunk(events ...[]byte) ([]byte,error){
	// Midi Track
	trackType,err := hexToBytes("4D","54","72","6B")
	if err != nil { return nil,err }
	trackLength,err := hexToBytes("00","00","00","3B")
	if err != nil { return nil,err }
	track := append(trackType,trackLength...)
	for _,event := range(events){
		track = append(track,event...)
	}
	return track,nil
}

func trackEvent(delta []byte,event []byte) []byte {
	var d []byte
	d = append(d,delta...)
	d = append(d,event...)
	return d
}

func main(){
	header,err := headerChunk()
	if err != nil { fmt.Println(err); return }
	// Set events
	d1,err := hexToBytes("00")
	if err != nil { fmt.Println(err); return }
	d2,err := hexToBytes("60")
	if err != nil { fmt.Println(err); return }
	d3,err := hexToBytes("81","40")
	if err != nil { fmt.Println(err); return }
	e1,err := hexToBytes("FF","58","04","04","02","18","08")
	if err != nil { fmt.Println(err); return }
	e2,err := hexToBytes("FF","51","03","07","A1","20")
	if err != nil { fmt.Println(err); return }
	e3,err := hexToBytes("C0","05")
	if err != nil { fmt.Println(err); return }
	e4,err := hexToBytes("C1","2E")
	if err != nil { fmt.Println(err); return }
	e5,err := hexToBytes("C2","46")
	if err != nil { fmt.Println(err); return }
	e6,err := hexToBytes("92","30","60")
	if err != nil { fmt.Println(err); return }
	e7,err := hexToBytes("3C","60")
	if err != nil { fmt.Println(err); return }
	e8,err := hexToBytes("91","43","40")
	if err != nil { fmt.Println(err); return }
	e9,err := hexToBytes("90","4C","20")
	if err != nil { fmt.Println(err); return }
	e10,err := hexToBytes("82","30","40")
	if err != nil { fmt.Println(err); return }
	e11,err := hexToBytes("3C","40")
	if err != nil { fmt.Println(err); return }
	e12,err := hexToBytes("81","43","40")
	if err != nil { fmt.Println(err); return }
	e13,err := hexToBytes("80","4C","40")
	if err != nil { fmt.Println(err); return }
	e14,err := hexToBytes("FF","2F","00")
	if err != nil { fmt.Println(err); return }
	et1 := trackEvent(d1,e1)
	et2 := trackEvent(d1,e2)
	et3 := trackEvent(d1,e3)
	et4 := trackEvent(d1,e4)
	et5 := trackEvent(d1,e5)
	et6 := trackEvent(d1,e6)
	et7 := trackEvent(d1,e7)
	et8 := trackEvent(d2,e8)
	et9 := trackEvent(d2,e9)
	et10 := trackEvent(d3,e10)
	et11 := trackEvent(d1,e11)
	et12 := trackEvent(d1,e12)
	et13 := trackEvent(d1,e13)
	et14 := trackEvent(d1,e14)
	// Create track
	track,err := trackChunk(et1,et2,et3,et4,et5,et6,et7,et8,et9,et10,et11,et12,et13,et14)
	if err != nil { fmt.Println(err); return }
	// join all
	midi := append(header,track...)
	fmt.Println(midi)
	for _,d := range(midi){
		fmt.Printf(" %02X",d)
	}
	fmt.Println()


	// Save to file
	err = os.WriteFile("output.mid",midi,0644)
	if err != nil { fmt.Println(err) }
}