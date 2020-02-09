package playlist_test

import (
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/html"

	. "github.com/rltvty/go-home/downloaders/heidi/bbc/playlist"
	"github.com/rltvty/go-home/htmlutils"
)

var _ = Describe("Playlist", func() {
	testFiles := []string{
		"./test_data/1-jun-12.htm",
		"./test_data/2-aug-13.htm",
		"./test_data/15-aug-16.htm",
		"./test_data/26-may-16.htm",
	}
	rootNodes := map[string]*html.Node{}
	BeforeEach(func() {
		for _, testFile := range testFiles {
			rootNodes[testFile] = htmlutils.GetHTMLFromFile(testFile)
		}
	})

	Describe("Remove NotNeeded", func() {
		BeforeEach(func() {
			for _, rootNode := range rootNodes {
				htmlutils.RemoveJunk(rootNode)
				RemoveBBCNotNeeded(rootNode)
			}
		})
		It("should return html with no bbc header, footer, navigation, images", func() {
			for _, rootNode := range rootNodes {
				formattedHtml, _ := htmlutils.RenderHTMLNode(rootNode)
				Expect(formattedHtml).ToNot(ContainSubstring("<img"))
				Expect(formattedHtml).ToNot(ContainSubstring("<header"))
				Expect(formattedHtml).ToNot(ContainSubstring("<footer"))
				Expect(formattedHtml).ToNot(ContainSubstring("<nav"))
			}
		})
	})

	Describe("GeneratePlaylist", func() {
		It("should return html with no bbc header, footer, navigation, images", func() {
			for testFile, rootNode := range rootNodes {
				playlist := GetPlaylist(rootNode)
				goal, _ := ioutil.ReadFile(fmt.Sprintf("%s.nfo", testFile))
				Expect(len(goal)).ToNot(Equal(0))
				Expect(playlist).To(Equal(string(goal)))
			}
		})
	})

})
