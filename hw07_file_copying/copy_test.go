package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("copy the entire file", func(t *testing.T) {
		src := "testdata/input.txt"
		dst := "output.txt"

		err := Copy(src, dst, 0, 0)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		defer os.Remove(dst)

		origContent, err := os.ReadFile("testdata/out_offset0_limit0.txt")
		if err != nil {
			t.Fatalf("Fail read origin file: %v", err)
		}

		copiedContent, err := os.ReadFile(dst)
		if err != nil {
			t.Fatalf("Fail read copied file: %v", err)
		}

		require.Equal(t, origContent, copiedContent)
	})

	t.Run("copy the middle of the file", func(t *testing.T) {
		src := "testdata/input.txt"
		dst := "output.txt"

		err := Copy(src, dst, 100, 1000)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		defer os.Remove(dst)

		origContent, err := os.ReadFile("testdata/out_offset100_limit1000.txt")
		if err != nil {
			t.Fatalf("Fail read origin file: %v", err)
		}

		copiedContent, err := os.ReadFile(dst)
		if err != nil {
			t.Fatalf("Fail read copied file: %v", err)
		}

		require.Equal(t, origContent, copiedContent)
	})
}
