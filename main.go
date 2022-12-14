package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type GpxTrackPoint struct {
	XMLName xml.Name  `xml:"trkpt"`
	Lat     float64   `xml:"lat,attr"`
	Lon     float64   `xml:"lon,attr"`
	Time    time.Time `xml:"time"`
	Ele     float64   `xml:"ele"`
}

type GpxTrackSeg struct {
	XMLName     xml.Name        `xml:"trkseg"`
	TrackPoints []GpxTrackPoint `xml:"trkpt"`
}

type GpxTrack struct {
	XMLName  xml.Name      `xml:"trk"`
	Name     string        `xml:"name"`
	TrackSeg []GpxTrackSeg `xml:"trkseg"`
}

type GpxDoc struct {
	XMLName xml.Name   `xml:"gpx"`
	Tracks  []GpxTrack `xml:"trk"`
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: gpx2geojson <infile.gpx>")
		return
	}

	inFile := os.Args[1]

	xmlFile, err := os.Open(inFile)
	defer xmlFile.Close()

	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}

	var doc GpxDoc
	err = xml.Unmarshal(data, &doc)
	if err != nil {
		panic(err)
	}

	fc := geojson.NewFeatureCollection()

	for _, track := range doc.Tracks {
		ls := make(orb.LineString, 0)
		for _, ptSeg := range track.TrackSeg[0].TrackPoints {

			pt := orb.Point{ptSeg.Lon, ptSeg.Lat}
			ls = append(ls, pt)
		}
		feat := geojson.NewFeature(ls)
		feat.Properties["name"] = track.Name

		fc.Append(feat)
	}

	rawJSON, _ := fc.MarshalJSON()

	fmt.Println(string(rawJSON))
}
