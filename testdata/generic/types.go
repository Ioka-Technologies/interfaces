package generic

import (
	"github.com/rjeczalik/interfaces/testdata/util"
)

type MyGeneric[A any] struct {
	Value []A
}

type MyGenericAlias = MyGeneric[util.MyUtil]

type MyGenericAliasWithTypeArg[B any] = MyGeneric[B]
