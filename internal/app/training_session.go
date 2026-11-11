package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/benabot/cli-drill/internal/progress"
)

type trainingPass struct {
	chapter.BatchResult
	immediate   progress.ChapterScore
	interrupted bool
}

func runTraining(opts Options, selected chapter.Chapter) error {
	state, err := progress.Load(progressPath(opts))
	if err != nil {
		return err
	}
	return runTrainingFrom(opts, selected, state)
}

func runTrainingFrom(opts Options, selected chapter.Chapter, state progress.Progress) error {
	if len(selected.Items) == 0 {
		return fmt.Errorf("chapter %s has no exercises", selected.ID)
	}
	hasKeys := selected.HasKeyExercises()
	if hasKeys && !canCaptureKeys(opts) {
		return errCaptureUnavailable
	}
	reader := bufio.NewReader(opts.In)
	view := newTextTrainingView(opts.Out)
	visit := chapter.NewVisit(selected, opts.AllExercises)

	for visit.HasNext() {
		batch := visit.Next()
		first, err := runTrainingPass(opts, reader, view, batch, &state)
		if err != nil {
			return err
		}
		if first.Total > 0 {
			view.summary("Bilan du lot", first, len(state.CompletedExercises))
		}
		if first.interrupted {
			return nil
		}

		latest := first
		for visit.HasNext() || len(latest.Missed) > 0 {
			action, err := readBatchAction(opts, reader, hasKeys, visit.HasNext(), len(latest.Missed) > 0)
			if err != nil {
				return err
			}
			switch action {
			case "c":
				// The next batch comes from the fixed visit, never from the retry.
			case "r":
				retry := batch
				retry.Items = latest.RetryItems()
				latest, err = runTrainingPass(opts, reader, view, retry, &state)
				if err != nil {
					return err
				}
				if latest.Total > 0 {
					view.summary("Bilan de la reprise", latest, len(state.CompletedExercises))
					_, _ = fmt.Fprintf(opts.Out, "Lot initial : %d/%d bonnes réponses\n", first.Correct, first.Total)
				}
				if latest.interrupted {
					return nil
				}
				continue
			default:
				return nil
			}
			break
		}
	}
	return nil
}

func runTrainingPass(opts Options, reader *bufio.Reader, view textTrainingView, batch chapter.Chapter, state *progress.Progress) (trainingPass, error) {
	pass := trainingPass{}
	textHeaderPrinted := false
	for index, item := range batch.Items {
		if item.ExerciseType == exercise.TypeKeySequence {
			for first := true; ; first = false {
				result, err := runKeySequenceExercise(opts, reader, batch, item, index, len(batch.Items), state)
				if err != nil && !errors.Is(err, io.EOF) {
					return pass, err
				}
				if result.Answered {
					if first {
						pass.Record(item, result.AnsweredKey, result.Correct)
					} else {
						pass.immediate.Attempts++
						if result.Correct {
							pass.immediate.Correct++
						}
					}
				}
				if result.Action == keySequenceQuit || errors.Is(err, io.EOF) {
					pass.interrupted = true
					return pass, nil
				}
				if result.Action != keySequenceRetry {
					break
				}
			}
			continue
		}
		if !textHeaderPrinted {
			view.chapterTitle(batch.Title)
			textHeaderPrinted = true
		}
		for {
			view.question(item, index+1, len(batch.Items))
			input, aborted, err := readTrainingInput(opts, reader, item.ExerciseType)
			if err != nil && !errors.Is(err, io.EOF) {
				return pass, err
			}
			if aborted {
				_, _ = fmt.Fprintln(opts.Out, "\nSession interrompue.")
				pass.interrupted = true
				return pass, nil
			}
			if strings.Trim(input, " \r\n") != "" {
				correct := item.MatchAnswer(input)
				recorded, saveErr := progress.RecordAnswer(progressPath(opts), batch.ID, item.ID, correct)
				if saveErr != nil {
					return pass, saveErr
				}
				*state = recorded
				pass.Record(item, strings.TrimRight(input, "\r\n"), correct)
				view.feedback(item, correct)
				if err == nil {
					break
				}
			}
			if errors.Is(err, io.EOF) {
				pass.interrupted = true
				return pass, nil
			}
		}
	}
	return pass, nil
}

func readBatchAction(opts Options, reader *bufio.Reader, raw, canContinue, canRetry bool) (string, error) {
	var actions []string
	if canContinue {
		actions = append(actions, "c Continuer")
	}
	if canRetry {
		actions = append(actions, "r Revoir mes erreurs")
	}
	actions = append(actions, "q Terminer")
	_, _ = fmt.Fprintln(opts.Out, strings.Join(actions, " · "))
	for {
		var input string
		var err error
		if raw {
			input, err = readKeySequence(opts, reader)
		} else {
			var aborted bool
			input, aborted, err = readTrainingInput(opts, reader, exercise.TypeFreeAnswer)
			if aborted {
				return "q", nil
			}
			if err == nil && strings.TrimSpace(input) == "" {
				input = "Enter"
			}
		}
		if err != nil && (!errors.Is(err, io.EOF) || input == "") {
			if errors.Is(err, io.EOF) {
				return "q", nil
			}
			return "", err
		}
		switch strings.ToLower(strings.TrimSpace(input)) {
		case "q", "esc", "ctrl+c":
			return "q", nil
		case "c":
			if canContinue {
				return "c", nil
			}
		case "r":
			if canRetry {
				return "r", nil
			}
		case "enter":
			if canContinue {
				return "c", nil
			}
			if canRetry {
				return "r", nil
			}
		}
	}
}
