package tafsir

import (
	"github.com/asim/reminder/internal/corpusversion"
	"sync"
)

var corpusVersion = sync.OnceValue(func() string { return corpusversion.Sum(files) })

func ContentVersion() string { return corpusVersion() }
