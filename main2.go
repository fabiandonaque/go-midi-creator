package main

///////////////
//  Imports  //
///////////////

import(
	"encoding/binary"
	"fmt"
	"os"
)
/////////////////
//  Constants  //
/////////////////

var Divisions uint16 = 96

///////////////
//  Structs  //
///////////////

type Channel struct {
	Track *Track
	Number byte
}

type Event struct {
	Delta int64
	Data []byte
}

type Track struct {
	BPM int64
	Events []Event
}

type Midi struct {
	Name string
	Tracks [](*Track)
}

/////////////////
//  Functions  //
/////////////////

func newMidi(name string) *Midi {
	m := Midi{
		Name: name,
	}
	return &m
}

func (m *Midi) newTrack(bpm int64) *Track {
	t := Track{
		BPM: bpm,
	}
	m.Tracks = append(m.Tracks,&t)
	// Set Time Signature
	(&t).newEvent(0,[]byte{0xFF,0x58,0x04,0x04,0x02,0x18,0x08})
	// Set Key Signature C Maj
	(&t).newEvent(0,[]byte{0xFF,0x59,0x02,0x00,0x00})
	// Set Tempo
	event := []byte{0xFF,0x51,0x03}
	micro := 60000000/bpm
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(micro))
	flag := false
	for _,b := range(buf) {
		if b != 0 || flag {
			flag = true
			event = append(event,b)
		}
	}
	(&t).newEvent(0,event)
	return &t
}

func (m *Midi) save() {
	data := []byte{0x4D,0x54,0x68,0x64,0x00,0x00,0x00,0x06,0x00,0x00,0x00,0x01}
	div := make([]byte, 2)
	binary.BigEndian.PutUint16(div, Divisions)
	data = append(data,div...)
	for _,t := range(m.Tracks){
		var subdata []byte
		for _,e := range(t.Events){
			delta := make([]byte, 8)
			binary.BigEndian.PutUint64(delta, uint64(e.Delta))
			flag := false
			for _,b := range(delta) {
				if b != 0 || flag {
					flag = true
					subdata = append(subdata,b)
				}
			}
			if !flag{
				subdata = append(subdata,0)
			}
			for _,d := range(e.Data){
				subdata = append(subdata,d)
			}
		}
		l := uint32(len(subdata))
		len := make([]byte, 4)
		binary.BigEndian.PutUint32(len, l)
		trackData := []byte{0x4D,0x54,0x72,0x6B}
		trackData = append(trackData,len...)
		trackData = append(trackData,subdata...)
		data = append(data,trackData...)
	}
	// Save to file
	err := os.WriteFile(m.Name+".mid",data,0644)
	if err != nil { fmt.Println(err) }
}

func (t *Track) newEvent(delta int64,data []byte) {
	e := Event{
		Delta: delta,
		Data: data,
	}
	t.Events = append(t.Events,e)
}

func (t *Track) endOfTrack() {
	t.newEvent(0,[]byte{0xFF,0x2F,0x00})
}

func (c *Channel) setInstrument(instrument byte){
	n := 192+c.Number
	c.Track.newEvent(0,[]byte{n,instrument})
}

func (c *Channel) setOnOffNote(duration float64,note byte,intensity float64){
	on := 144+c.Number
	off := 128+c.Number
	i := byte(intensity*127)
	d := int64(duration*float64(Divisions))
	c.Track.newEvent(0,[]byte{on,note,i})
	c.Track.newEvent(d,[]byte{off,note,0})
}

func (c *Channel) setOnNote(duration float64,note byte,intensity float64){
	off := 128+c.Number
	i := byte(intensity*127)
	d := int64(duration*float64(Divisions))
	c.Track.newEvent(d,[]byte{off,note,i})
}

func (c *Channel) setOffNote(duration float64,note byte,intensity float64){
	on := 144+c.Number
	i := byte(intensity*127)
	d := int64(duration*float64(Divisions))
	c.Track.newEvent(d,[]byte{on,note,i})
}

////////////
//  Main  //
////////////

func main(){
	// Create Midi
	m := newMidi("prueba")
	// Set new track
	t := m.newTrack(120)
	// Create Channel
	c := Channel{
		Track: t,
		Number: 1,
	}
	// Set instrument
	c.setInstrument(15)
	// Set arpeggio
	c.setOnOffNote(1.0,48,0.7)
	c.setOnOffNote(1.0,51,0.7)
	c.setOnOffNote(1.0,55,0.7)
	c.setOnOffNote(1.0,58,0.7)
	// Set chord
	c.setOnNote(0.0,48,0.7)
	c.setOnNote(0.0,51,0.7)
	c.setOnNote(0.0,55,0.7)
	c.setOnNote(0.0,58,0.7)
	c.setOffNote(1.0,48,0.7)
	c.setOffNote(0.0,51,0.7)
	c.setOffNote(0.0,55,0.7)
	c.setOffNote(0.0,58,0.7)
	// End of track
	t.endOfTrack()

	m.save()
}