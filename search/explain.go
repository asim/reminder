package search

import (
	"fmt"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/quran"
)

// ExplanationPassage resolves the anchor from trusted corpus records before search.
func ExplanationPassage(q *quran.Quran, c model.ExplainContext) ([]*Result, error) {
	if c.Chapter == 0 {
		return nil, nil
	}
	if q == nil || c.Chapter > len(q.Chapters) || c.Chapter < 1 {
		return nil, fmt.Errorf("chapter not found")
	}
	chapter := q.Get(c.Chapter)
	if c.Verse > 0 {
		if chapter.Verse(c.Verse) == nil {
			return nil, fmt.Errorf("verse not found")
		}
		refs := VerseReferences(q, c.Chapter, c.Verse, 1)
		// Keep the focus first; supply the immediate passage context automatically.
		for _, number := range []int{c.Verse - 1, c.Verse + 1} {
			if number > 0 {
				refs = append(refs, VerseReferences(q, c.Chapter, number, 1)...)
			}
		}
		return refs, nil
	}
	// Put the chapter text before supplementary commentary so the budget covers as
	// much of the chapter as possible. The model must disclose partial coverage.
	var primary, extra []*Result
	for _, verse := range chapter.Verses {
		if verse.Number < 1 {
			continue
		}
		refs := VerseReferences(q, c.Chapter, verse.Number, 1)
		selected := 0
		for i, ref := range refs {
			if ref.Metadata["edition"] == c.Translation {
				selected = i
				break
			}
		}
		if len(refs) > 0 {
			primary = append(primary, refs[selected])
			extra = append(extra, refs[:selected]...)
			extra = append(extra, refs[selected+1:]...)
		}
	}
	return append(primary, extra...), nil
}

// HadithExplanationPassage resolves the exact collection record before ranked search.
func HadithExplanationPassage(b *hadith.Collection, c model.ExplainContext) ([]*Result, error) {
	collection := hadith.Select(c.Collection, b)
	if collection == nil {
		return nil, fmt.Errorf("hadith collection not found")
	}
	book := collection.Get(c.Book)
	if book == nil {
		return nil, fmt.Errorf("hadith book not found")
	}
	for _, h := range book.Hadiths {
		if h.Number == c.Hadith {
			return []*Result{hadithDocument(collection, book, h)}, nil
		}
	}
	return nil, fmt.Errorf("hadith not found")
}
