package parser

import (
	"html/template"
	"log"
	"os"
	"strings"
)

type Les struct {
	Name    string
	Id      LessonId
	Content template.HTML
}

type Chapter struct {
	Name    string
	Lessons []Les
}

type LessonManager struct {
	Chapters []Chapter
}

func LoadLessons(path string) LessonManager {

	chapterDirs, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	chapters := []Chapter{}

	for chapI, chapDir := range chapterDirs {
		_, chapName, _ := strings.Cut(chapDir.Name(), " ")

		lessonDirs, err := os.ReadDir(path + "/" + chapDir.Name())
		if err != nil {
			panic("Failed to read chapter: '" + chapName + "'" + err.Error())
		}

		lessons := []Les{}
		for i, lessonFile := range lessonDirs {
			_, lessonName, _ := strings.Cut(lessonFile.Name(), " ")
			lessonName = strings.TrimSuffix(lessonName, ".md")

			content, err := os.ReadFile(path + "/" + chapDir.Name() + "/" + lessonFile.Name())
			if err != nil {
				panic("Failed to read lesson: '" + lessonName + "'" + err.Error())
			}
			parsedContent := ParseFinSyn(string(content))

			les := Les{Name: lessonName, Id: LessonId{chapter: uint16(chapI), lesson: uint16(i)}, Content: parsedContent}

			log.Printf("Loaded lesson: %s", les.Name)
			lessons = append(lessons, les)
		}

		chapters = append(chapters, Chapter{Name: chapName, Lessons: lessons})

	}

	return LessonManager{Chapters: chapters}
}

func (mgr *LessonManager) GetChapterById(id uint16) *Chapter {
	if id >= uint16(len(mgr.Chapters)) {
		return nil
	}
	return &mgr.Chapters[id]
}

func (mgr *LessonManager) GetChapterByLessonId(id LessonId) *Chapter {
	return &mgr.Chapters[id.chapter]
}

func (mgr *LessonManager) GetLessonById(id LessonId) *Les {
	for _, chap := range mgr.Chapters {
		for _, lesson := range chap.Lessons {
			if lesson.Id == id {
				return &lesson
			}
		}
	}
	return nil
}
