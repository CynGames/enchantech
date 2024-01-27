package use_cases

import "fmt"

type TrackArticleReadUseCase struct {
}

func NewTrackArticleReadUseCase() *TrackArticleReadUseCase {
	return &TrackArticleReadUseCase{}
}

func (u *TrackArticleReadUseCase) Execute() error {
	return fmt.Errorf("unknown error")
}
