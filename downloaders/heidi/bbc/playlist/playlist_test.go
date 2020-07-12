package playlist_test

import (
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

	Describe("ParseTrack", func() {
		var rootNode *html.Node
		BeforeEach(func() {
			rootNode = htmlutils.GetHTMLFromFile("./test_data/playlist.htm")
			htmlutils.RemoveEmpty(rootNode)
		})
		It("should return info for each track", func() {
			trackNodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
				return htmlutils.IncludesAttr(n, "class", "segment__track")
			})
			tracks := []Track{}
			for _, trackNode := range trackNodes {
				tracks = append(tracks, ParseTrack(trackNode))
			}
			Expect(len(tracks)).To(Equal(8))

			Expect(tracks[0].Artist).To(Equal("Emmanuelle"))
			Expect(tracks[0].Name).To(Equal("Italove"))
			Expect(tracks[0].Album).To(Equal(""))
			Expect(tracks[0].Label).To(Equal("DEEWEE"))
			Expect(tracks[0].Info).To(Equal(""))

			Expect(tracks[1].Artist).To(Equal("Krankbrother"))
			Expect(tracks[1].Name).To(Equal("Right There With You"))
			Expect(tracks[1].Album).To(Equal(""))
			Expect(tracks[1].Label).To(Equal(""))
			Expect(tracks[1].Info).To(Equal(""))

			Expect(tracks[2].Artist).To(Equal("Stereogamous"))
			Expect(tracks[2].Name).To(Equal("Don’t Fight It (Remix) (feat. Shaun J. Wright)"))
			Expect(tracks[2].Album).To(Equal(""))
			Expect(tracks[2].Label).To(Equal("Twirl Recordings"))
			Expect(tracks[2].Info).To(Equal("Remix Artist: Alinka"))

			Expect(tracks[3].Artist).To(Equal("Farley & Severino"))
			Expect(tracks[3].Name).To(Equal("Music Fills the Air (Remix) (feat. Roy Inc)"))
			Expect(tracks[3].Album).To(Equal(""))
			Expect(tracks[3].Label).To(Equal("SoSure Music"))
			Expect(tracks[3].Info).To(Equal("Remix Artist: Hard Ton"))

			Expect(tracks[4].Artist).To(Equal("Seth Troxler & Tom Trago"))
			Expect(tracks[4].Name).To(Equal("De Natte Cel"))
			Expect(tracks[4].Album).To(Equal(""))
			Expect(tracks[4].Label).To(Equal(""))
			Expect(tracks[4].Info).To(Equal("Remix Artist: Prins Thomas Diskomiks"))

			Expect(tracks[5].Artist).To(Equal("Felix Da Housecat, Jamie Principle, & Vince Lawrence AKA The 312"))
			Expect(tracks[5].Name).To(Equal("Touch Your Body"))
			Expect(tracks[5].Album).To(Equal(""))
			Expect(tracks[5].Label).To(Equal("Crosstown Rebels"))
			Expect(tracks[5].Info).To(Equal("Remix Artist: Moodymann"))

			Expect(tracks[6].Artist).To(Equal("Solomun"))
			Expect(tracks[6].Name).To(Equal("Watergate"))
			Expect(tracks[6].Album).To(Equal("Watergate 11"))
			Expect(tracks[6].Label).To(Equal("Watergate Records"))
			Expect(tracks[6].Info).To(Equal(""))

			Expect(tracks[7].Artist).To(Equal("Max Berlin"))
			Expect(tracks[7].Name).To(Equal("Elle Et Moi (Joakim Remix)"))
			Expect(tracks[7].Album).To(Equal(""))
			Expect(tracks[7].Label).To(Equal("Eighttrack Recordings"))
			Expect(tracks[7].Info).To(Equal(""))
		})
	})

	/*
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
	*/

})
