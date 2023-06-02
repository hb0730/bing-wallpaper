package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBingWallpaperUpload(t *testing.T) {
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"upload", "--type", "s3", "--index", "0", "--market", "zh-CN", "--resolution", "1920"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
	t.Log(actual.String())
}
