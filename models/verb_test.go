package models

import (
	"testing"
)

type verbTestCase struct {
	Word Verb

	Past string
	Future string
	Present_It string
	Present_You string
}

var verbtestdata []verbTestCase = []verbTestCase{
	verbTestCase{
		Word: Verb{
				Regular: true,
				Word: "pull",
			},
		Past: "pulled",
		Future: "will pull",
		Present_It: "pulls",
		Present_You: "pull",
	},

	verbTestCase{
		Word: Verb{
				Regular: true,
				Word: "push",
			},
		Past: "pushed",
		Future: "will push",
		Present_It: "pushes",
		Present_You: "push",
	},

	verbTestCase{
		Word: Verb{
				Regular: true,
				Word: "clone",
			},
		Past: "cloned",
		Future: "will clone",
		Present_It: "clones",
		Present_You: "clone",
	},

	verbTestCase{
		Word: Verb{
				Regular: true,
				Word: "identify",
			},
		Past: "identified",
		Future: "will identify",
		Present_It: "identifies",
		Present_You: "identify",
	},

	verbTestCase{
		Word: Verb{
				Regular: true,
				Word: "fetch",
			},
		Past: "fetched",
		Future: "will fetch",
		Present_It: "fetches",
		Present_You: "fetch",
	},

	verbTestCase{
		Word: Verb{
				Regular: true,
				Word: "access",
			},
		Past: "accessed",
		Future: "will access",
		Present_It: "accesses",
		Present_You: "access",
	},

	verbTestCase{
		Word: Verb{
				Regular: true,
				Word: "index",
			},
		Past: "indexed",
		Future: "will index",
		Present_It: "indexes",
		Present_You: "index",
	},
}

func Test_Verbs(t *testing.T) {
	for _, td := range verbtestdata {
		past := td.Word.Conjugate(CT_It, CM_Past, false)
		presentIt := td.Word.Conjugate(CT_It, CM_Present, false)
		presentYou := td.Word.Conjugate(CT_You, CM_Present, false)
		future := td.Word.Conjugate(CT_It, CM_Future, false)

		if past != td.Past {
			t.Errorf("past failed: %s != %s", past, td.Past)
		}

		if presentIt != td.Present_It {
			t.Errorf("presentIt failed: %s != %s", presentIt, td.Present_It)
		}

		if presentYou != td.Present_You {
			t.Errorf("presentYou failed: %s != %s", presentYou, td.Present_You)
		}

		if future != td.Future {
			t.Errorf("future failed: %s != %s", future, td.Future)
		}
	}
}
