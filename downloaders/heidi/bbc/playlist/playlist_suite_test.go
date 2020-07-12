package playlist_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlaylist(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Playlist Suite")
}
