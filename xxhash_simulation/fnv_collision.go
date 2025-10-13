package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"time"
)

// CollisionStats holds the results of a collision test
type CollisionStats struct {
	HashType          string
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

// ComputeFNV32a computes FNV-1a 32-bit hash
func ComputeFNV32a(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32()
}

// ComputeFNV64a computes FNV-1a 64-bit hash
func ComputeFNV64a(data string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(data))
	return h.Sum64()
}

// TestCollisionRateFNV32 tests FNV-1a 32-bit collision rate with random inputs
func TestCollisionRateFNV32(numInputs int, stringLength int) CollisionStats {
	start := time.Now()
	
	// Map hash to original input to detect TRUE collisions
	hashToInput := make(map[uint32]string)
	collisions := 0
	duplicateInputs := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < numInputs; i++ {
		// Generate random input
		data := GenerateRandomString(stringLength, r)
		
		// Compute FNV-1a 32-bit hash
		hash := ComputeFNV32a(data)
		
		if originalInput, exists := hashToInput[hash]; exists {
			// Hash collision detected - check if it's a TRUE collision
			if originalInput != data {
				collisions++
				fmt.Printf("FNV32 TRUE COLLISION FOUND!\n")
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
	hashBits := 32.0
	expectedCollisions := (float64(numInputs) * float64(numInputs-1)) / (2.0 * math.Pow(2, hashBits))
	
	if duplicateInputs > 0 {
		fmt.Printf("Note: %d duplicate inputs were generated (not counted as collisions)\n", duplicateInputs)
	}
	
	return CollisionStats{
		HashType:          "FNV-1a 32-bit",
		NumInputs:         numInputs,
		UniqueHashes:      len(hashToInput),
		Collisions:        collisions,
		CollisionRate:     collisionRate,
		ExpectedCollisions: expectedCollisions,
		Duration:          duration,
		DuplicateInputs:   duplicateInputs,
	}
}

// TestCollisionRateFNV64 tests FNV-1a 64-bit collision rate with random inputs
func TestCollisionRateFNV64(numInputs int, stringLength int) CollisionStats {
	start := time.Now()
	
	// Map hash to original input to detect TRUE collisions
	hashToInput := make(map[uint64]string)
	collisions := 0
	duplicateInputs := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < numInputs; i++ {
		// Generate random input
		data := GenerateRandomString(stringLength, r)
		
		// Compute FNV-1a 64-bit hash
		hash := ComputeFNV64a(data)
		
		if originalInput, exists := hashToInput[hash]; exists {
			// Hash collision detected - check if it's a TRUE collision
			if originalInput != data {
				collisions++
				fmt.Printf("FNV64 TRUE COLLISION FOUND!\n")
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
	hashBits := 64.0
	expectedCollisions := (float64(numInputs) * float64(numInputs-1)) / (2.0 * math.Pow(2, hashBits))
	
	if duplicateInputs > 0 {
		fmt.Printf("Note: %d duplicate inputs were generated (not counted as collisions)\n", duplicateInputs)
	}
	
	return CollisionStats{
		HashType:          "FNV-1a 64-bit",
		NumInputs:         numInputs,
		UniqueHashes:      len(hashToInput),
		Collisions:        collisions,
		CollisionRate:     collisionRate,
		ExpectedCollisions: expectedCollisions,
		Duration:          duration,
		DuplicateInputs:   duplicateInputs,
	}
}

// TestSequentialInputsFNV32 tests collision rate with sequential numeric inputs (32-bit)
func TestSequentialInputsFNV32(numInputs int) CollisionStats {
	start := time.Now()
	
	hashToInput := make(map[uint32]string)
	collisions := 0
	
	for i := 0; i < numInputs; i++ {
		data := fmt.Sprintf("input_%d", i)
		hash := ComputeFNV32a(data)
		
		if originalInput, exists := hashToInput[hash]; exists {
			if originalInput != data {
				collisions++
				fmt.Printf("FNV32 TRUE COLLISION FOUND!\n")
				fmt.Printf("  Input 1: %s (hash: %d)\n", originalInput, hash)
				fmt.Printf("  Input 2: %s (hash: %d)\n", data, hash)
			}
		} else {
			hashToInput[hash] = data
		}
	}
	
	duration := time.Since(start)
	collisionRate := float64(collisions) / float64(numInputs)
	hashBits := 32.0
	expectedCollisions := (float64(numInputs) * float64(numInputs-1)) / (2.0 * math.Pow(2, hashBits))
	
	return CollisionStats{
		HashType:          "FNV-1a 32-bit",
		NumInputs:         numInputs,
		UniqueHashes:      len(hashToInput),
		Collisions:        collisions,
		CollisionRate:     collisionRate,
		ExpectedCollisions: expectedCollisions,
		Duration:          duration,
		DuplicateInputs:   0,
	}
}

// TestSequentialInputsFNV64 tests collision rate with sequential numeric inputs (64-bit)
func TestSequentialInputsFNV64(numInputs int) CollisionStats {
	start := time.Now()
	
	hashToInput := make(map[uint64]string)
	collisions := 0
	
	for i := 0; i < numInputs; i++ {
		data := fmt.Sprintf("input_%d", i)
		hash := ComputeFNV64a(data)
		
		if originalInput, exists := hashToInput[hash]; exists {
			if originalInput != data {
				collisions++
				fmt.Printf("FNV64 TRUE COLLISION FOUND!\n")
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
		HashType:          "FNV-1a 64-bit",
		NumInputs:         numInputs,
		UniqueHashes:      len(hashToInput),
		Collisions:        collisions,
		CollisionRate:     collisionRate,
		ExpectedCollisions: expectedCollisions,
		Duration:          duration,
		DuplicateInputs:   0,
	}
}

// TestWithDuplicateInputsFNV32 demonstrates the difference between duplicate inputs and collisions
func TestWithDuplicateInputsFNV32(numInputs int, duplicatePercent float64) CollisionStats {
	start := time.Now()
	
	hashToInput := make(map[uint32]string)
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
		
		if r.Float64() < duplicatePercent && len(uniqueStrings) > 0 {
			data = uniqueStrings[r.Intn(len(uniqueStrings))]
		} else {
			data = GenerateRandomString(20, r)
		}
		
		hash := ComputeFNV32a(data)
		
		if originalInput, exists := hashToInput[hash]; exists {
			if originalInput != data {
				collisions++
				fmt.Printf("FNV32 TRUE COLLISION FOUND!\n")
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
	hashBits := 32.0
	expectedCollisions := (float64(numInputs) * float64(numInputs-1)) / (2.0 * math.Pow(2, hashBits))
	
	return CollisionStats{
		HashType:          "FNV-1a 32-bit",
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
	fmt.Printf("\n=== %s [%s] ===\n", testName, stats.HashType)
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
	fmt.Println("FNV-1a Collision Rate Simulation")
	fmt.Println("=================================")
	fmt.Println("\nNOTE: We only count TRUE collisions (different inputs → same hash)")
	fmt.Println("Duplicate inputs (same input → same hash) are NOT collisions!\n")
	
	fmt.Println("\n========== FNV-1a 32-bit Tests ==========")
	
	// FNV32 Test 1: Small scale (10K inputs - likely to see collisions)
	stats1 := TestCollisionRateFNV32(10_000, 20)
	PrintStats("Random Strings (10K, length=20)", stats1)
	
	// FNV32 Test 2: Medium scale (100K inputs - very likely to see collisions)
	stats2 := TestCollisionRateFNV32(100_000, 20)
	PrintStats("Random Strings (100K, length=20)", stats2)
	
	// FNV32 Test 3: Sequential inputs (100K)
	stats3 := TestSequentialInputsFNV32(100_000)
	PrintStats("Sequential Inputs (100K)", stats3)
	
	// FNV32 Test 4: With duplicate inputs
	stats4 := TestWithDuplicateInputsFNV32(50_000, 0.20)
	PrintStats("With Duplicate Inputs (50K, 20% duplicates)", stats4)

	stats21 := TestCollisionRateFNV32(10_000_000, 20)
	PrintStats("Random Strings (10M, length=20)", stats21)
	
	fmt.Println("\n========== FNV-1a 64-bit Tests ==========")
	
	// FNV64 Test 1: Small scale (100K inputs)
	stats5 := TestCollisionRateFNV64(100_000, 20)
	PrintStats("Random Strings (100K, length=20)", stats5)
	
	// FNV64 Test 2: Medium scale (1M inputs)
	stats6 := TestCollisionRateFNV64(1_000_000, 20)
	PrintStats("Random Strings (1M, length=20)", stats6)
	
	// FNV64 Test 3: Sequential inputs (1M)
	stats7 := TestSequentialInputsFNV64(1_000_000)
	PrintStats("Sequential Inputs (1M)", stats7)

	stats51 := TestCollisionRateFNV64(10_000_000, 20)
	PrintStats("Random Strings (10M, length=20)", stats51)

	stats52 := TestCollisionRateFNV64(1_00_000_000, 20)
	PrintStats("Rangom Strings (100M, length=20)", stats52)

	fmt.Println("\n=== Summary ===")
	fmt.Println("\nFNV-1a 32-bit (32-bit output):")
	fmt.Println("- Hash space size: 2^32 = 4,294,967,296")
	fmt.Println("- Expected first collision after ~77,000 inputs (birthday paradox)")
	fmt.Println("- Collision probability for 100K inputs: ~0.001 (0.1%)")
	fmt.Println("- You WILL see collisions with 100K+ inputs")
	
	fmt.Println("\nFNV-1a 64-bit (64-bit output):")
	fmt.Println("- Hash space size: 2^64 = 18,446,744,073,709,551,616")
	fmt.Println("- Expected first collision after ~5 billion inputs (birthday paradox)")
	fmt.Println("- Collision probability for 1M inputs: ~0.00000000000005")
	fmt.Println("- You should NOT see collisions with < 100M inputs")
	
	fmt.Println("\nKey Distinction:")
	fmt.Println("✓ TRUE COLLISION: Different inputs produce the same hash (BAD)")
	fmt.Println("✓ DUPLICATE INPUT: Same input produces the same hash (EXPECTED)")

}
