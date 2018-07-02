package main

import (
	"fmt"
	"math/rand"
)

var list1 = []string{
	"rare",
	"sugared",
	"precious",
	"dutiful",
	"damasked",
	"flowering",
	"gallant",
	"celestial",
	"sweet",
	"saucy",
	"sportful",
	"artful",
	"heavenly",
	"yarely",
	"tuneful",
	"courteous",
	"delicate",
	"silken",
	"brave",
	"complete",
	"vasty",
	"pleasing",
	"cheek - rosy",
	"deserving",
	"melting",
	"wholesome",
	"fruitful",
}

var list2 = []string{
	"honey-tongued",
	"well-wishing",
	"berhyming",
	"fair-faced",
	"five-fingered-tied",
	"heart-inflaming",
	"not-answering",
	"spleenative",
	"softly-sprighted",
	"smooth-faced",
	"sweet-suggesting",
	"swinge-buckling",
	"tender-hearted",
	"tender-feeling",
	"thunder-darting",
	"tiger-booted",
	"lustyhooded",
	"time-pleasing welsh",
	"superstitious",
	"sympathizing",
	"sweet-tongued",
	"weeping-ripe",
	"well-favoured",
	"young-eyed",
	"sweet-mouthing primrose",
	"best-tempered",
	"well-graced",
}

var list3 = []string{
	"nymph",
	"ornament",
	"toast",
	"curiosity",
	"apple-john",
	"bilbo",
	"cuckoo-bud",
	"nose-herb",
	"gamester",
	"ouch",
	"goddess",
	"night-cap",
	"delight",
	"watercake",
	"umpire",
	"sprite",
	"song",
	"cheese",
	"kissing-comfit",
	"wit-cracker",
	"hawthorn-bud",
	"valentine",
	"smilet",
	"true-penny",
	"path",
	"gaudy-night",
	"pigeon-egg",
}

func compliment() string {
	first := list1[rand.Intn(len(list1))]
	second := list2[rand.Intn(len(list2))]
	third := list3[rand.Intn(len(list3))]
	return fmt.Sprintf("%s %s %s", first, second, third)
}
