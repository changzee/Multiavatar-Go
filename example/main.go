package main

import (
	"fmt"
	"log"
	"os"

	"github.com/changzee/multiavatar-go"
)

func main() {
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// --- Example 1: Generate a default avatar ---
	inputString1 := "Binx Bond123"
	svgDefault := multiavatar.Generate(inputString1)
	filePathDefault := "output/avatar_default.svg"
	err := os.WriteFile(filePathDefault, []byte(svgDefault), 0644)
	if err != nil {
		log.Fatalf("Failed to write default avatar to file: %v", err)
	}
	fmt.Printf("Successfully generated default avatar for '%s' at %s\n", inputString1, filePathDefault)

	// --- Example 2: Generate an avatar with a transparent background ---
	inputString2 := "John Doe"
	svgTransparent := multiavatar.Generate(inputString2, multiavatar.WithoutBackground())
	filePathTransparent := "output/avatar_transparent.svg"
	err = os.WriteFile(filePathTransparent, []byte(svgTransparent), 0644)
	if err != nil {
		log.Fatalf("Failed to write transparent avatar to file: %v", err)
	}
	fmt.Printf("Successfully generated transparent avatar for '%s' at %s\n", inputString2, filePathTransparent)

	// --- Example 3: Generate another avatar to show determinism ---
	// Using the same input string as example 1 should produce the exact same SVG
	svgDeterministic := multiavatar.Generate(inputString1)
	if svgDefault != svgDeterministic {
		log.Fatal("Error: Deterministic algorithm failed. Avatars with same input are different.")
	}
	fmt.Printf("Successfully confirmed deterministic generation for '%s'\n", inputString1)

	fmt.Println("\nAll examples executed successfully.")
}
