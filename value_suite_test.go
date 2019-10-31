package value

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGotable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gotable Suite")
}
