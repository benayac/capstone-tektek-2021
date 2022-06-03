package helper

import (
	"fmt"
	"math"
)

func GetDistanceFromLatLng(latFirst float64, lonFirst float64, latSecond float64, lonSecond float64) float64 {
	R := 6371.0
	dLat := degToRad(latSecond - latFirst)
	dLon := degToRad(lonSecond - lonFirst)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(degToRad(latFirst))*math.Cos(degToRad(latSecond))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := R * c
	fmt.Printf("latFirst: %v\nLonFirst: %v\nLatSecond: %v\nLonSecond: %v\n", latFirst, lonFirst, latSecond, lonSecond)
	fmt.Printf("Distance: %v\n", distance)
	return distance
}

func degToRad(deg float64) float64 {
	return deg * (math.Phi / 180)
}
