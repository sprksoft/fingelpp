package lessons

import (
	"fmt"
	"strconv"
	"strings"
)

type LessonId struct {
	chapter uint16
	lesson  uint16
}

func (lid LessonId) String() string {
	return fmt.Sprintf("%v.%v", lid.chapter, lid.lesson)
}

func (lid LessonId) ChapterId() uint16 {
	return lid.chapter
}

func ParseLessonId(str string) (LessonId, error) {
	chapIndexStr, lessonIndexStr, _ := strings.Cut(str, ".")

	chapIndex, err := strconv.ParseUint(chapIndexStr, 10, 16)
	if err != nil {
		return LessonId{}, err
	}
	lessonIndex, err := strconv.ParseUint(lessonIndexStr, 10, 16)
	if err != nil {
		return LessonId{}, err
	}

	return LessonId{chapter: uint16(chapIndex), lesson: uint16(lessonIndex)}, nil
}
