package words

import (
	"context"
	"log"
	"strings"

	"github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wordspb "yadro.com/course/proto/words"
)

func Norm(_ context.Context, in *wordspb.WordsRequest) (*wordspb.WordsReply, error) {
	if len(in.Phrase) > 4096 {
		return nil, status.Errorf(codes.ResourceExhausted, "Phrase lenght > 4 KiB")
	}

	lowerPhrase := strings.ToLower(in.Phrase)
	phraseSlice := strings.FieldsFunc(lowerPhrase, func(r rune) bool {
		isAlphanumeric := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		return !isAlphanumeric
	})

	uniqueStems := make(map[string]struct{})

	for _, word := range phraseSlice {
		if word == "" {
			continue
		}

		if english.IsStopWord(word) {
			continue
		}

		stemmedWord, err := snowball.Stem(word, "english", false)
		if err != nil {
			log.Printf("Error stemming word '%s': %v", word, err)
			continue
		}

		uniqueStems[stemmedWord] = struct{}{}
	}

	result := make([]string, 0, len(uniqueStems))
	for stem := range uniqueStems {
		result = append(result, stem)
	}

	return &wordspb.WordsReply{
		Words: result,
	}, nil
}
