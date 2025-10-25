package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/cespare/xxhash/v2"
)

// CollisionStats holds the results of a collision test
type CollisionStats struct {
	NumInputs         int
	UniqueHashes      int
	Collisions        int
	CollisionRate     float64
	ExpectedCollisions float64
	Duration          time.Duration
	DuplicateInputs   int
}

// GenerateRandomString creates a random string of given length
func GenerateRandomString(length int, r *rand.Rand) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// TestCollisionRate tests xxHash collision rate with random inputs
func TestCollisionRate(numInputs int, stringLength int) CollisionStats {
	start := time.Now()
	
	// Map hash to original input to detect TRUE collisions
	hashToInput := make(map[uint64]string)
	collisions := 0
	duplicateInputs := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < numInputs; i++ {
		// Generate random input
		data := GenerateRandomString(stringLength, r)
		
		// Compute xxHash64
		hash := xxhash.Sum64String(data)
		
		if originalInput, exists := hashToInput[hash]; exists {
			// Hash collision detected - check if it's a TRUE collision
			if originalInput != data {
				collisions++
				fmt.Printf("TRUE COLLISION FOUND!\n")
				fmt.Printf("  Input 1: %s (hash: %d)\n", originalInput, hash)
				fmt.Printf("  Input 2: %s (hash: %d)\n", data, hash)
			} else {
				duplicateInputs++
			}
		} else {
			hashToInput[hash] = data
		}
	}
	
	duration := time.Since(start)
	collisionRate := float64(collisions) / float64(numInputs)
	
	// Calculate expected collisions using birthday paradox approximation
	// E[collisions] ≈ n² / (2 * 2^bits)
	hashBits := 64.0
	expectedCollisions := (float64(numInputs) * float64(numInputs-1)) / (2.0 * math.Pow(2, hashBits))
	
	if duplicateInputs > 0 {
		fmt.Printf("Note: %d duplicate inputs were generated (not counted as collisions)\n", duplicateInputs)
	}
	
	return CollisionStats{
		NumInputs:         numInputs,
		UniqueHashes:      len(hashToInput),
		Collisions:        collisions,
		CollisionRate:     collisionRate,
		ExpectedCollisions: expectedCollisions,
		Duration:          duration,
		DuplicateInputs:   duplicateInputs,
	}
}

// TestSequentialInputs tests collision rate with sequential numeric inputs
func TestSequentialInputs(numInputs int) CollisionStats {
	start := time.Now()
	
	// Map hash to original input to detect TRUE collisions
	hashToInput := make(map[uint64]string)
	collisions := 0
	
	for i := 0; i < numInputs; i++ {
		// Generate sequential input
		data := fmt.Sprintf("input_%d", i)
		
		// Compute xxHash64
		hash := xxhash.Sum64String(data)
		
		if originalInput, exists := hashToInput[hash]; exists {
			// Hash collision detected - check if it's a TRUE collision
			if originalInput != data {
				collisions++
				fmt.Printf("TRUE COLLISION FOUND!\n")
				fmt.Printf("  Input 1: %s (hash: %d)\n", originalInput, hash)
				fmt.Printf("  Input 2: %s (hash: %d)\n", data, hash)
			}
		} else {
			hashToInput[hash] = data
		}
	}
	
	duration := time.Since(start)
	collisionRate := float64(collisions) / float64(numInputs)
	
	hashBits := 64.0
	expectedCollisions := (float64(numInputs) * float64(numInputs-1)) / (2.0 * math.Pow(2, hashBits))
	
	return CollisionStats{
		NumInputs:         numInputs,
		UniqueHashes:      len(hashToInput),
		Collisions:        collisions,
		CollisionRate:     collisionRate,
		ExpectedCollisions: expectedCollisions,
		Duration:          duration,
		DuplicateInputs:   0,
	}
}

// TestSimilarInputs tests collision rate with similar strings
func TestSimilarInputs(numInputs int, baseString string) CollisionStats {
	start := time.Now()
	
	// Map hash to original input to detect TRUE collisions
	hashToInput := make(map[uint64]string)
	collisions := 0
	
	for i := 0; i < numInputs; i++ {
		// Generate similar input by appending number
		data := fmt.Sprintf("%s_%d", baseString, i)
		
		// Compute xxHash64
		hash := xxhash.Sum64String(data)
		
		if originalInput, exists := hashToInput[hash]; exists {
			// Hash collision detected - check if it's a TRUE collision
			if originalInput != data {
				collisions++
				fmt.Printf("TRUE COLLISION FOUND!\n")
				fmt.Printf("  Input 1: %s (hash: %d)\n", originalInput, hash)
				fmt.Printf("  Input 2: %s (hash: %d)\n", data, hash)
			}
		} else {
			hashToInput[hash] = data
		}
	}
	
	duration := time.Since(start)
	collisionRate := float64(collisions) / float64(numInputs)
	
	hashBits := 64.0
	expectedCollisions := (float64(numInputs) * float64(numInputs-1)) / (2.0 * math.Pow(2, hashBits))
	
	return CollisionStats{
		NumInputs:         numInputs,
		UniqueHashes:      len(hashToInput),
		Collisions:        collisions,
		CollisionRate:     collisionRate,
		ExpectedCollisions: expectedCollisions,
		Duration:          duration,
		DuplicateInputs:   0,
	}
}

// TestWithDuplicateInputs demonstrates the difference between duplicate inputs and collisions
func TestWithDuplicateInputs(numInputs int, duplicatePercent float64) CollisionStats {
	start := time.Now()
	
	hashToInput := make(map[uint64]string)
	collisions := 0
	duplicateInputs := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Generate a pool of unique strings
	uniqueStrings := make([]string, int(float64(numInputs)*(1.0-duplicatePercent)))
	for i := range uniqueStrings {
		uniqueStrings[i] = GenerateRandomString(20, r)
	}
	
	for i := 0; i < numInputs; i++ {
		var data string
		
		// Randomly decide if we should use a duplicate input
		if r.Float64() < duplicatePercent && len(uniqueStrings) > 0 {
			// Use an existing string (duplicate input)
			data = uniqueStrings[r.Intn(len(uniqueStrings))]
		} else {
			// Generate new string
			data = GenerateRandomString(20, r)
		}
		
		hash := xxhash.Sum64String(data)
		
		if originalInput, exists := hashToInput[hash]; exists {
			if originalInput != data {
				collisions++
				fmt.Printf("TRUE COLLISION FOUND!\n")
				fmt.Printf("  Input 1: %s (hash: %d)\n", originalInput, hash)
				fmt.Printf("  Input 2: %s (hash: %d)\n", data, hash)
			} else {
				duplicateInputs++
			}
		} else {
			hashToInput[hash] = data
		}
	}
	
	duration := time.Since(start)
	collisionRate := float64(collisions) / float64(numInputs)
	
	hashBits := 64.0
	expectedCollisions := (float64(numInputs) * float64(numInputs-1)) / (2.0 * math.Pow(2, hashBits))
	
	return CollisionStats{
		NumInputs:         numInputs,
		UniqueHashes:      len(hashToInput),
		Collisions:        collisions,
		CollisionRate:     collisionRate,
		ExpectedCollisions: expectedCollisions,
		Duration:          duration,
		DuplicateInputs:   duplicateInputs,
	}
}

// PrintStats displays collision statistics
func PrintStats(testName string, stats CollisionStats) {
	fmt.Printf("\n=== %s ===\n", testName)
	fmt.Printf("Number of inputs:      %d\n", stats.NumInputs)
	fmt.Printf("Unique hashes:         %d\n", stats.UniqueHashes)
	fmt.Printf("TRUE collisions:       %d (different inputs, same hash)\n", stats.Collisions)
	if stats.DuplicateInputs > 0 {
		fmt.Printf("Duplicate inputs:      %d (same input, same hash - NOT collisions)\n", stats.DuplicateInputs)
	}
	fmt.Printf("Collision rate:        %.10f (%.4e)\n", stats.CollisionRate, stats.CollisionRate)
	fmt.Printf("Expected collisions:   %.10f (%.4e)\n", stats.ExpectedCollisions, stats.ExpectedCollisions)
	fmt.Printf("Duration:              %v\n", stats.Duration)
	
	if stats.Collisions > 0 {
		fmt.Printf("⚠️  WARNING: TRUE COLLISIONS DETECTED!\n")
	} else {
		fmt.Printf("✓ No true collisions detected\n")
	}
}

func main() {
	fmt.Println("xxHash64 Collision Rate Simulation")
	fmt.Println("===================================")
	fmt.Println("\nNOTE: We only count TRUE collisions (different inputs → same hash)")
	fmt.Println("Duplicate inputs (same input → same hash) are NOT collisions!\n")
	
	// Test 1: Small scale random test (100K inputs)
	stats1 := TestCollisionRate(100_000, 20)
	PrintStats("Random Strings (100K, length=20)", stats1)
	
	// Test 2: Medium scale random test (1M inputs)
	stats2 := TestCollisionRate(1_000_000, 20)
	PrintStats("Random Strings (1M, length=20)", stats2)
	
	// Test 3: Large scale random test (10M inputs)
	stats3 := TestCollisionRate(10_000_000, 20)
	PrintStats("Random Strings (10M, length=20)", stats3)
	
	// Test 4: Sequential inputs
	stats4 := TestSequentialInputs(1_000_000)
	PrintStats("Sequential Inputs (1M)", stats4)
	
	// Test 5: Similar strings
	stats5 := TestSimilarInputs(1_000_000, "user_data")
	PrintStats("Similar Strings (1M)", stats5)
	
	// Test 6: Short strings
	stats6 := TestCollisionRate(1_000_000, 5)
	PrintStats("Short Random Strings (1M, length=5)", stats6)
	
	// Test 7: Demonstration with duplicate inputs (20% duplicates)
	stats7 := TestWithDuplicateInputs(100_000, 0.20)
	PrintStats("With Duplicate Inputs (100K, 20% duplicates)", stats7)

	stats8 := TestCollisionRate(1_00_000_000, 20)
	PrintStats("Random String (100M, length=20)", stats8)

	//stats9 := TestCollisionRate(1_000_000_000, 20)
	//PrintStats("Random String (100M, length=20)", stats9)
	
	fmt.Println("\n=== Summary ===")
	fmt.Println("For xxHash64 (64-bit output):")
	fmt.Println("- Expected first collision after ~5 billion inputs (birthday paradox)")
	fmt.Println("- Collision probability for 1M inputs: ~0.00000000000005")
	fmt.Println("- Hash space size: 2^64 = 18,446,744,073,709,551,616")
	fmt.Println("\nKey Distinction:")
	fmt.Println("✓ TRUE COLLISION: Different inputs produce the same hash (BAD)")
	fmt.Println("✓ DUPLICATE INPUT: Same input produces the same hash (EXPECTED)")
}
