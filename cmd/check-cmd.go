package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"
	"wait4it/model"
)

func RunCheck(c model.CheckContext) {
	m, err := findCheckModule(c.Config.CheckType)
	if err != nil {
		wStdErr(err)
		os.Exit(2)
	}

	cx := m.(model.CheckInterface)

	cx.BuildContext(c)
	r, err := cx.Validate()
	if !r {
		wStdErr(err)
		os.Exit(2)
	}

	fmt.Print("Wait4it...")

	t := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go ticker(cx, t, done)

	time.Sleep(time.Duration(c.Config.Timeout) * time.Second)
	done <- true

	fmt.Print("failed")
	os.Exit(1)
}
func findCheckModule(ct string) (interface{}, error) {
	m, ok := cm[ct]
	if !ok {
		return nil, errors.New("unsupported check type")
	}

	return m, nil
}

func ticker(cs model.CheckInterface, t *time.Ticker, d chan bool) {
	for {
		select {
		case <-d:
			return
		case <-t.C:
			check(cs)
		}
	}
}

func check(cs model.CheckInterface) {
	r, eor, err := cs.Check()
	if err != nil && eor {
		wStdErr(err.Error())
		os.Exit(2)
	}

	wStdOut(r)
}
