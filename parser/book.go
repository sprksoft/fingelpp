package parser

import (
	"fingelpp/finsyn"
	"html/template"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

type Les struct {
	Name    string
	Id      LessonId
	Content template.HTML
	Src     string
	path    string
}

type Chapter struct {
	Name    string
	Lessons []*Les
}

type Book struct {
	Chapters []*Chapter
	path     string
}

func LoadBook(path string) *Book {
	book := Book{Chapters: nil, path: path}
	book.Reload()
	return &book
}

func (book *Book) Reload() {
	chapterDirs, err := os.ReadDir(book.path)
	if err != nil {
		panic(err)
	}
	book.Chapters = nil

	for chapI, chapDir := range chapterDirs {
		_, chapName, _ := strings.Cut(chapDir.Name(), " ")

		lessonDirs, err := os.ReadDir(book.path + "/" + chapDir.Name())
		if err != nil {
			panic("Failed to read chapter: '" + chapName + "'" + err.Error())
		}

		lessons := []*Les{}
		for i, lessonFile := range lessonDirs {
			_, lessonName, _ := strings.Cut(lessonFile.Name(), " ")
			lessonName = strings.TrimSuffix(lessonName, ".md")
			path := book.path + "/" + chapDir.Name() + "/" + lessonFile.Name()
			les := Les{Name: lessonName, Id: LessonId{chapter: uint16(chapI), lesson: uint16(i)}, path: path}
			les.Reload()

			lessons = append(lessons, &les)
		}

		book.Chapters = append(book.Chapters, &Chapter{Name: chapName, Lessons: lessons})
	}
}

func (les *Les) Reload() {
	log.Info("Reading lesson " + les.path)
	content, err := os.ReadFile(les.path)
	if err != nil {
		panic("Failed to read lesson: '" + les.Name + "': " + err.Error())
	}
	les.Content = finsyn.ParseFinSyn(string(content))
	les.Src = string(content)
	log.Infof("Loaded lesson: %s", les.Name)
}

func (book *Book) GetChapterById(id uint16) *Chapter {
	if id >= uint16(len(book.Chapters)) {
		return nil
	}
	return book.Chapters[id]
}

func (book *Book) GetChapterByLessonId(id LessonId) *Chapter {
	return book.Chapters[id.chapter]
}

func (book *Book) GetLessonById(id LessonId) *Les {
	for _, chap := range book.Chapters {
		for _, lesson := range chap.Lessons {
			if lesson.Id == id {
				return lesson
			}
		}
	}
	return nil
}
