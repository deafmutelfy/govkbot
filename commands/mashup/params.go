package mashup

import (
	"errors"
	"strconv"
	"vkbot/core"
)

type params struct {
	Longest bool

	FirstBass   int
	FirstTreble int
	FirstOffset int

	SecondBass   int
	SecondTreble int
	SecondOffset int
}

func parseParams(args *[]string) (*params, error) {
	p := &params{}

	n := len(*args)

	if n >= 1 {
		if core.IsInArray([]string{"1", "да", "+", "true"}, (*args)[0]) {
			p.Longest = true
		}
	}
	if n >= 2 {
		m, err := strconv.Atoi((*args)[1])
		if err != nil {
			return nil, errors.New("недопустимое значение громкости нижних частот первой аудиозаписи")
		}

		p.FirstBass = m
	}
	if n >= 3 {
		m, err := strconv.Atoi((*args)[2])
		if err != nil {
			return nil, errors.New("недопустимое значение громкости верхних частот первой аудиозаписи")
		}

		p.FirstTreble = m
	}
	if n >= 4 {
		m, err := strconv.ParseFloat((*args)[3], 32)
		if err != nil {
			return nil, errors.New("недопустимое значение задержки первой аудиозаписи")
		}

		p.FirstOffset = int(m * 1000)
	}
	if n >= 5 {
		m, err := strconv.Atoi((*args)[4])
		if err != nil {
			return nil, errors.New("недопустимое значение громкости нижних частот второй аудиозаписи")
		}

		p.SecondBass = m
	}
	if n >= 6 {
		m, err := strconv.Atoi((*args)[5])
		if err != nil {
			return nil, errors.New("недопустимое значение громкости верхних частот второй аудиозаписи")
		}

		p.SecondTreble = m
	}
	if n >= 7 {
		m, err := strconv.ParseFloat((*args)[6], 32)
		if err != nil {
			return nil, errors.New("недопустимое значение задержки второй аудиозаписи")
		}

		p.SecondOffset = int(m * 1000)
	}

	return p, nil
}
