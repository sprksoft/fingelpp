package parser

import (
	"fingelpp/finsyn"
	"github.com/charmbracelet/log"
	"html/template"
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
	path     string
}

func LoadLessons(path string) LessonManager {
	mgr := LessonManager{Chapters: nil, path: path}
	mgr.Reload()
	return mgr
}

func (mgr *LessonManager) Reload() {
	chapterDirs, err := os.ReadDir(mgr.path)
	if err != nil {
		panic(err)
	}
	mgr.Chapters = nil

	for chapI, chapDir := range chapterDirs {
		_, chapName, _ := strings.Cut(chapDir.Name(), " ")

		lessonDirs, err := os.ReadDir(mgr.path + "/" + chapDir.Name())
		if err != nil {
			panic("Failed to read chapter: '" + chapName + "'" + err.Error())
		}

		lessons := []Les{}
		for i, lessonFile := range lessonDirs {
			_, lessonName, _ := strings.Cut(lessonFile.Name(), " ")
			lessonName = strings.TrimSuffix(lessonName, ".md")

			content, err := os.ReadFile(mgr.path + "/" + chapDir.Name() + "/" + lessonFile.Name())
			if err != nil {
				panic("Failed to read lesson: '" + lessonName + "'" + err.Error())
			}
			parsedContent := finsyn.ParseFinSyn(string(content))

			les := Les{Name: lessonName, Id: LessonId{chapter: uint16(chapI), lesson: uint16(i)}, Content: parsedContent}

			log.Infof("Loaded lesson: %s", les.Name)
			lessons = append(lessons, les)
		}

		mgr.Chapters = append(mgr.Chapters, Chapter{Name: chapName, Lessons: lessons})
	}
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
