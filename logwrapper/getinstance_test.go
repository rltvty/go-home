package logwrapper_test

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/rltvty/go-home/logwrapper"
	"github.com/rltvty/go-home/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rltvty/go-home/logwrapper"
)

var _ = Describe("GetInstance", func() {
	Context("When in Production", func() {
		var tempFile *os.File
		var err error
		var standardLogger *StandardLogger

		BeforeEach(func() {
			tempFile, err = ioutil.TempFile("", "getinstance_prod")
			if err != nil {
				log.Fatal(err)
			}
			standardLogger = logwrapper.GetInstance(func(config *Config) {
				config.Stdout = tempFile
			})
		})

		It("Should return a logger that outputs logs to stdout in json", func() {
			standardLogger.Info("Production Info Message")

			tempFile.Seek(0, 0)
			scanner := bufio.NewScanner(tempFile)
			scanner.Scan()
			jsonText := scanner.Text()

			Expect(testhelpers.IsJSON(jsonText)).To(BeTrue())
			jsonData := testhelpers.UnmarshalJSON(jsonText)
			Expect(jsonData["level"]).To(Equal("info"))
			Expect(jsonData["msg"]).To(Equal("Production Info Message"))
		})
		AfterEach(func() {
			os.Remove(tempFile.Name())
			logwrapper.ResetConfig()
		})
	})

	Context("When in Development", func() {
		var tempFile *os.File
		var err error
		var standardLogger *StandardLogger

		BeforeEach(func() {
			tempFile, err = ioutil.TempFile("", "getinstance_dev")
			if err != nil {
				log.Fatal(err)
			}
			standardLogger = GetInstance(func(config *Config) {
				config.Env = DEVELOPMENT
				config.Stdout = tempFile
			})
		})

		It("Should return a logger that outputs logs to stdout in text", func() {
			standardLogger.Info("Development Info Message")

			tempFile.Seek(0, 0)
			scanner := bufio.NewScanner(tempFile)
			scanner.Scan()
			text := scanner.Text()

			Expect(testhelpers.IsJSON(text)).To(BeFalse())
			Expect(text).To(ContainSubstring("INFO"))
			Expect(text).To(ContainSubstring("Development Info Message"))
		})
		AfterEach(func() {
			//log.Println(tempFile.Name())
			os.Remove(tempFile.Name())
			logwrapper.ResetConfig()
		})
	})
})
